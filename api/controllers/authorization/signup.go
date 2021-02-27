package authorization

import (
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
)

// SignUpRoute is used to sign up users
func SignUpRoute(router *gin.RouterGroup) {
	router.POST("/", SignUp)
}

// SignUp Controller for /signup route
func SignUp(c *gin.Context) {
	var info serializers.SignupInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email,Password,FirstName are required"})
		return
	}
	info.Email = strings.ToLower(info.Email)
	info.Role = constants.SystemAdminRole // signup user will always be system admin

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

	var user models.User

	serializers.SignupInfoToUserModel(&info, &user)
	err = users.CreateUser(&user)
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
		err = users.DeleteUserPermanently(&user)
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
