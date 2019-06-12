package models

import "github.com/jinzhu/gorm"

type Game struct {
	gorm.Model
	CardsOrder       string
	Set              string
	LastCardPosition uint8
}
