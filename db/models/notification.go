package models

import (
	"github.com/jinzhu/gorm"
)

// Notification model to create 'notifications' table in the database
type Notification struct {
	gorm.Model
	Priority int    `gorm:"not null"`
	Title    string `gorm:"not null"`
	Body     string `gorm:"not null"`
}
