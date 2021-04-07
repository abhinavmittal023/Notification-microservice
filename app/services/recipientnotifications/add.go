package recipientnotifications

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// AddRecipientNotification func creates a new entry in the database and returns nil/error
func AddRecipientNotification(recipientNotification *models.RecipientNotifications) error {
	return db.Get().Create(recipientNotification).Error
}
