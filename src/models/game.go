package models

import "github.com/jinzhu/gorm"

type Game struct {
	gorm.Model
	RoomID           uint
	Room             Room `gorm:"foreignkey:RoomID"`
	CardsOrder       string
	Set              string
	LastCardPosition uint8
}

func (game *Game) Check(guess [3][2]int) (bool, [3]int) {
	return false, [3]int{0, 0, 0}
}
