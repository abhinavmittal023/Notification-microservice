package authorization

import (
	"fmt"
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
	"github.com/go-playground/validator/v10"
)

// SignUpRoute is used to sign up users
func SignUpRoute(router *gin.RouterGroup) {
	router.POST("", SignUp)
}

// SignUp Controller for /signup route
func SignUp(c *gin.Context) {
	var info serializers.SignupInfo
	if err := c.BindJSON(&info); err != nil {
		ve, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidDataType})
			return
		}
		var errors []string
		for _, value := range ve {
			if value.Tag() != "max" {
				errors = append(errors, fmt.Sprintf("%s is %s", value.Field(), value.Tag()))
			} else {
				errors = append(errors, fmt.Sprintf("%s cannot have more than %s characters", value.Field(), value.Param()))
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errors})
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
	status, err = users.CreateUserAndVerify(&user, true)
	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
