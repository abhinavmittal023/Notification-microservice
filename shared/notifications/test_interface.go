package notifications

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// MockNotification is the struct that implements the NewNotification interface
type MockNotification struct{}

// New is the function that mocks notification send service
func (notification MockNotification) New(recipientNotification *models.RecipientNotifications, to string, title string, body string, notificationType Notifications) (int, error) {
	return http.StatusOK, nil
}
