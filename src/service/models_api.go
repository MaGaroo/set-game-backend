package service

import (
	"set-game/src/models"
	"github.com/google/uuid"
)

func (service *Service) getOrCreatePlayer(playerToken string, username string, room *models.Room) (*models.Player, error) {
	if player, err := service.getPlayer(playerToken); err != nil {
		player, err = service.createPlayer(username, room)
		if err != nil {
			return nil, err
		}
		return player, nil
	} else {
		player.Username = username
		service.DB.Save(&player)
		return player, nil
	}
}

func (service *Service) getRoom(roomToken string) (*models.Room, error) {
	var room models.Room
	if err := service.DB.Where("token = ?", roomToken).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (service *Service) getPlayer(playerToken string) (*models.Player, error) {
	var player models.Player
	if err := service.DB.Where("token = ?", playerToken).First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (service *Service) createPlayer(username string, room *models.Room) (*models.Player, error) {
	player := models.Player{
		Username: username,
		RoomID:   room.ID,
		Token:    uuid.New().String(),
		Score:    0,
	}
	if err := service.DB.Create(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}
