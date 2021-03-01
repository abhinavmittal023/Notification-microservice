package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// DeleteUser sets the deletedAt field to current time in the database, and returns error/nil
func DeleteUser(user *models.User) error {
	dbg := db.Get()
	res := dbg.Delete(&user)
	return res.Error
}
