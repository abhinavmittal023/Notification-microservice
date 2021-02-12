package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// GetUserWithID gets the user with specified ID from the database, and returns error/nil
func GetUserWithID(userID uint64) (*models.User, error) {
	var user models.User
	res := db.Get().First(&user, userID)
	return &user, res.Error
}

// GetFirstUser gets the details of the first user in the database
func GetFirstUser() (*models.User, error) {
	var user models.User
	res := db.Get().First(&user)
	return &user, res.Error
}

// GetUserWithEmail gets the user with specified Email from the database
func GetUserWithEmail(email string) (*models.User, error) {
	var user models.User
	res := db.Get().Where("email = ?", email).First(&user)
	return &user, res.Error
}

//GetAllUsers gets all users from the database and returns []models.User,err
func GetAllUsers() ([]models.User,error){
	var users []models.User
	dbg := db.Get()
	res := dbg.Find(&users)
	return users,res.Error
}