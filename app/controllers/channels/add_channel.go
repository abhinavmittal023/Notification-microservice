package channels

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, type and priority are required"})
		return
	}
	var channel models.Channel
	serializers.ChannelInfoToChannelModel(&info, &channel)

	err := channels.AddChannel(&channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("AddChannel service error")
		return
	}

	serializers.ChannelModelToChannelInfo(&info, &channel)
	c.JSON(http.StatusOK, info)
}
