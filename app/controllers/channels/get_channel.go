package channels

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetChannelRoute is used to get a channel from the database given its id
func GetChannelRoute(router *gin.RouterGroup) {
	router.GET(":id", GetChannel)
}

// GetChannel function is a controller for get channels/:id route
func GetChannel(c *gin.Context) {

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
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	var channelInfo serializers.ChannelInfo
	serializers.ChannelModelToChannelInfo(&channelInfo, channel)

	c.JSON(http.StatusOK, channelInfo)
}
