package users

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	li "code.jtg.tools/ayush.singhal/notifications-microservice/shared/logwrapper"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/misc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

// UpdateUserRoute is used to change users email in database
func UpdateUserRoute(router *gin.RouterGroup) {
	router.PUT("/:id/update", UpdateUser)
}

// UpdateUser Controller for put /users/:id/update route
func UpdateUser(c *gin.Context) {
	f, err := li.OpenFile()
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	}
	defer f.Close()
	var standardLogger = li.NewLogger()
	standardLogger.SetOutput(f)
	var info serializers.ChangeCredentialsInfo
	if err = c.BindJSON(&info); err != nil {
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
	_, found := misc.FindInSlice(constants.RoleType(), info.Role)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidRole})
		return
	}

	info.Email = strings.ToLower(info.Email)

	info.ID, err = strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidID,
		})
		return
	}

	status, err := serializers.EmailRegexCheck(info.Email)

	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	testUser, err := users.GetUserWithEmail(info.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		standardLogger.InternalServerError("Get User with email in update user")
		return
	} else if testUser.ID != uint(info.ID) && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().EmailAlreadyPresent,
		})
	}

	user, err := users.GetUserWithID(uint64(info.ID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().IDNotInRecords})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		standardLogger.InternalServerError("Get User with id in update user")
		return
	}

	serializers.ChangeCredentialsInfoToUserModel(&info, user)
	err = users.PatchUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		standardLogger.InternalServerError("Patch User to database")
		return
	}
	standardLogger.EntityUpdated(fmt.Sprintf("user with email %s", user.Email))

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
