package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

//DeleteUserPermanently deletes the user permanently from the database, and returns error/nil
func DeleteUserPermanently(user *models.User) error {
	dbg := db.Get()
	res := dbg.Unscoped().Delete(&user)
	return res.Error
}

//SoftDeleteUser sets the deletedAt field to current time in the database, and returns error/nil
func SoftDeleteUser(user *models.User) error{
	dbg := db.Get()
	res := dbg.Delete(&user)
	return res.Error
}
