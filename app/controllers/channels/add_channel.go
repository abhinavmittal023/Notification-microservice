package channels

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/misc"
	"github.com/gin-gonic/gin"
)

// AddChannelRoute is used to add channels to database
func AddChannelRoute(router *gin.RouterGroup) {
	router.POST("", AddChannel)
}

// AddChannel controller for the post channels/ route
func AddChannel(c *gin.Context) {
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

	var channel models.Channel
	serializers.ChannelInfoToChannelModel(&info, &channel)

	err := channels.AddChannel(&channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("AddChannel service error")
		return
	}

	serializers.ChannelModelToChannelInfo(&info, &channel)
	c.JSON(http.StatusOK, info)
}
