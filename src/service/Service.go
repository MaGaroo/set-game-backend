package service

import "github.com/jinzhu/gorm"

type Service struct {
	DB *gorm.DB
}

func NewService(cfg *Config) (service *Service, err error) {
	if service.DB, err = gorm.Open("sqlite3", cfg.DBAddress); err != nil {
		return nil, err
	}
	return service, nil
}
