package recipients

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipients"
	"github.com/gin-gonic/gin"
)

// GetAllRecipientRoute is used to get recipients from database
func GetAllRecipientRoute(router *gin.RouterGroup) {
	router.GET("", GetAllRecipient)
	router.OPTIONS("", preflight.Preflight)
}

// GetAllRecipient Controller for get /recipient route
func GetAllRecipient(c *gin.Context) {

	var err error

	var pagination serializers.Pagination
	err = c.BindQuery(&pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit and offset",
		})
		return
	}

	var recipientFilter filter.Recipient
	err = c.BindQuery(&recipientFilter)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid filter Parameters",
		})
		return
	}

	recipientArray, err := recipients.GetAllRecipients(pagination, recipientFilter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	var infoArray []serializers.RecipientInfo

	for _, recipient := range recipientArray {
		var info serializers.RecipientInfo
		serializers.RecipientModelToRecipientInfo(&info, &recipient)
		if info.PreferredChannelID != 0 {
			channel, err := channels.GetChannelWithID(uint(info.PreferredChannelID))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			info.ChannelType = uint(channel.Type)
		}
		infoArray = append(infoArray, info)
	}

	c.JSON(http.StatusOK, infoArray)
}
