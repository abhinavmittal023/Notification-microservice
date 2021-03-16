package channels

import (
	"fmt"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
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
	f, err := li.OpenFile()
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	}
	defer f.Close()
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
		standardLogger.InternalServerError("Get Channel with ID in delete channel")
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	err = channels.DeleteChannel(channel)
	if err != nil {
		standardLogger.InternalServerError(fmt.Sprintf("Delete Channel %s from database", channel.Name))
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	standardLogger.EntityDeleted(fmt.Sprintf("channel %s", channel.Name))
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
