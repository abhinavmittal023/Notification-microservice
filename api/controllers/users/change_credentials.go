package users

import (
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChangeUserCredentialsRoute is used to change users email in database
func ChangeUserCredentialsRoute(router *gin.RouterGroup) {
	router.PUT("/changecredentials", ChangeCredentials)
	router.OPTIONS("/changecredentials", preflight.Preflight)
}

// ChangeCredentials Controller for /users/changeemail route
func ChangeCredentials(c *gin.Context) {
	var info serializers.ChangeCredentialsInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	if info.Email != ""{
		info.Email = strings.ToLower(info.Email)
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
	}

	user, err := users.GetUserWithID(uint64(info.ID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID not in database"})
		return
	}
	if info.Role == 0{
		info.Role = user.Role
	}
	if info.Email == ""{
		info.Email = user.Email
	}

	serializers.ChangeCredentialsInfoToUserModel(&info, user)
	err = users.PatchUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Update User service error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
