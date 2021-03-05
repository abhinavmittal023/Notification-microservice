package notifications

import (
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
	var errors serializers.ErrorInfo
	errors.Error = make(map[int][]string)
	var errorFound = false

	for index, recipient := range info.Notifications.Recipients {
		var errorMap []string
		channelSent := map[string]bool{}
		recipientModel, err := recipients.GetRecipientWithRecipientID(recipient)
		if err == gorm.ErrRecordNotFound {
			errorMap = append(errorMap, constants.Errors().RecipientIDIncorrect)
			errors.Error[index] = errorMap
			errorFound = true
			continue
		} else if err != nil {
			errorMap = append(errorMap, constants.Errors().InternalError)
			errors.Error[index] = errorMap
			errorFound = true
			c.JSON(http.StatusInternalServerError, errors)
			return
		}
		channelList, err := channels.GetChannelsWithPriorityLessThan(uint(notification.Priority))
		if err != nil {
			errorMap = append(errorMap, constants.Errors().InternalError)
			errors.Error[index] = errorMap
			errorFound = true
			c.JSON(http.StatusInternalServerError, errors)
			return
		}
		for _, channel := range *channelList {

			recipientNotification := models.RecipientNotifications{
				NotificationID: uint64(notification.ID),
				RecipientID:    uint64(recipientModel.ID),
				ChannelName:    channel.Name,
				Status:         constants.Pending,
			}

			if constants.ChannelType(uint(channel.Type)) == "Email" && recipientModel.Email != "" {
				channelSent["Email"] = true
				err = recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					errorMap = append(errorMap, constants.Errors().InternalError)
					errors.Error[index] = errorMap
					errorFound = true
					c.JSON(http.StatusInternalServerError, errors)
					return
				}
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

			} else if constants.ChannelType(uint(channel.Type)) == "Push" && recipientModel.PushToken != "" {
				channelSent["Push"] = true
				err = recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					errorMap = append(errorMap, constants.Errors().InternalError)
					errors.Error[index] = errorMap
					errorFound = true
					c.JSON(http.StatusInternalServerError, errors)
					return
				}
				push := sendNotification.Push{
					To:    recipientModel.PushToken,
					Title: info.Notifications.Title,
					Body:  info.Notifications.Body,
				}
				err = push.SendNotification()
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

			} else if constants.ChannelType(uint(channel.Type)) == "Web" && recipientModel.WebToken != "" {
				channelSent["Web"] = true
				err = recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					errorMap = append(errorMap, constants.Errors().InternalError)
					errors.Error[index] = errorMap
					errorFound = true
					c.JSON(http.StatusInternalServerError, errors)
					return
				}
				web := sendNotification.Web{
					To:    recipientModel.WebToken,
					Title: info.Notifications.Title,
					Body:  info.Notifications.Body,
				}
				err = web.SendNotification()
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
		if recipientModel.PreferredChannelType > 0 && !channelSent[constants.ChannelType(recipientModel.PreferredChannelType)] {
			channel, err := channels.GetChannelWithType(recipientModel.PreferredChannelType)
			if err == gorm.ErrRecordNotFound {
				errorMap = append(errorMap, "Preferred Channel was Deleted")
				errors.Error[index] = errorMap
				errorFound = true
				continue
			}
			if err != nil {
				errorMap = append(errorMap, constants.Errors().InternalError)
				errors.Error[index] = errorMap
				errorFound = true
				c.JSON(http.StatusInternalServerError, errors)
				return
			}

			recipientNotification := models.RecipientNotifications{
				NotificationID: uint64(notification.ID),
				RecipientID:    uint64(recipientModel.ID),
				ChannelName:    channel.Name,
				Status:         constants.Pending,
			}

			if constants.ChannelType(uint(channel.Type)) == "Email" && recipientModel.Email != "" {
				channelSent["Email"] = true
				err = recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					errorMap = append(errorMap, constants.Errors().InternalError)
					errors.Error[index] = errorMap
					errorFound = true
					c.JSON(http.StatusInternalServerError, errors)
					return
				}
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

			} else if constants.ChannelType(uint(channel.Type)) == "Push" && recipientModel.PushToken != "" {
				channelSent["Push"] = true
				err = recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					errorMap = append(errorMap, constants.Errors().InternalError)
					errors.Error[index] = errorMap
					errorFound = true
					c.JSON(http.StatusInternalServerError, errors)
					return
				}
				push := sendNotification.Push{
					To:    recipientModel.PushToken,
					Title: info.Notifications.Title,
					Body:  info.Notifications.Body,
				}
				err = push.SendNotification()
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

			} else if constants.ChannelType(uint(channel.Type)) == "Web" && recipientModel.WebToken != "" {
				channelSent["Web"] = true
				err = recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					errorMap = append(errorMap, constants.Errors().InternalError)
					errors.Error[index] = errorMap
					errorFound = true
					c.JSON(http.StatusInternalServerError, errors)
					return
				}
				web := sendNotification.Web{
					To:    recipientModel.WebToken,
					Title: info.Notifications.Title,
					Body:  info.Notifications.Body,
				}
				err = web.SendNotification()
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
	}
	if errorFound {
		c.JSON(http.StatusBadRequest, errors)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
