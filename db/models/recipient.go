package models

import (
	"github.com/jinzhu/gorm"
)

// Recipient model to create 'recipients' table in the database
type Recipient struct {
	gorm.Model
	Email              string
	PushToken          string
	WebToken           string
	PreferredChannel   Channel
	PreferredChannelID uint64
}

//RecipientNotifications model to create 'recipientnotifications' table in the database
type RecipientNotifications struct {
	gorm.Model
	Notification   Notification
	NotificationID uint64 `gorm:"not null"`
	Recipient      Recipient
	RecipientID    uint64 `gorm:"not null"`
	Channel        Channel
	ChannelID      uint64 `gorm:"not null;index"`
	Status         uint64 `gorm:"not null;index"`
}
