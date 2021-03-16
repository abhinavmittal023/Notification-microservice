package users

import (
	"fmt"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/logs"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	li "code.jtg.tools/ayush.singhal/notifications-microservice/shared/logwrapper"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/misc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

// AddUserRoute is used to add users to database
func AddUserRoute(router *gin.RouterGroup) {
	router.POST("", AddUser)
}

// AddUser Controller for post /users/ route
func AddUser(c *gin.Context) {
	f, err := li.OpenFile()
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	}
	defer f.Close()
	var standardLogger = li.NewLogger()
	standardLogger.SetOutput(f)
	var info serializers.AddUserInfo
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

	user, err := users.GetUserWithEmail(info.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().EmailAlreadyPresent})
		return
	}
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	user = serializers.AddUserInfoToUserModel(&info)
	status, err = users.CreateUserAndVerify(user, false)
	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	logs.AddLogs(constants.InfoLog,fmt.Sprintf("User with email %s added",user.Email))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
