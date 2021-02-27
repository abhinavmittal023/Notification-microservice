package authorization

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email,Password are required"})
		return
	}
	info.Email = strings.ToLower(info.Email)

	status, message := serializers.EmailRegexCheck(info.Email)

	if status != http.StatusOK {
		c.JSON(status, gin.H{
			"error": message,
		})
		return
	}

	user, err := users.GetUserWithEmail(info.Email)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "EmailId or Passwords mismatch"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Get user with email error")
		return
	}

	match, err := hash.Validate(info.Password, user.Password, configuration.GetResp().PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Error while validating the password")
		return
	}

	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "EmailId or Passwords mismatch"})
		return
	}

	if !user.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "EmailId not verified"})
		return
	}

	info.FirstName = user.FirstName
	info.LastName = user.LastName
	info.Password = ""
	info.Role = user.Role
	var token serializers.RefreshToken

	token.AccessToken, err = auth.GenerateAccessToken(uint64(user.ID), user.Role, configuration.GetResp().Token.ExpiryTime.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Access Token not generated")
		return
	}
	token.RefreshToken, err = auth.GenerateRefreshToken(uint64(user.ID), configuration.GetResp().Token.ExpiryTime.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Refresh Token not generated")
		return
	}

	js, err := json.Marshal(&serializers.LoginResponse{
		LoginInfo:    info,
		RefreshToken: token,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("JSON marshalling error")
		return
	}
	c.Data(http.StatusOK, "application/json", js)
}
