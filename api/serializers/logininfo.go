package serializers

import (
	"net/http"
	"regexp"

	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// LoginInfo serializer to bind request data
type LoginInfo struct {
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

// LoginInfoToUserModel converts LoginInfo serializer to User model
func LoginInfoToUserModel(info LoginInfo, user *models.User) {
	user.Email = info.Email
	user.Password = info.Password
}

// EmailRegexCheck checks for email id in valid format
func EmailRegexCheck(email string) (int, string) {
	match, err := regexp.MatchString(constants.EmailRegex, email)
	if err != nil {
		return http.StatusInternalServerError, "Internal Server Error"
	}
	if !match {
		return http.StatusBadRequest, "Email ID is invalid"
	}
	return http.StatusOK, ""
}
