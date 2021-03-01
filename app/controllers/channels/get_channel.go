package channels

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	miscquery "code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/misc_query"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
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

	var iteratorInfo miscquery.Iterator
	err = c.BindQuery(&iteratorInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Non boolean next or prev provided",
		})
		return
	}

	if iteratorInfo.Next && iteratorInfo.Previous {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Provide either next or previous",
		})
		return
	}

	var channel *models.Channel
	if iteratorInfo.Next {
		channel, err = channels.GetNextChannelfromID(channelID)
	} else if iteratorInfo.Previous {
		channel, err = channels.GetPreviousChannelfromID(channelID)
	} else {
		channel, err = channels.GetChannelWithID(uint(channelID))
	}

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

	firstRecord, err := channels.GetFirstChannel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	prevAvail := firstRecord.ID != channel.ID

	lastRecord, err := channels.GetLastChannel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	nextAvail := lastRecord.ID != channel.ID

	var channelInfo serializers.ChannelInfo
	serializers.ChannelModelToChannelInfo(&channelInfo, channel)

	c.JSON(http.StatusOK, gin.H{
		"channel_details": channelInfo,
		"next":            nextAvail,
		"previous":        prevAvail,
	})
}
