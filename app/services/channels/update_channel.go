package channels

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// PatchChannel function saves the updated parameters to the database records
func PatchChannel(channel *models.Channel) error {
	return db.Get().Save(channel).Error
}
