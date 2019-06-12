package service

import (
	"github.com/jinzhu/gorm"
	"net/http"
)

type Service struct {
	DB     *gorm.DB
	config Config
}

func NewService(cfg *Config) (service *Service, err error) {
	service = &Service{}
	err = nil
	if service.DB, err = gorm.Open("sqlite3", cfg.DBAddress); err != nil {
		return nil, err
	}
	service.config = *cfg
	return service, nil
}

func (service *Service) Run() (err error) {
	service.setupRoutes()
	err = http.ListenAndServe(":"+service.config.Port, nil)
	return err
}
