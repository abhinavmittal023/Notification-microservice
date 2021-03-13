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
type CreateNotification struct{}

// New is the function to send real notifications
func (notification CreateNotification) New(recipientNotification *models.RecipientNotifications, to string, title string, body string, notificationType Notifications) (int, error) {
	notificationType.New(to, title, body)
	err := notificationType.SendNotification()
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

// MockNotification is the struct that implements the NewNotification interface
type MockNotification struct{}

// New is the function that mocks notification send service
func (notification MockNotification) New(recipientNotification *models.RecipientNotifications, to string, title string, body string, notificationType Notifications) (int, error) {
	return http.StatusOK, nil
}
