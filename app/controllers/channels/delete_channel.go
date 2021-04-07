package channels

import (
	"fmt"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/logs"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	li "code.jtg.tools/ayush.singhal/notifications-microservice/shared/logwrapper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DeleteChannelRoute is used to delete existing channels
func DeleteChannelRoute(router *gin.RouterGroup) {
	router.DELETE(":id", DeleteChannel)
}

// DeleteChannel function is a controller for delete channels/:id route
func DeleteChannel(c *gin.Context) {
	f := li.GetFile()
	var err error
	var standardLogger = li.NewLogger()
	standardLogger.SetOutput(f)
	channelID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidID,
		})
		return
	}

	channel, err := channels.GetChannelWithID(uint(channelID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().IDNotInRecords,
		})
		return
	} else if err != nil {
		standardLogger.InternalServerError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	err = channels.DeleteChannel(channel)
	if err != nil {
		standardLogger.InternalServerError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	logs.AddLogs(constants.InfoLog, fmt.Sprintf("Channel %s deleted", channel.Name))
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
