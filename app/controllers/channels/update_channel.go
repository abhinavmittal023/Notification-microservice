package channels

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UpdateChannelRoute is used to update existing channels
func UpdateChannelRoute(router *gin.RouterGroup) {
	router.PUT(":id", UpdateChannel)
	router.OPTIONS(":id", preflight.Preflight)
}

// UpdateChannel controller for put the channels/:id route
func UpdateChannel(c *gin.Context) {
	var info serializers.ChannelInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "valid name, type and priority are required"})
		return
	}
	if info.Type > constants.MaxType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Type provided"})
		return
	}
	if info.Priority > constants.MaxPriority {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Priority provided"})
		return
	}

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

	testChannel, err := channels.GetChannelWithName(strings.ToLower(info.Name))
	if err != gorm.ErrRecordNotFound && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	} else if testChannel.ID != channel.ID && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel with provided name already exists"})
		return
	}

	testChannel, err = channels.GetChannelWithType(info.Type)
	if err != gorm.ErrRecordNotFound && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	} else if testChannel.ID != channel.ID && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel with provided type already exists"})
		return
	}

	serializers.ChannelInfoToChannelModel(&info, channel)

	err = channels.PatchChannel(channel)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
