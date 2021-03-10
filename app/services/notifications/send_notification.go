package notifications

import (
	"net/http"
	"sync"

	apimessage "code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/api_message"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipientnotifications"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	sendNotification "code.jtg.tools/ayush.singhal/notifications-microservice/shared/notifications"
	"github.com/jinzhu/gorm"
)

// SendAllNotifications functon sends the notification to the specific recipient
func SendAllNotifications(errChan chan error, stopChan chan bool, notification models.Notification, recipientModel models.Recipient, channelList []models.Channel, openAPI *apimessage.OpenAPI, wg *sync.WaitGroup, mu *sync.Mutex) {

	defer wg.Done()

	select {
	case <-stopChan:
		return
	default:
	}

	channelSent := map[string]bool{}

	for _, channel := range channelList {

		recipientNotification := models.RecipientNotifications{
			NotificationID: uint64(notification.ID),
			RecipientID:    uint64(recipientModel.ID),
			ChannelName:    channel.Name,
			Status:         constants.Pending,
		}

		if constants.ChannelType(uint(channel.Type)) == "Email" && recipientModel.Email != "" {
			channelSent["Email"] = true
			err := recipientnotifications.AddRecipientNotification(&recipientNotification)
			if err != nil {
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				select {
				case <-stopChan:
					return
				case errChan <- err:
					return
				}
			}
			email := sendNotification.Email{To: recipientModel.Email, Subject: notification.Title, Message: notification.Body}
			status, err := send(&recipientNotification, &email)
			if err != nil {
				if status == http.StatusInternalServerError {
					openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
					select {
					case <-stopChan:
						return
					case errChan <- err:
						return
					}
				}
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				continue
			}
			openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, true, mu)
		} else if constants.ChannelType(uint(channel.Type)) == "Push" && recipientModel.PushToken != "" {
			channelSent["Push"] = true
			err := recipientnotifications.AddRecipientNotification(&recipientNotification)
			if err != nil {
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				select {
				case <-stopChan:
					return
				case errChan <- err:
					return
				}
			}
			push := sendNotification.Push{To: recipientModel.PushToken, Title: notification.Title, Body: notification.Body}
			status, err := send(&recipientNotification, &push)
			if err != nil {
				if status == http.StatusInternalServerError {
					openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
					select {
					case <-stopChan:
						return
					case errChan <- err:
						return
					}
				}
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				continue
			}
			openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, true, mu)
		} else if constants.ChannelType(uint(channel.Type)) == "Web" && recipientModel.WebToken != "" {
			channelSent["Web"] = true
			err := recipientnotifications.AddRecipientNotification(&recipientNotification)
			if err != nil {
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				select {
				case <-stopChan:
					return
				case errChan <- err:
					return
				}
			}
			web := sendNotification.Web{To: recipientModel.WebToken, Title: notification.Title, Body: notification.Body}
			status, err := send(&recipientNotification, &web)
			if err != nil {
				if status == http.StatusInternalServerError {
					openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
					select {
					case <-stopChan:
						return
					case errChan <- err:
						return
					}
				}
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				continue
			}
			openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, true, mu)
		}
	}

	if recipientModel.PreferredChannelType > 0 && !channelSent[constants.ChannelType(recipientModel.PreferredChannelType)] {
		channel, err := channels.GetChannelWithType(recipientModel.PreferredChannelType)
		if err == gorm.ErrRecordNotFound {
			openAPI.PreferredChannelTypeDeleted = append(openAPI.PreferredChannelTypeDeleted, recipientModel.RecipientID)
			return
		}
		if err != nil {
			openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
			select {
			case <-stopChan:
				return
			case errChan <- err:
				return
			}
		}

		recipientNotification := models.RecipientNotifications{
			NotificationID: uint64(notification.ID),
			RecipientID:    uint64(recipientModel.ID),
			ChannelName:    channel.Name,
			Status:         constants.Pending,
		}
		if constants.ChannelType(uint(channel.Type)) == "Email" && recipientModel.Email != "" {
			channelSent["Email"] = true
			err := recipientnotifications.AddRecipientNotification(&recipientNotification)
			if err != nil {
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				select {
				case <-stopChan:
					return
				case errChan <- err:
					return
				}
			}
			email := sendNotification.Email{To: recipientModel.Email, Subject: notification.Title, Message: notification.Body}
			status, err := send(&recipientNotification, &email)
			if err != nil {
				if status == http.StatusInternalServerError {
					openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
					select {
					case <-stopChan:
						return
					case errChan <- err:
						return
					}
				}
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				return
			}
			openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, true, mu)
		} else if constants.ChannelType(uint(channel.Type)) == "Push" && recipientModel.PushToken != "" {
			channelSent["Push"] = true
			err := recipientnotifications.AddRecipientNotification(&recipientNotification)
			if err != nil {
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				select {
				case <-stopChan:
					return
				case errChan <- err:
					return
				}
			}
			push := sendNotification.Push{To: recipientModel.PushToken, Title: notification.Title, Body: notification.Body}
			status, err := send(&recipientNotification, &push)
			if err != nil {
				if status == http.StatusInternalServerError {
					openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
					select {
					case <-stopChan:
						return
					case errChan <- err:
						return
					}
				}
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				return
			}
			openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, true, mu)
		} else if constants.ChannelType(uint(channel.Type)) == "Web" && recipientModel.WebToken != "" {
			channelSent["Web"] = true
			err := recipientnotifications.AddRecipientNotification(&recipientNotification)
			if err != nil {
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				select {
				case <-stopChan:
					return
				case errChan <- err:
					return
				}
			}
			web := sendNotification.Web{To: recipientModel.WebToken, Title: notification.Title, Body: notification.Body}
			status, err := send(&recipientNotification, &web)
			if err != nil {
				if status == http.StatusInternalServerError {
					openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
					select {
					case <-stopChan:
						return
					case errChan <- err:
						return
					}
				}
				openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, false, mu)
				return
			}
			openAPI.AddRecipientID(recipientModel.RecipientID, channel.Name, true, mu)
		}
	}
}

func send(recipientNotification *models.RecipientNotifications, notification sendNotification.Notifications) (int, error) {

	err := notification.SendNotification()
	if err != nil {
		recipientNotification.Status = constants.Failure
		err2 := recipientnotifications.PatchRecipientNotification(recipientNotification)
		if err2 != nil {
			return http.StatusInternalServerError, err2
		}
		return http.StatusOK, err
	}
	recipientNotification.Status = constants.Success
	err = recipientnotifications.PatchRecipientNotification(recipientNotification)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
