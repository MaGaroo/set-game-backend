package service

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"fmt"
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

	service.reader(ws)
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
