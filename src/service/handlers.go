package service

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"fmt"
	"set-game/src/models"
	"set-game/src/service/messages"
	"strings"
	"strconv"
)

func (service *Service) setupRoutes() {
	http.HandleFunc("/ws", service.wsEndpoint)
}

func (service *Service) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	service.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := service.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
	}
	defer ws.Close()

	player, room := service.register(ws)
	if player != nil {
		service.play(player, room, ws)
	}
}

func (service *Service) register(conn *websocket.Conn) (*models.Player, *models.Room) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
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
			return nil, nil
		}

		player, err := service.getOrCreatePlayer(req.PlayerToken, req.Username, room)
		if err != nil {
			return nil, nil
		}
		service.connections[player.Token] = conn

		// TODO write response

		return player, room
	}
	return nil, nil
}

func (service *Service) play(player *models.Player, room *models.Room, conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// TODO move the following part(parse) to messages package
		// basic format checking of the message
		messageSlice := strings.Fields(string(p))
		if messageType != 1 || len(messageSlice) != 7 || messageSlice[0] != "guess" {
			continue
		}

		// parse the message
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
		trueGuess, newCards := game.Check(guess)
		if !trueGuess {
			continue
		}

		// increase player score
		player.Score += 1
		service.savePlayer(player)

		// find members of this room
		var listeners []models.Player
		service.getRoomPeople(room, &listeners)

		// send updates to people
		for _, listener := range listeners {
			go func(places [3][2]int, values [3]int, conn *websocket.Conn) {
				if conn == nil {
					return
				}
				conn.WriteMessage(1, messages.UpdateScore{
					Username: player.Username,
					Score:    player.Score,
				}.ToBytes())
				for i := 0; i < 3; i++ {
					conn.WriteMessage(1, messages.UpdateCard{
						Row:    places[i][0],
						Column: places[i][1],
						Card:   values[i],
					}.ToBytes())
				}
			}(guess, newCards, service.connections[listener.Token])
		}
	}

}

func (service *Service) reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}
