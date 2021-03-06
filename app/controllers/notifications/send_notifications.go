package notifications

import (
	"log"
	"net/http"
	"sync"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	apimessage "code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/api_message"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/notifications"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipients"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// PostSendNotificationsRoute is used to make new API Key
func PostSendNotificationsRoute(router *gin.RouterGroup) {
	router.POST("", PostSendNotifications)
}

// PostSendNotifications controller is used to send notifications from notifications route
func PostSendNotifications(c *gin.Context) {
	var info serializers.SendNotifications
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().NotificationInfoRequired})
		return
	}
	var notification models.Notification
	err := serializers.NotificationsInfoToNotificationModel(&info, &notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = notifications.AddNotification(&notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.Errors().InternalError,
		})
		return
	}
	openAPI := apimessage.OpenAPI{
		NotificationStatus:          make(map[string]apimessage.OpenAPIChannel),
		RecipientIDIncorrect:        []string{},
		PreferredChannelTypeDeleted: []string{},
	}
	var (
		wg sync.WaitGroup
		mu sync.Mutex = sync.Mutex{}
	)

	errorChan := make(chan error)
	stopChan := make(chan bool)
	go func() {
		for _, recipient := range info.Notifications.Recipients {
			recipientModel, err := recipients.GetRecipientWithRecipientID(recipient)
			if err == gorm.ErrRecordNotFound {
				openAPI.RecipientIDIncorrect = append(openAPI.RecipientIDIncorrect, recipient)
				continue
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, openAPI)
				return
			}
			channelList, err := channels.GetChannelsWithPriorityLessThan(uint(notification.Priority))
			if err != nil {
				c.JSON(http.StatusInternalServerError, openAPI)
				return
			}
			wg.Add(1)
			go notifications.SendAllNotifications(errorChan, stopChan, notification, *recipientModel, *channelList, &openAPI, &wg, &mu)
		}
		wg.Wait()
		close(errorChan)
		close(stopChan)
	}()
	err = <-errorChan
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, openAPI)
		return
	}
	c.JSON(http.StatusOK, openAPI)
}
