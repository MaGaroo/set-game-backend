package models

import "github.com/jinzhu/gorm"

type Player struct {
	gorm.Model
	Username string
	RoomID   int64
	Room     Room `gorm:"foreignkey:RoomID"`
	Score    int64
}
