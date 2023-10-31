package services

import (
	"practice04/entity"

	"gorm.io/gorm"
)

type MessageService struct {
	DB *gorm.DB
}

func CreateMessageService(db *gorm.DB) MessageService {
	return MessageService{db}
}

func (m MessageService) GetAll() ([]*entity.Message, error) {
	var messages []*entity.Message
	results := m.DB.Find(&messages)
	if results.Error != nil {
		return messages, results.Error
	}
	return messages, nil
}

func (m MessageService) GetAllFromChannel(channelName string) ([]*entity.Message, error) {
	var messages []*entity.Message
	c := CreateChannelService(m.DB)
	channel, err := c.GetOne(channelName)
	if err != nil {
		return messages, err
	}
	results := m.DB.Preload("Messages").Find(&channel)
	messages = channel.Messages
	if results.Error != nil {
		return messages, results.Error
	}
	return messages, nil
}

func (m MessageService) AddOne(channel_name string, user_id uint, content string) (*entity.Message, error) {
	cs := CreateChannelService(m.DB)
	channel, err := cs.GetOne(channel_name)
	if err != nil {
		return nil, err
	}
	us := CreateUserService(m.DB)
	user, err := us.GetOne(user_id)
	if err != nil {
		return nil, err
	}
	message := &entity.Message{Channel: channel, User: user, Content: content}
	result := m.DB.Create(message)
	if result.Error != nil {
		return nil, result.Error
	}
	return message, nil
}

func (m MessageService) DeleteOne(id uint) error {
	result := m.DB.Delete(&entity.Message{}, "ID = ?", id)
	return result.Error
}

func (m MessageService) DeleteAll() error {
	result := m.DB.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.Message{})
	return result.Error
}
