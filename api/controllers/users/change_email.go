package users

import (
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ChangeUserEmailRoute is used to change users email in database
func ChangeUserEmailRoute(router *gin.RouterGroup) {
	router.PUT("/changeemail", ChangeEmail)
	router.OPTIONS("/changeemail", preflight.Preflight)
}

//ChangeEmail Controller for /users/changeemail route
func ChangeEmail(c *gin.Context) {
	val, _ := c.Get(constants.Role)
	if val != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var info serializers.ChangeEmailInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old Email, New Email are required"})
		return
	}
	info.OldEmail = strings.ToLower(info.OldEmail)
	info.NewEmail = strings.ToLower(info.NewEmail)

	er := serializers.EmailRegexCheck(info.OldEmail)

	if er == "internal_server_error" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Internal Server Error due to email regex")
		return
	}
	if er == "bad_request" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is invalid"})
		return
	}

	er = serializers.EmailRegexCheck(info.NewEmail)

	if er == "internal_server_error" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Internal Server Error due to email regex")
		return
	}
	if er == "bad_request" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is invalid"})
		return
	}

	user, err := users.GetUserWithEmail(info.OldEmail)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EmailId not in database"})
		return
	}

	serializers.ChangeEmailInfoToUserModel(&info, user)
	err = users.PatchUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Update User service error")
		return
	}
	to := []string{
		info.NewEmail,
	}
	err = auth.SendValidationEmail(to, uint64(user.ID))
	if err != nil {
		serializers.RevertChangesToUserEmailModel(&info, user)
		err = users.PatchUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			log.Println("Revert changes Error")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("SMTP Error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
