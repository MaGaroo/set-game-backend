package service

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"fmt"
	"set-game/src/models"
	"set-game/src/service/messages"
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

	player := service.register(ws)
	if player != nil {
		service.play(player)
	}
}

func (service *Service) register(conn *websocket.Conn) *models.Player {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return nil
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
			return nil
		}

		player, err := service.getOrCreatePlayer(req.PlayerToken, req.Username, room)
		if err != nil {
			return nil
		}
		service.connections[player.Token] = conn

		return player
	}
	return nil
}

func (service *Service) play(player *models.Player) {

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
