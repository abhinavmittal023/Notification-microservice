package users

import (
	"fmt"
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
)

// CreateUser creates a new user in the database, and returns error/nil
func CreateUser(user *models.User) error {
	return db.Get().Create(user).Error
}

// CreateUserAndVerify creates a new user and sends a verification mail
func CreateUserAndVerify(user *models.User) (int, error) {

	tx := db.Get().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println(r)
			return
		}
	}()

	err := tx.Create(user).Error
	if err != nil {
		log.Println("Create User error", err.Error())
		return http.StatusInternalServerError, fmt.Errorf(constants.Errors().InternalError)
	}
	to := []string{
		user.Email,
	}
	err = auth.SendValidationEmail(to, uint64(user.ID))
	if err != nil {
		log.Println("SMTP Error", err.Error())
		err = tx.Rollback().Error
		if err != nil {
			log.Println("Transaction Rollback Error", err.Error())
			return http.StatusInternalServerError, fmt.Errorf(constants.Errors().InternalError)
		}
		return http.StatusInternalServerError, fmt.Errorf(constants.Errors().InternalError)
	}
	err = tx.Commit().Error
	if err != nil {
		log.Println("Transaction Commit Error", err.Error())
		return http.StatusInternalServerError, fmt.Errorf(constants.Errors().InternalError)
	}
	return http.StatusOK, nil

}
