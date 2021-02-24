package recipientnotifications

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// PatchRecipientNotification function saves the updated parameters to the database records
func PatchRecipientNotification(recipientNotification *models.RecipientNotifications) error {
	return db.Get().Save(recipientNotification).Error
}
