package models

import (
	"github.com/jinzhu/gorm"
	"math/rand"
	"encoding/json"
	"set-game/src/setgame"
)

type Game struct {
	gorm.Model
	RoomID uint
	Room   Room `gorm:"foreignkey:RoomID"`
	Deck   string
	Table  string
	Gone   string
}

func (game *Game) Start() error {
	// table and gone deck will be empty
	game.Table = "[]"
	game.Gone = "[]"

	deck := make([]int, 81)
	gone := make([]int, 0)
	table := make([]int, 0)
	for i := 0; i < 81; i++ {
		deck[i] = i
	}
	// TODO seed the rand package
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	table, gone, deck, _ = setgame.Normalize(table, gone, deck)

	bytesDeck, err := json.Marshal(deck)
	if err != nil {
		return err
	}
	game.Deck = string(bytesDeck)

	bytesTable, err := json.Marshal(table)
	if err != nil {
		return err
	}
	game.Table = string(bytesTable)

	bytesGone, err := json.Marshal(gone)
	if err != nil {
		return err
	}
	game.Gone = string(bytesGone)

	return nil
}

func (game *Game) Check(guessPositions []int) (bool, []int, bool) {
	var tableCards []int
	if err := json.Unmarshal([]byte(game.Table), &tableCards); err != nil {
		return false, []int{}, false
	}
	if !setgame.PositionsMakeSet(tableCards, guessPositions) {
		return false, []int{}, false
	}

	var goneCards []int
	var deckCards []int

	if err := json.Unmarshal([]byte(game.Deck), &deckCards); err != nil {
		return false, []int{}, false
	}
	if err := json.Unmarshal([]byte(game.Gone), &goneCards); err != nil {
		return false, []int{}, false
	}

	tableCards, goneCards, deckCards, endGame := setgame.RemoveAndNormalize(tableCards, goneCards, deckCards, guessPositions)
	byteTable, err := json.Marshal(tableCards)
	if err != nil {
		return false, []int{}, false
	}
	byteGone, err := json.Marshal(goneCards)
	if err != nil {
		return false, []int{}, false
	}
	byteDeck, err := json.Marshal(deckCards)
	if err != nil {
		return false, []int{}, false
	}
	game.Table = string(byteTable)
	game.Gone = string(byteGone)
	game.Deck = string(byteDeck)
	return true, tableCards, endGame
}
