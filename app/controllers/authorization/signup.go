package authorization

import (
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
)

// SignUpRoute is used to sign up users
func SignUpRoute(router *gin.RouterGroup) {
	router.POST("", SignUp)
}

// SignUp Controller for /signup route
func SignUp(c *gin.Context) {
	var info serializers.SignupInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().EmailPasswordNameRequired})
		return
	}
	info.Email = strings.ToLower(info.Email)
	info.Role = constants.SystemAdminRole // signup user will always be system admin

	status, err := serializers.EmailRegexCheck(info.Email)

	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	info.Password, err = hash.Message(info.Password, configuration.GetResp().PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Error while hashing the password", err.Error())
		return
	}
	var user models.User

	serializers.SignupInfoToUserModel(&info, &user)
	status, err = users.CreateUserAndVerify(&user)
	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
