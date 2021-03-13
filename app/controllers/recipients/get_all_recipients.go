package recipients

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipients"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAllRecipientRoute is used to get recipients from database
func GetAllRecipientRoute(router *gin.RouterGroup) {
	router.GET("", GetAllRecipient)
}

// GetAllRecipient Controller for get /recipient route
func GetAllRecipient(c *gin.Context) {

	var err error

	var pagination serializers.Pagination
	err = c.BindQuery(&pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidPagination,
		})
		return
	}

	var recipientFilter filter.Recipient
	err = c.BindQuery(&recipientFilter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidFilter,
		})
		return
	}
	filter.ConvertRecipientStringToLower(&recipientFilter)

	recordsCount, err := recipients.GetAllRecipientsCount(&recipientFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("find all recipients query error", err.Error())
		return
	}

	var infoArray []serializers.RecipientInfo
	if recordsCount == 0 {
		c.JSON(http.StatusOK, gin.H{
			"records_count":     recordsCount,
			"recipient_records": infoArray})
		return
	}

	recipientArray, err := recipients.GetAllRecipients(&pagination, &recipientFilter)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	warning := ""

	for _, recipient := range recipientArray {
		info := serializers.RecipientModelToRecipientInfo(&recipient)
		if info.PreferredChannelType != 0 {
			channel, err := channels.GetChannelWithType(uint(info.PreferredChannelType))
			if err == gorm.ErrRecordNotFound {
				// TODO: Should the PreferredChannelID field be cleared or just hidden
				// when channel is corresponding deleted
				warning = "Some Preferred Channels were Deleted"
				channel.Type = 0
				info.PreferredChannelType = 0
			} else if err != nil {
				log.Println(err.Error())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
				return
			}
			channelInfo := serializers.ChannelModelToChannelInfo(channel)
			info.PreferredChannel = *channelInfo
			info.ChannelType = uint(channel.Type)
		}
		infoArray = append(infoArray, *info)
	}

	if warning != "" {
		c.JSON(http.StatusOK, gin.H{
			"records_count":     recordsCount,
			"recipient_records": infoArray,
			"warning":           warning})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"records_count":     recordsCount,
			"recipient_records": infoArray})
	}
}
