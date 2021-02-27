package users

import (
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChangeUserCredentialsRoute is used to change users email in database
func ChangeUserCredentialsRoute(router *gin.RouterGroup) {
	router.PUT("/changecredentials", ChangeCredentials)
}

// ChangeCredentials Controller for /users/changeemail route
func ChangeCredentials(c *gin.Context) {
	var info serializers.ChangeCredentialsInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old Email, New Email are required"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "EmailId not in database"})
		return
	}
	if info.Role == 0 {
		info.Role = user.Role
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
