package channels

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AddChannelRoute is used to add channels to database
func AddChannelRoute(router *gin.RouterGroup) {
	router.POST("", AddChannel)
	router.OPTIONS("", preflight.Preflight)
}

// AddChannel controller for the post channels/ route
func AddChannel(c *gin.Context) {
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

	var channel models.Channel
	serializers.ChannelInfoToChannelModel(&info, &channel)

	_, err := channels.GetChannelWithName(channel.Name)
	if err != gorm.ErrRecordNotFound && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel with provided name already exists"})
		return
	}

	_, err = channels.GetChannelWithType(info.Type)
	if err != gorm.ErrRecordNotFound && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel with provided type already exists"})
		return
	}

	err = channels.AddChannel(&channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("AddChannel service error")
		return
	}

	serializers.ChannelModelToChannelInfo(&info, &channel)
	c.JSON(http.StatusOK, info)
}
