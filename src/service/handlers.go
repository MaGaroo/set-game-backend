package service

import "net/http"

func (service *Service) setupRoutes() {
	http.HandleFunc("/ws", service.wsEndpoint)
}

func (service *Service) wsEndpoint(w http.ResponseWriter, r *http.Request) {

}
