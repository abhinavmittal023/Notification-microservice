package users

import (
	"fmt"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/logs"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	li "code.jtg.tools/ayush.singhal/notifications-microservice/shared/logwrapper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DeleteUserRoute is used to delete users from database
func DeleteUserRoute(router *gin.RouterGroup) {
	router.DELETE("/:id", DeleteUser)
}

// DeleteUser Controller for delete /users/:id route
func DeleteUser(c *gin.Context) {
	f := li.GetFile()
	var err error
	var standardLogger = li.NewLogger()
	standardLogger.SetOutput(f)
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidID})
		return
	}
	user, err := users.GetUserWithID(uint64(userID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().IDNotInRecords})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		standardLogger.InternalServerError(err.Error())
		return
	}

	err = users.DeleteUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		standardLogger.InternalServerError(err.Error())
		return
	}
	logs.AddLogs(constants.InfoLog, fmt.Sprintf("User with email %s deleted", user.Email))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
