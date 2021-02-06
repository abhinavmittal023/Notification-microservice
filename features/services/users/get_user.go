package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

//GetUserWithID gets the user with specified ID from the database, and returns error/nil
func GetUserWithID(user *models.User, userID uint64) error {
	dbg := db.Get()
	res := dbg.First(user, userID)
	return res.Error
}

//GetFirstUser gets the details of the first user in the database
func GetFirstUser(user *models.User) error {
	dbg := db.Get()
	res := dbg.First(user)
	return res.Error
}

//GetUserWithEmail gets the user with specified Email from the database
func GetUserWithEmail(user *models.User,email string) error {
	dbg := db.Get()
	res := dbg.Where("email = ?", email).First(user)
	return res.Error
}