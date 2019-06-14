package models

import "github.com/jinzhu/gorm"

type Player struct {
	gorm.Model
	Username string
	RoomID   uint
	Room     Room `gorm:"foreignkey:RoomID"`
	Score    int8
	Token    string
}
