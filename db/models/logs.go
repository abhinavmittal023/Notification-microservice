package models

import "github.com/jinzhu/gorm"

type Logs struct{
	gorm.Model
	Level	uint
	Msg		string
}