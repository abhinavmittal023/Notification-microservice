package services

import (
	"log"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
)

// CreateToken creates and saves the reset token into the users table
func CreateToken(user *models.User) (string, error) {

	tx := db.Get().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println(r)
			return
		}
	}()
	var err error

	resetToken := hash.GenerateSecureToken(constants.ResetTokenLength)
	user.ResetToken, err = hash.Message(resetToken, configuration.GetResp().ResetTokenHash)
	if err != nil {
		errt := tx.Rollback().Error
		log.Println("Transaction Rollback Error", errt.Error())
		return "", err
	}
	err = tx.Save(user).Error
	if err != nil {
		errt := tx.Rollback().Error
		log.Println("Transaction Rollback Error", errt.Error())
		return "", err
	}
	err = tx.Commit().Error
	if err != nil {
		log.Println("Transaction Commit Error", err.Error())
		return "", err
	}
	return resetToken, nil
}
