package authservice

import (
	"errors"
	"log"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
)

// ErrInvalidToken is generated if token is invalid
var ErrInvalidToken = errors.New("Invalid Token")

// ErrAlreadyVerfied is generated if user clicks on validation link after being verification
var ErrAlreadyVerfied = errors.New("User Already Verified")

// ValidateToken function is used to check if token is valid, and return user model, if valid
func ValidateToken(tokenString string, tokenType string) (*models.User, error) {
	token, err := auth.ValidateToken(tokenString)
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims := token.Claims.(*auth.CustomClaims)

	if claims.TokenType == tokenType {

		userDetails, err := users.GetUserWithID(claims.UserID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if tokenType == "validation" && userDetails.Verified == true {
			return userDetails, ErrAlreadyVerfied
		}
		if token.Valid {
			return userDetails, nil
		}
		return nil, ErrInvalidToken
	}
	return nil, ErrInvalidToken
}
