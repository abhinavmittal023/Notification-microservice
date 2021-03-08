package channels

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/misc"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UpdateChannelRoute is used to update existing channels
func UpdateChannelRoute(router *gin.RouterGroup) {
	router.PUT(":id", UpdateChannel)
}

// UpdateChannel controller for put the channels/:id route
func UpdateChannel(c *gin.Context) {
	var info serializers.ChannelInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().NameTypePriorityRequired})
		return
	}
	_, found := misc.FindInSlice(constants.ChannelIntType(), int(info.Type))
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidType})
		return
	}
	if info.Priority > constants.MaxPriority {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidPriority})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	testChannel, err := channels.GetChannelWithName(strings.ToLower(info.Name))
	if err != gorm.ErrRecordNotFound && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	} else if testChannel.ID != channel.ID && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().ChannelNamePresent})
		return
	}

	testChannel, err = channels.GetChannelWithType(info.Type)
	if err != gorm.ErrRecordNotFound && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	} else if testChannel.ID != channel.ID && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().ChannelTypePresent})
		return
	}

	err = serializers.ChannelConfigValidation(&info)

	if err != nil && err.Error() == constants.Errors().InternalError {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serializers.ChannelInfoToChannelModel(&info, channel)

	err = channels.PatchChannel(channel)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
