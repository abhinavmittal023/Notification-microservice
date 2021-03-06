package notifications

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipientnotifications"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// Notifications interface is used to send different types of notifications
type Notifications interface {
	SendNotification() error
	New(to string, title string, body string)
}

// NewNotification interface is used to send notifications directly
type NewNotification interface {
	New(recipientNotification *models.RecipientNotifications, to string, title string, body string, notificationType Notifications) (int, error)
}

// CreateNotification is a struct that implements NewNotification Interface
type CreateNotification struct {
	Retry int
}

// New is the function to send real notifications
func (notification CreateNotification) New(recipientNotification *models.RecipientNotifications, to string, title string, body string, notificationType Notifications) (int, error) {
	notificationType.New(to, title, body)
	indx := 0
	var err error = nil
	for ok := true; ok; ok = indx <= notification.Retry {
		err = notificationType.SendNotification()
		if err != nil {
			recipientNotification.Status = constants.Failure
			err2 := recipientnotifications.PatchRecipientNotification(recipientNotification)
			if err2 != nil {
				indx = notification.Retry
				return http.StatusInternalServerError, err2
			}
			indx++
		} else {
			recipientNotification.Status = constants.Success
			err2 := recipientnotifications.PatchRecipientNotification(recipientNotification)
			if err2 != nil {
				indx = notification.Retry
				return http.StatusInternalServerError, err2
			}
			indx = notification.Retry
			return http.StatusOK, nil
		}
	}
	return http.StatusOK, err
}
