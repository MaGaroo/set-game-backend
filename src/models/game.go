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
	// table and gone cards will be empty
	game.Table = "[]"
	game.Gone = "[]"

	var cards [81]int
	for i := 0; i < 81; i++ {
		cards[i] = i
	}
	// TODO seed the rand package
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	bytesDeck, err := json.Marshal(cards)
	if err != nil {
		return err
	}
	game.Deck = string(bytesDeck)
	return nil
}

func (game *Game) Check(guess [3][2]int) (bool, []int, bool) {
	var tableCards []int
	if err := json.Unmarshal([]byte(game.Table), &tableCards); err != nil {
		return false, []int{}, false
	}
	guessPositions := make([]int, 3)
	for i := 0; i < 3; i++ {
		guessPositions[i] = guess[i][0]*3 + guess[i][1]
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
