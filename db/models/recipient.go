package models

import (
	"github.com/jinzhu/gorm"
)

// Recipient model to create 'recipients' table in the database
type Recipient struct {
	gorm.Model
	RecipientID          string `gorm:"type:varchar(100);not null"`
	Email                string
	PushToken            string
	WebToken             string
	PreferredChannel     Channel
	PreferredChannelType uint
}

// RecipientNotifications model to create 'recipientnotifications' table in the database
type RecipientNotifications struct {
	gorm.Model
	Notification   Notification
	NotificationID uint64 `gorm:"not null"`
	Recipient      Recipient
	RecipientID    uint64 `gorm:"not null"`
	Channel        Channel
	ChannelName    string `gorm:"type:varchar(100);not null;index"`
	Status         uint64 `gorm:"not null;index"`
}
