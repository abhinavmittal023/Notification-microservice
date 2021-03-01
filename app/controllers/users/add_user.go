package users

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/misc"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AddUserRoute is used to add users to database
func AddUserRoute(router *gin.RouterGroup) {
	router.POST("", AddUser)
}

// AddUser Controller for post /users/ route
func AddUser(c *gin.Context) {
	var info serializers.AddUserInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Role, %s", constants.Errors().EmailPasswordNameRequired)})
		return
	}
	_, found := misc.FindInSlice(constants.RoleType(), int(info.Role))
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidRole})
		return
	}

	info.Email = strings.ToLower(info.Email)

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
		log.Println("Error while hashing the password")
		return
	}

	user, err := users.GetUserWithEmail(info.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().EmailAlreadyPresent})
		return
	}
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("GetUserWithEmail service error")
		return
	}

	serializers.AddUserInfoToUserModel(&info, user)
	status, err = users.CreateUserAndVerify(user)
	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
