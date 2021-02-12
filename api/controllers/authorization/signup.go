package authorization

import (
	"log"
	"net/http"

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
	info.Role = constants.SystemAdminRole // signup user will always be system admin

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

	var user models.User

	serializers.SignupInfoToUserModel(info, &user)
	err := users.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not created Error"})
		return
	}
	to := []string{
		info.Email,
	}
	err = auth.SendValidationEmail(to, uint64(user.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email couldn't be sent"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
