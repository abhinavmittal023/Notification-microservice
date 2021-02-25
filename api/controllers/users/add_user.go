package users

import (
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AddUserRoute is used to add users to database
func AddUserRoute(router *gin.RouterGroup) {
	router.POST("/add", AddUser)
}

// AddUser Controller for /users/add route
func AddUser(c *gin.Context) {
	val, _ := c.Get(constants.Role)
	if val != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var info serializers.AddUserInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, Password, FirstName are required"})
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

	var err error

	info.Password, err = hash.Message(info.Password, configuration.GetResp().PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Error while hashing the password")
		return
	}

	user, err := users.GetUserWithEmail(info.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EmailId already in database"})
		return
	}
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("GetUserWithEmail service error")
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
