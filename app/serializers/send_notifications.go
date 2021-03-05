package serializers

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/pkg/errors"
)

// SendNotifications struct holds information about APIKey and Notifications
type SendNotifications struct {
	Notifications Notifications `json:"notifications" binding:"required"`
}

// Notifications serializer holds the information about notifications
type Notifications struct {
	Recipients []string `json:"recipients" binding:"required"`
	Priority   string   `json:"priority" binding:"required"`
	Title      string   `json:"title" binding:"required"`
	Body       string   `json:"body" binding:"required"`
}

// NotificationsInfoToNotificationModel converts the serializer to model
func NotificationsInfoToNotificationModel(info *SendNotifications, notification *models.Notification) error {
	notification.Priority = constants.PriorityTypeToInt(info.Notifications.Priority)
	if notification.Priority == 0 {
		return errors.New("Invalid Priority")
	}
	notification.Title = info.Notifications.Title
	notification.Body = info.Notifications.Body
	return nil
}
