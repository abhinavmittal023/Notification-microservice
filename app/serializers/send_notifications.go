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
	Title      string   `json:"title" binding:"required,max=300"`
	Body       string   `json:"body" binding:"required"`
	Retry      int      `json:"retry"`
}

// NotificationsInfoToNotificationModel converts the serializer to model
func NotificationsInfoToNotificationModel(info *SendNotifications) (*models.Notification, error) {
	var notification models.Notification
	notification.Priority = constants.PriorityTypeToInt(info.Notifications.Priority)
	if notification.Priority == 0 {
		return nil, errors.New("Invalid Priority")
	}
	notification.Title = info.Notifications.Title
	notification.Body = info.Notifications.Body
	return &notification, nil
}
