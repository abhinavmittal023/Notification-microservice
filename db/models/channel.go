package models

import (
	"github.com/jinzhu/gorm"
)

// Channel model to create 'channels' table in the database
type Channel struct {
	gorm.Model
	Name             string `gorm:"not null;type:varchar(255)"`
	ShortDescription string
	Type             int `gorm:"not null"`
	Priority         int `gorm:"not null"`
	Configuration    string
}
