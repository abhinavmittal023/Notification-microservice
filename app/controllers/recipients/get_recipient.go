package recipients

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipients"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetRecipientRoute is used to get recipients from database
func GetRecipientRoute(router *gin.RouterGroup) {
	router.GET("/:id", GetRecipient)
}

// GetRecipient Controller for get /recipient/:id route
func GetRecipient(c *gin.Context) {
	recipientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidID})
		return
	}
	recipient, err := recipients.GetRecipientWithID(uint64(recipientID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().IDNotInRecords})
		return
	} else if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	var info serializers.RecipientInfo
	serializers.RecipientModelToRecipientInfo(&info, recipient)

	if info.PreferredChannelID != 0 {
		var channelInfo serializers.ChannelInfo
		channel, err := channels.GetChannelWithID(uint(info.PreferredChannelID))
		if err == gorm.ErrRecordNotFound {
			// TODO: Should the PreferredChannelID field be cleared or just hidden
			// when channel is corresponding deleted
			channelID := info.PreferredChannelID
			info.PreferredChannelID = 0
			c.JSON(http.StatusOK, gin.H{
				"recipient_details": info,
				"warning":           fmt.Sprintf("Preferred Channel %v was Deleted", channelID),
			})
			return
		} else if err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
			return
		}
		serializers.ChannelModelToChannelInfo(&channelInfo, channel)
		c.JSON(http.StatusOK, gin.H{
			"recipient_details": info,
			"preferred_channel": channelInfo,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recipient_details": info,
	})
}
