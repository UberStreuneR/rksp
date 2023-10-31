package services

import (
	"fmt"
	"practice04/entity"

	"gorm.io/gorm"
)

var Channel ChannelService

type ChannelService struct {
	DB *gorm.DB
}

func CreateChannelService(db *gorm.DB) ChannelService {
	return ChannelService{db}
}

func (s ChannelService) GetAll() ([]*entity.Channel, error) {
	var channel []*entity.Channel
	results := s.DB.Preload("Users").Find(&channel)
	if results.Error != nil {
		return channel, results.Error
	}
	return channel, nil
}

func (s ChannelService) GetOne(name string) (*entity.Channel, error) {
	var channel *entity.Channel
	result := s.DB.Preload("Users").Preload("Messages").First(&channel, "name = ?", name)
	if result.Error != nil {
		return channel, result.Error
	}
	return channel, nil
}

func (s ChannelService) GetChannelsForUser(id uint) ([]*entity.Channel, error) {
	var user *entity.User
	result := s.DB.Preload("Channel").Find(&user, "id = ?", fmt.Sprint(id))
	if result.Error != nil {
		return nil, result.Error
	}
	return user.Channels, nil
}

func (s ChannelService) AddOne(name string) (*entity.Channel, error) {
	channel := &entity.Channel{Name: name}
	result := s.DB.Create(channel)
	if result.Error != nil {
		return nil, result.Error
	}
	return channel, nil
}

func (s ChannelService) UpdateOne(name, newName string) (*entity.Channel, error) {
	channel, err := s.GetOne(name)
	if err != nil {
		return nil, err
	}
	channel.Name = newName
	result := s.DB.Model(&channel).Where("name = ?", name).Update("name", newName)
	if result.Error != nil {
		return nil, result.Error
	}
	return channel, nil
}

func (s ChannelService) DeleteOne(name string) error {
	result := s.DB.Delete(&entity.Channel{}, "name = ?", name)
	return result.Error
}

func (s ChannelService) DeleteAll() error {
	result := s.DB.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.Channel{})
	return result.Error
}

func (s ChannelService) AddChannelToUser(id uint, strChannel []string) error {
	var channel []*entity.Channel
	result := s.DB.Model(&entity.Channel{}).Where("name IN (?)", strChannel).Find(&channel)
	if result.Error != nil {
		return result.Error
	}
	var user *entity.User
	result = s.DB.First(&user, id)
	if result.Error != nil {
		return result.Error
	}
	user.Channels = append(user.Channels, channel...)
	result = s.DB.Save(user)
	s.DB.Save(user.Channels)
	ul := CreateUserLogService(s.DB)
	if result.Error != nil {
		for _, channelName := range strChannel {
			ul.AddOne(user.ID, channelName, "added")
		}
	}
	return result.Error
}

func (s ChannelService) RemoveChannelFromUser(id uint, strChannel []string) error {
	var user *entity.User
	result := s.DB.First(&user, id)
	if result.Error != nil {
		return result.Error
	}
	channelHash := make(map[string]bool)
	for _, seg := range strChannel {
		channelHash[seg] = true
	}
	channels, err := s.GetChannelsForUser(id)
	if err != nil {
		return err
	}
	for _, seg := range channels {
		if channelHash[seg.Name] {
			s.DB.Model(user).Association("Channel").Delete(seg)
		}
	}
	return nil
}

func (s ChannelService) RemoveUserFromChannel(id uint, channelName string) error {
	channel, err := s.GetOne(channelName)
	if err != nil {
		return err
	}
	var user *entity.User
	result := s.DB.Preload("Channels").Find(&user, "id = ?", fmt.Sprint(id))
	if result.Error != nil {
		return result.Error
	}
	if len(user.Channels) == 0 {
		return nil
	}
	ul := CreateUserLogService(s.DB)
	err = s.DB.Model(user).Association("Channels").Delete(channel)
	if err != nil {
		ul.AddOne(id, channelName, "deleted")
	}
	return err
}
