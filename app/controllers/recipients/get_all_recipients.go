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
	"github.com/jinzhu/gorm"
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
	warning := ""

	for _, recipient := range recipientArray {
		var info serializers.RecipientInfo
		serializers.RecipientModelToRecipientInfo(&info, &recipient)
		if info.PreferredChannelType != 0 {
			channel, err := channels.GetChannelWithType(uint(info.PreferredChannelType))
			if err == gorm.ErrRecordNotFound {
				// TODO: Should the PreferredChannelID field be cleared or just hidden
				// when channel is corresponding deleted
				warning = "Some Preferred Channels were Deleted"
				channel.Type = 0
				info.PreferredChannelType = 0
			} else if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			info.ChannelType = uint(channel.Type)
		}
		infoArray = append(infoArray, info)
	}

	if warning != "" {
		c.JSON(http.StatusOK, gin.H{
			"recipient_records": infoArray,
			"warning":           warning})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"recipient_records": infoArray})
	}
}
