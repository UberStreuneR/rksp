package entity

import "time"

type User struct {
	ID uint `gorm:"primarykey"`
	// Channel []*Channel `gorm:"many2many:user_channels;constraint:OnDelete:CASCADE;"`
	Channels []*Channel `gorm:"many2many:user_channels;constraint:OnDelete:CASCADE;"`
}

type Message struct {
	ID uint `gorm:"primarykey;autoIncrement"`
	// Channel []*Channel `gorm:"many2many:user_channels;constraint:OnDelete:CASCADE;"`
	ChannelID uint
	Channel   *Channel `gorm:"constraint:OnDelete:CASCADE;"`
	UserID    uint
	User      *User `gorm:"constraint:OnDelete:CASCADE;"`
	Content   string
	CreatedAt time.Time
}

type Channel struct {
	ID       uint       `gorm:"primarykey"`
	Name     string     `gorm:"unique"`
	Users    []*User    `gorm:"many2many:user_channels;constraint:OnDelete:CASCADE;"`
	Messages []*Message `gorm:"constraint:OnDelete:CASCADE;"`
}

type UserLog struct {
	ID          uint `gorm:"primarykey;autoIncrement;"`
	UserID      uint
	ChannelName string
	Operation   string
	CreatedAt   time.Time
}
