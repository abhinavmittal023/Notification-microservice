package channels

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/misc"
	"github.com/gin-gonic/gin"
)

// GetChannelRoute is used to get channel from the database
func GetChannelRoute(router *gin.RouterGroup) {
	router.GET("/available", GetAvailableChannels)
}

// GetAvailableChannels is the controller
func GetAvailableChannels(c *gin.Context) {

	var pagination = serializers.Pagination{}
	var channelFilter = filter.Channel{}
	channelList, err := channels.GetAllChannels(&pagination, &channelFilter)
	if err != nil {
		log.Println("find all channels query error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	channelTypes := constants.ChannelIntType()
	var takenChannelTypes []int
	for _, channel := range channelList {
		if len(channelTypes) == 0 {
			break
		}
		index, found := misc.FindInSlice(channelTypes, channel.Type)
		if !found {
			log.Println("Not found channel type from db in constants")
			c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
			return
		}
		takenChannelTypes = append(takenChannelTypes, channel.Type)
		channelTypes[index], channelTypes[len(channelTypes)-1] = channelTypes[len(channelTypes)-1], channelTypes[index]
		channelTypes = channelTypes[:len(channelTypes)-1]
	}
	channelsInfo, err := channels.GetChannelConfig(channelTypes)
	if err != nil {
		log.Println("JSON marshalling error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.Errors().InternalError,
		})
		return
	}
	takenChannelsInfo, err := channels.GetChannelConfig(takenChannelTypes)
	if err != nil {
		log.Println("JSON marshalling error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.Errors().InternalError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"available_channels": channelsInfo,
		"taken_channels":     takenChannelsInfo,
	})
}
