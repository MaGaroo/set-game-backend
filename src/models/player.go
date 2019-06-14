package models

import "github.com/jinzhu/gorm"

type Player struct {
	gorm.Model
	Username string
	RoomID   uint
	Room     Room `gorm:"foreignkey:RoomID"`
	Score    int64
	Token    string
}
