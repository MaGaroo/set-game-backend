package service

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"set-game/src/models"
	"set-game/src/service/messages"
	"strings"
	"strconv"
)

func (service *Service) setupRoutes() {
	log.Printf("Setting up routes\n")
	http.HandleFunc("/ws", service.wsEndpoint)
	http.HandleFunc("/create-room", service.createRoomRequestHandler)
}

func (service *Service) createRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PostForm.Get("token")
	log.Printf("Creating room(%v)\n", token)
	if err := service.createRoom(token); err != nil {
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
		// TODO write response

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
		messageSlice := strings.Fields(string(p))
		if messageType != 1 || len(messageSlice) != 7 || messageSlice[0] != "guess" {
			continue
		}

		// parse the message
		log.Println("Parsing a guess")
		var guess [3][2]int
		badFormat := false
		for i := 0; i < 6; i++ {
			guess[i>>1][i&1], err = strconv.Atoi(messageSlice[i+1])
			if err != nil {
				badFormat = true
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
				conn.WriteMessage(1, messages.UpdateScore{
					Username: player.Username,
					Score:    player.Score,
				}.ToBytes())
				conn.WriteMessage(1, messages.UpdateCard{
					Cards: newCards,
				}.ToBytes())
				if endGame {
					conn.WriteMessage(1, messages.EndGame{
						Winner: winner,
					}.ToBytes())
				}
			}(newCards, service.connections[listener.Token])
		}
	}

}
