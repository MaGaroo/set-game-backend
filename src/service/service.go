package service

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

type Service struct {
	DB          *gorm.DB
	config      Config
	upgrader    websocket.Upgrader
	connections map[string]*websocket.Conn
}

func NewService(cfg *Config) (service *Service, err error) {
	log.Println("Creating a new service")
	service = &Service{}
	err = nil

	log.Printf("Openning database with address %s\n", cfg.DBAddress)
	if service.DB, err = gorm.Open("sqlite3", cfg.DBAddress); err != nil {
		return nil, err
	}

	service.config = *cfg

	service.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	if err = service.AutoMigrate(); err != nil {
		log.Fatalf("Migration failed with error %s\n", err.Error())
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
