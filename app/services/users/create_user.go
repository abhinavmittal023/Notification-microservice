package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// CreateUser creates a new user in the database, and returns error/nil
func CreateUser(user *models.User) error {
	return db.Get().Create(user).Error
}
