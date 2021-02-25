package notifications

import (
	"fmt"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/notifications"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipientnotifications"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipients"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	sendNotification "code.jtg.tools/ayush.singhal/notifications-microservice/shared/notifications"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notifications info is required"})
		return
	}
	var notification models.Notification
	serializers.NotificationsInfoToNotificationModel(&info, &notification)
	err := notifications.AddNotification(&notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	}
	var errors serializers.ErrorInfo
	errors.Error = make(map[int][]string)
	var errorFound = false
	for index, recipient := range info.Notifications.Recipients {
		var errorMap []string
		recipientModel, err := recipients.GetRecipientWithRecipientID(recipient)
		if err == gorm.ErrRecordNotFound {
			errorMap = append(errorMap, "Recipient ID incorrect")
			errors.Error[index] = errorMap
			errorFound = true
			continue
		} else if err != nil {
			errorMap = append(errorMap, "Internal Server Error")
			errors.Error[index] = errorMap
			errorFound = true
			c.JSON(http.StatusInternalServerError, errors)
			return
		}
		channelType := 1 // Email
		channel, err := channels.GetChannelWithType(uint(channelType))
		if err == gorm.ErrRecordNotFound {
			errorMap = append(errorMap, fmt.Sprintf("Channel Type %s was deleted", constants.ChannelType(uint(channelType))))
			errors.Error[index] = errorMap
			errorFound = true
			continue
		} else if err != nil {
			errorMap = append(errorMap, "Internal Server Error")
			errors.Error[index] = errorMap
			errorFound = true
			c.JSON(http.StatusInternalServerError, errors)
			return
		}
		recipientNotification := models.RecipientNotifications{
			NotificationID: uint64(notification.ID),
			RecipientID:    uint64(recipientModel.ID),
			ChannelID:      uint64(channel.ID),
			Status:         constants.Pending,
		}
		err = recipientnotifications.AddRecipientNotification(&recipientNotification)
		if err != nil {
			errorMap = append(errorMap, "Internal Server Error")
			errors.Error[index] = errorMap
			errorFound = true
			c.JSON(http.StatusInternalServerError, errors)
			return
		}
		if recipientModel.Email != "" {
			email := sendNotification.Email{
				To:      recipientModel.Email,
				Subject: info.Notifications.Title,
				Message: info.Notifications.Body,
			}
			err = email.SendNotification()
			if err != nil {
				errorMap = append(errorMap, "Failure")
				errors.Error[index] = errorMap
				errorFound = true
				recipientNotification.Status = constants.Failure
				recipientnotifications.PatchRecipientNotification(&recipientNotification)
				continue
			}
			recipientNotification.Status = constants.Success
			recipientnotifications.PatchRecipientNotification(&recipientNotification)
		}
	}
	if errorFound {
		c.JSON(http.StatusBadRequest, errors)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
