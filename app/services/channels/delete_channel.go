package channels

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// DeleteChannel function soft deletes the channel and returns nil/error
func DeleteChannel(channel *models.Channel) error {
	return db.Get().Delete(channel).Error
}
