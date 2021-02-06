package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

//PatchUser just saves the results made to user table
func PatchUser(user *models.User) error {
	dbg := db.Get()
	res := dbg.Save(user)
	return res.Error
}
