package channels

import (
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
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
			"error": "ID should be a unsigned integer",
		})
		return
	}

	channel, err := channels.GetChannelWithID(uint(channelID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID is not in the database",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	var channelInfo serializers.ChannelInfo
	serializers.ChannelModelToChannelInfo(&channelInfo, channel)

	c.JSON(http.StatusOK, channelInfo)
}
