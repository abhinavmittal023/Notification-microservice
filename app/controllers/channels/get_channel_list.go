package channels

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
)

// GetAllChannelsRoute is used to get all the channels from the database
func GetAllChannelsRoute(router *gin.RouterGroup) {
	router.GET("", GetAllChannels)
}

// GetAllChannels function is a controller for the get channels/ route
func GetAllChannels(c *gin.Context) {

	var err error
	var pagination serializers.Pagination
	err = c.BindQuery(&pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidPagination,
		})
		return
	}

	var channelFilter filter.Channel
	err = c.BindQuery(&channelFilter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidFilter,
		})
		return
	}
	filter.ConvertChannelStringToLower(&channelFilter)

	channelList, err := channels.GetAllChannels(&pagination, &channelFilter)
	if err != nil {
		log.Println("find all channels query error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	var infoArray []serializers.ChannelInfo
	var info serializers.ChannelInfo

	for _, channel := range channelList {
		serializers.ChannelModelToChannelInfo(&info, &channel)
		infoArray = append(infoArray, info)
	}
	c.JSON(http.StatusOK, infoArray)
}
