package recipients

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	miscquery "code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/misc_query"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipients"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
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

	var iteratorInfo miscquery.Iterator
	err = c.BindQuery(&iteratorInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().NextPrevNonBool,
		})
		return
	}

	if iteratorInfo.Next && iteratorInfo.Previous {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().NextPrevBothProvided,
		})
		return
	}

	var recipient *models.Recipient
	if iteratorInfo.Next {
		recipient, err = recipients.GetNextRecipientfromID(uint64(recipientID))
	} else if iteratorInfo.Previous {
		recipient, err = recipients.GetPreviousRecipientfromID(uint64(recipientID))
	} else {
		recipient, err = recipients.GetRecipientWithID(uint64(recipientID))
	}

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().IDNotInRecords})
		return
	} else if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	firstRecord, err := recipients.GetFirstRecipient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	prevAvail := firstRecord.ID != recipient.ID

	lastRecord, err := recipients.GetLastRecipient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	nextAvail := lastRecord.ID != recipient.ID

	var info serializers.RecipientInfo
	serializers.RecipientModelToRecipientInfo(&info, recipient)

	if info.PreferredChannelType != 0 {
		var channelInfo serializers.ChannelInfo
		channel, err := channels.GetChannelWithType(uint(info.PreferredChannelType))
		if err == gorm.ErrRecordNotFound {
			// TODO: Should the PreferredChannelID field be cleared or just hidden
			// when channel is corresponding deleted
			channelType := info.PreferredChannelType
			info.PreferredChannelType = 0
			c.JSON(http.StatusOK, gin.H{
				"recipient_details": info,
				"next":              nextAvail,
				"previous":          prevAvail,
				"warning":           fmt.Sprintf("Preferred Channel %s was Deleted", constants.ChannelType(channelType)),
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
			"next":              nextAvail,
			"previous":          prevAvail,
			"preferred_channel": channelInfo,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recipient_details": info,
		"next":              nextAvail,
		"previous":          prevAvail,
	})
}
