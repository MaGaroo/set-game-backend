package service

import (
	"github.com/jinzhu/gorm"
	"net/http"
	"github.com/gorilla/websocket"
)

type Service struct {
	DB       *gorm.DB
	config   Config
	upgrader websocket.Upgrader
}

func NewService(cfg *Config) (service *Service, err error) {
	service = &Service{}
	err = nil
	if service.DB, err = gorm.Open("sqlite3", cfg.DBAddress); err != nil {
		return nil, err
	}
	service.config = *cfg
	service.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return service, nil
}

func (service *Service) Run() (err error) {
	service.setupRoutes()
	err = http.ListenAndServe(":"+service.config.Port, nil)
	return err
}
