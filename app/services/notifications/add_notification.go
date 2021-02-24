package notifications

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// AddNotification func creates a new channel in the database and returns nil/error
func AddNotification(notification *models.Notification) error {
	return db.Get().Create(notification).Error
}
