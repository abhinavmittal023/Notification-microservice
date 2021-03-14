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
func SendAllNotifications(errChan chan error, stopChan chan bool, notification models.Notification, recipientModel models.Recipient, channelList []models.Channel, messageChan chan apimessage.Message, notificationInterface sendNotification.NewNotification, wg *sync.WaitGroup) {

	defer wg.Done()

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
			select {
			case <-stopChan:
				return
			default:
				err := recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
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
						select {
						case <-stopChan:
							return
						case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
						}
						select {
						case <-stopChan:
							return
						case errChan <- err:
							return
						}
					}
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
					continue
				}
				select {
				case <-stopChan:
					return
				case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Success, ChannelName: channel.Name}:
				}
			}
		} else if constants.ChannelType(uint(channel.Type)) == "Push" && recipientModel.PushToken != "" {
			channelSent["Push"] = true
			select {
			case <-stopChan:
				return
			default:
				err := recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
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
						select {
						case <-stopChan:
							return
						case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
						}
						select {
						case <-stopChan:
							return
						case errChan <- err:
							return
						}
					}
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
					continue
				}
				select {
				case <-stopChan:
					return
				case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Success, ChannelName: channel.Name}:
				}
			}
		} else if constants.ChannelType(uint(channel.Type)) == "Web" && recipientModel.WebToken != "" {
			channelSent["Web"] = true
			select {
			case <-stopChan:
				return
			default:
				err := recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
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
						select {
						case <-stopChan:
							return
						case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
						}
						select {
						case <-stopChan:
							return
						case errChan <- err:
							return
						}
					}
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
					continue
				}
				select {
				case <-stopChan:
					return
				case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Success, ChannelName: channel.Name}:
				}
			}
		}
	}

	if recipientModel.PreferredChannelType > 0 && !channelSent[constants.ChannelType(recipientModel.PreferredChannelType)] {
		channel, err := channels.GetChannelWithType(recipientModel.PreferredChannelType)
		if err == gorm.ErrRecordNotFound {
			select {
			case <-stopChan:
				return
			case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.PreferredChannelTypeDeleted, ChannelName: channel.Name}:
			}
			return
		}
		if err != nil {
			select {
			case <-stopChan:
				return
			case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
			}
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
			select {
			case <-stopChan:
				return
			default:
				err := recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
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
						select {
						case <-stopChan:
							return
						case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
						}
						select {
						case <-stopChan:
							return
						case errChan <- err:
							return
						}
					}
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
					return
				}
				select {
				case <-stopChan:
					return
				case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Success, ChannelName: channel.Name}:
				}
			}
		} else if constants.ChannelType(uint(channel.Type)) == "Push" && recipientModel.PushToken != "" {
			channelSent["Push"] = true
			select {
			case <-stopChan:
				return
			default:
				err := recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
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
						select {
						case <-stopChan:
							return
						case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
						}
						select {
						case <-stopChan:
							return
						case errChan <- err:
							return
						}
					}
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
					return
				}
				select {
				case <-stopChan:
					return
				case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Success, ChannelName: channel.Name}:
				}
			}
		} else if constants.ChannelType(uint(channel.Type)) == "Web" && recipientModel.WebToken != "" {
			channelSent["Web"] = true
			select {
			case <-stopChan:
				return
			default:
				err := recipientnotifications.AddRecipientNotification(&recipientNotification)
				if err != nil {
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
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
						select {
						case <-stopChan:
							return
						case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
						}
						select {
						case <-stopChan:
							return
						case errChan <- err:
							return
						}
					}
					select {
					case <-stopChan:
						return
					case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Failure, ChannelName: channel.Name}:
					}
					return
				}
				select {
				case <-stopChan:
					return
				case messageChan <- apimessage.Message{ID: recipientModel.RecipientID, Option: apimessage.Success, ChannelName: channel.Name}:
				}
			}
		}
	}
}

// SendToRecipients function sends the notification to all recipients
func SendToRecipients(channelList []models.Channel, recipientList []models.Recipient, openAPI *apimessage.OpenAPI, errorChan chan error, stopChan chan bool, notification models.Notification, notificationInterface sendNotification.NewNotification, mainWaitGroup *sync.WaitGroup) {

	defer mainWaitGroup.Done()
	var recipientWaitGroup sync.WaitGroup

	messageChan := make(chan apimessage.Message)
	mainWaitGroup.Add(1)
	go openAPI.AddStatus(stopChan, messageChan, mainWaitGroup)

	for _, recipient := range recipientList {
		recipientWaitGroup.Add(1)
		go SendAllNotifications(errorChan, stopChan, notification, recipient, channelList, messageChan, notificationInterface, &recipientWaitGroup)
	}
	recipientWaitGroup.Wait()
	close(errorChan)
	close(messageChan)
}
