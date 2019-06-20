package service

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"set-game/src/models"
	"set-game/src/service/messages"
	"strings"
	"encoding/json"
)

func (service *Service) setupRoutes() {
	log.Printf("Setting up routes\n")
	http.HandleFunc("/ws", service.wsEndpoint)
	http.HandleFunc("/create-room", service.createRoomRequestHandler)
}

func (service *Service) createRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	// TODO remove this TOF!
	w.Header().Set("Access-Control-Allow-Origin", "*")
	token, ok := r.URL.Query()["token"]
	if !ok || len(token) == 0 || len(token[0]) != 10 {
		log.Printf("Bad token")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad token"))
		return
	}
	log.Printf("Creating room(%s)\n", token[0])
	if err := service.createRoom(token[0]); err != nil {
		log.Fatalf("Failed to create room(%v): %s\n", token, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (service *Service) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Printf("Creating a new websoket\n")
	service.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := service.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to upgrade connection to websocket\n")
		log.Println(err.Error())
	}
	defer ws.Close()

	player, room := service.register(ws)
	if player != nil {
		service.play(player, room, ws)
	}
}

func (service *Service) register(conn *websocket.Conn) (*models.Player, *models.Room) {
	log.Printf("Registering a new player\n")
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Could not register: read error: %s\n", err.Error())
			return nil, nil
		}
		if messageType != 1 {
			continue
		}

		req := messages.ParseIntroRequest(string(p))
		if req == nil {
			continue
		}

		room, err := service.getRoom(req.RoomToken)
		if err != nil {
			log.Fatalf("Could not register: could not get room: %s", err.Error())
			return nil, nil
		}

		player, err := service.getOrCreatePlayer(req.PlayerToken, req.Username, room)
		if err != nil {
			log.Fatalf("Could not register: could not get or create player: %s", err.Error())
			return nil, nil
		}
		service.connections[player.Token] = conn
		log.Printf("Successfully registered player(%v) in room(%v)", player.Username, room.Token)

		log.Println("Sending initial data to the joined player")
		game := service.getGame(room)
		var tableCards []int
		json.Unmarshal([]byte(game.Table), &tableCards)
		cardsMessage := messages.UpdateCard{
			Cards: tableCards,
		}
		conn.WriteMessage(1, cardsMessage.ToBytes())

		var members []models.Player
		service.getRoomPeople(room, &members)
		for _, member := range members {
			scoreMessage := messages.UpdateScore{
				Username: member.Username,
				Score:    member.Score,
			}
			conn.WriteMessage(1, scoreMessage.ToBytes())
		}

		return player, room
	}
	return nil, nil
}

func (service *Service) play(player *models.Player, room *models.Room, conn *websocket.Conn) {
	log.Printf("Player(%v) started to play in room(%v)\n", player.Token, room.Token)
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Could not play: read error: %s\n", err.Error())
			return
		}

		// TODO move the following part(parse) to messages package
		// basic format checking of the message
		log.Println("Parsing a guess")
		messageSlice := strings.Fields(string(p))
		if messageType != 1 || len(messageSlice) != 4 || messageSlice[0] != "guess" {
			continue
		}

		// parse the message
		guess := make([]int, 3)
		badFormat := false
		for i := 0; i < 3; i++ {
			for j := 0; j < i; j++ {
				if guess[i] == guess[j] {
					badFormat = true
				}
			}
		}
		if badFormat {
			continue
		}

		// TODO use mutex

		// get game of this room
		game := service.getGame(room)
		if game == nil {
			continue
		}

		// check if the guess was true
		log.Println("Checking player's guess")
		trueGuess, newCards, endGame := game.Check(guess)
		service.saveGame(game)
		if !trueGuess {
			continue
		}

		// increase player score
		log.Printf("Increasing player(%v) score\n", player.Username)
		player.Score += 1
		service.savePlayer(player)

		// find members of this room
		var listeners []models.Player
		service.getRoomPeople(room, &listeners)

		winner := ""
		winnerScore := int8(0)
		if endGame {
			log.Println("Finding winner")
			for _, person := range listeners {
				if person.Score > winnerScore {
					winner = person.Username
					winnerScore = person.Score
				}
			}
		}

		// send updates to people
		log.Println("Sending updates to members of a room")
		for _, listener := range listeners {
			go func(cards []int, conn *websocket.Conn) {
				if conn == nil {
					return
				}
				updateScore := messages.UpdateScore{
					Username: player.Username,
					Score:    player.Score,
				}
				conn.WriteMessage(1, updateScore.ToBytes())
				updateCard := messages.UpdateCard{
					Cards: newCards,
				}
				conn.WriteMessage(1, updateCard.ToBytes())
				if endGame {
					endGame := messages.EndGame{
						Winner: winner,
					}
					conn.WriteMessage(1, endGame.ToBytes())
				}
			}(newCards, service.connections[listener.Token])
		}
	}

}
