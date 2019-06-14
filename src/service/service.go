package service

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"github.com/gorilla/websocket"
)

type Service struct {
	DB          *gorm.DB
	config      Config
	upgrader    websocket.Upgrader
	connections map[string]*websocket.Conn
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

	if err = service.AutoMigrate(); err != nil {
		service.DB.Close()
		return nil, err
	}

	service.connections = make(map[string]*websocket.Conn)

	return service, nil
}

func (service *Service) Run() (err error) {
	service.setupRoutes()
	err = http.ListenAndServe(":"+service.config.Port, nil)
	return err
}

func (service *Service) AutoMigrate() error {
	return nil
}
