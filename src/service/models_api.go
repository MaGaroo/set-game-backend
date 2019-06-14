package service

import (
	"set-game/src/models"
	"github.com/google/uuid"
)

func (service *Service) AutoMigrate() error {
	if err := service.DB.AutoMigrate(&models.Room{}).Error; err != nil {
		return err
	}
	if err := service.DB.AutoMigrate(&models.Game{}).Error; err != nil {
		return err
	}
	if err := service.DB.AutoMigrate(&models.Player{}).Error; err != nil {
		return err
	}
	return nil
}

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

func (service *Service) getGame(room *models.Room) *models.Game {
	var game models.Game
	if err := service.DB.Where("room_id = ?", room.ID).First(&game).Error; err != nil {
		return nil
	}
	return &game
}

func (service *Service) getRoomPeople(room *models.Room, players *[]models.Player) {
	service.DB.Where("room_id = ?", room.ID).Find(players)
}

func (service *Service) savePlayer(player *models.Player) {
	service.DB.Save(player)
}
