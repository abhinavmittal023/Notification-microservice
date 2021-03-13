package serializers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
)

// LoginInfo serializer to bind request data
type LoginInfo struct {
	ID        uint   `json:"user_id"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password,omitempty" binding:"required"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      int    `json:"role"`
}

// LoginResponse formats the userDetails and TokenDetails into one struct
type LoginResponse struct {
	LoginInfo    LoginInfo    `json:"user_details"`
	RefreshToken RefreshToken `json:"token_details"`
}

// EmailRegexCheck checks for email id in valid format
func EmailRegexCheck(email string) (int, error) {
	match, err := regexp.MatchString(constants.EmailRegex, email)
	if err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, fmt.Errorf(constants.Errors().InternalError)
	}
	if !match {
		return http.StatusBadRequest, fmt.Errorf(constants.Errors().InvalidEmail)
	}
	return http.StatusOK, nil
}
