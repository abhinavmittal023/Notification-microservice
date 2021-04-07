package models

import "github.com/jinzhu/gorm"

// User model to create 'users' table in the database
type User struct {
	gorm.Model
	FirstName string `gorm:"not null;type:varchar(255)"`
	LastName  string `gorm:"type:varchar(255)"`
	Email     string `gorm:"type:text;not null"`
	Password  string `gorm:"type:text;not null"`
	Verified  bool
	Role      int
}
