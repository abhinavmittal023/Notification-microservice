package models

import "github.com/jinzhu/gorm"

// Logs is a model used to create the table for Logs
type Logs struct {
	gorm.Model
	Level uint
	Msg   string
}
