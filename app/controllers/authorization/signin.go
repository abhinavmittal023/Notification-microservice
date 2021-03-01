package authorization

import (
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SignInRoute is used to sign in users
func SignInRoute(router *gin.RouterGroup) {
	router.POST("", SignIn)
}

// SignIn Controller for /login route
func SignIn(c *gin.Context) {
	var info serializers.LoginInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().EmailPasswordRequired})
		return
	}
	info.Email = strings.ToLower(info.Email)

	status, err := serializers.EmailRegexCheck(info.Email)

	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := users.GetUserWithEmail(info.Email)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.Errors().CredentialsMismatch})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Get user with email error", err.Error())
		return
	}

	match, err := hash.Validate(info.Password, user.Password, configuration.GetResp().PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Error while validating the password", err.Error())
		return
	}

	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.Errors().CredentialsMismatch})
		return
	}

	if !user.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.Errors().EmailNotVerified})
		return
	}

	info.ID = user.ID
	info.FirstName = user.FirstName
	info.LastName = user.LastName
	info.Password = ""
	info.Role = user.Role
	var token serializers.RefreshToken

	token.AccessToken, err = auth.GenerateAccessToken(uint64(user.ID), user.Role, configuration.GetResp().Token.ExpiryTime.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Access Token not generated", err.Error())
		return
	}
	token.RefreshToken, err = auth.GenerateRefreshToken(uint64(user.ID), configuration.GetResp().Token.ExpiryTime.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Refresh Token not generated", err.Error())
		return
	}

	loginResponse := serializers.LoginResponse{
		LoginInfo:    info,
		RefreshToken: token,
	}
	c.JSON(http.StatusOK, loginResponse)
}
