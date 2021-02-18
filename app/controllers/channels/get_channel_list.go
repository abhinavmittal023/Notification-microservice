package channels

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"github.com/gin-gonic/gin"
)

// GetAllChannelsRoute is used to get all the channels from the database
func GetAllChannelsRoute(router *gin.RouterGroup) {
	router.GET("", GetAllChannels)
}

// GetAllChannels function is a controller for the get channels/ route
func GetAllChannels(c *gin.Context) {

	var limit uint64
	var err error
	var offset uint64
	limitString := c.Query("limit")
	offsetString := c.Query("offset")

	if limitString != "" {
		limit, err = strconv.ParseUint(limitString, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "limit should be a unsigned integer",
			})
			return
		}
	} else {
		limit = 20
	}

	if offsetString != "" {
		offset, err = strconv.ParseUint(offsetString, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "offset should be a unsigned integer",
			})
			return
		}
	} else {
		offset = 0
	}

	channelList, err := channels.GetAllChannels(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("find all channels query error")
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
