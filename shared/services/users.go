package services

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
)

// CreateToken creates and saves the reset token into the users table
func CreateToken(user *models.User) (string, error) {

	var err error

	resetToken := hash.GenerateSecureToken(constants.ResetTokenLength)
	user.ResetToken, err = hash.Message(resetToken, configuration.GetResp().ResetTokenHash)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}
