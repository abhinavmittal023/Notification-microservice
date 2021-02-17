package users

import (
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AddUserRoute is used to add users to database
func AddUserRoute(router *gin.RouterGroup) {
	router.POST("/", AddUser)
	router.OPTIONS("/", preflight.Preflight)
}

// AddUser Controller for post /users/ route
func AddUser(c *gin.Context) {
	var info serializers.AddUserInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email,Password,FirstName are required"})
		return
	}
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

	info.Password = hash.Message(info.Password, configuration.GetResp().PasswordHash)

	user, err := users.GetUserWithEmail(info.Email)
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EmailId already in database"})
		return
	}

	serializers.AddUserInfoToUserModel(&info, user)
	err = users.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("CreateUser service error")
		return
	}
	to := []string{
		info.Email,
	}
	err = auth.SendValidationEmail(to, uint64(user.ID))
	if err != nil {
		err = users.DeleteUserPermanently(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			log.Println("Delete User Service Error")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("SMTP Error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
