package authorization

import (
	"encoding/json"
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SignInRoute is used to sign in users
func SignInRoute(router *gin.RouterGroup) {
	router.POST("/", SignIn)
}

// SignIn Controller for /signin route
func SignIn(c *gin.Context) {
	var info serializers.LoginInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email,Password are required"})
		return
	}

	er := serializers.EmailRegexCheck(info.Email)

	if er == "internal_server_error" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Internal Server Error due to email regex")
		return
	}
	if er == "bad_request" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is invalid"})
		return
	}

	user,err := users.GetUserWithEmail(info.Email)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "EmailId not in database"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if !user.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "EmailId not verified"})
		return
	}

	if !hash.Validate(info.Password, user.Password, configuration.GetResp().PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Passwords mismatch"})
		return
	}

	info.FirstName = user.FirstName
	info.LastName = user.LastName
	info.Password = ""
	var token serializers.RefreshToken

	token.AccessToken, err = auth.GenerateAccessToken(uint64(user.ID), user.Role, configuration.GetResp().Token.ExpiryTime.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Access Token not generated"})
		return
	}
	token.RefreshToken, err = auth.GenerateRefreshToken(uint64(user.ID), configuration.GetResp().Token.ExpiryTime.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refresh Token not generated"})
		return
	}

	js, err := json.Marshal(&serializers.LoginResponse{
		LoginInfo:    info,
		RefreshToken: token,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON marshalling error"})
		return
	}
	c.Data(http.StatusOK, "application/json", js)
}
