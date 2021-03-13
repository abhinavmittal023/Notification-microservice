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
func SendAllNotifications(errChan chan error, stopChan chan bool, notification models.Notification, recipientModel models.Recipient, channelList []models.Channel, openAPI *apimessage.OpenAPI, notificationInterface sendNotification.NewNotification, wg *sync.WaitGroup, mu *sync.Mutex) {

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
			email := sendNotification.Email{}
			status, err := notificationInterface.New(&recipientNotification, recipientModel.Email, notification.Title, notification.Body, &email)
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
			push := sendNotification.Push{}
			status, err := notificationInterface.New(&recipientNotification, recipientModel.PushToken, notification.Title, notification.Body, &push)
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
			web := sendNotification.Web{}
			status, err := notificationInterface.New(&recipientNotification, recipientModel.WebToken, notification.Title, notification.Body, &web)
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
			email := sendNotification.Email{}
			status, err := notificationInterface.New(&recipientNotification, recipientModel.Email, notification.Title, notification.Body, &email)
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
			push := sendNotification.Push{}
			status, err := notificationInterface.New(&recipientNotification, recipientModel.PushToken, notification.Title, notification.Body, &push)
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
			web := sendNotification.Web{}
			status, err := notificationInterface.New(&recipientNotification, recipientModel.WebToken, notification.Title, notification.Body, &web)
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

func SendToRecipients(channelList []models.Channel, recipientList []models.Recipient, openAPI *apimessage.OpenAPI, errorChan chan error, stopChan chan bool, notification models.Notification, notificationInterface sendNotification.NewNotification, mainWaitGroup *sync.WaitGroup) {

	defer mainWaitGroup.Done()
	var recipientWaitGroup sync.WaitGroup
	var mu sync.Mutex = sync.Mutex{}
	for _, recipient := range recipientList {
		recipientWaitGroup.Add(1)
		go SendAllNotifications(errorChan, stopChan, notification, recipient, channelList, openAPI, notificationInterface, &recipientWaitGroup, &mu)
	}
	recipientWaitGroup.Wait()
	errorChan <- nil
}
