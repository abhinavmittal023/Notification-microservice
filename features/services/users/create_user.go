package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

//CreateUser creates a new user in the database, and returns error/nil
func CreateUser(user *models.User) error {
	dbg := db.Get()
	res := dbg.Create(user)
	return res.Error
}
