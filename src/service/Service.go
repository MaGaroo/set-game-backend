package service

import "github.com/jinzhu/gorm"

type Service struct {
	DB     *gorm.DB
	config Config
}

func NewService(cfg *Config) (service *Service, err error) {
	if service.DB, err = gorm.Open("sqlite3", cfg.DBAddress); err != nil {
		return nil, err
	}
	service.config = *cfg
	return service, nil
}
