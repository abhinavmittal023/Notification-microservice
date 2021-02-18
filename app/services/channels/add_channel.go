package channels

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// AddChannel func creates a new channel in the database and returns nil/error
func AddChannel(channel *models.Channel) error {
	return db.Get().Create(channel).Error
}
