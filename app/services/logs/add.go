package logs

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// AddLogs func creates a new log entry in the database and returns nil/error
func AddLogs(level uint, msg string) error {
	var channel models.Logs
	channel.Level = level
	channel.Msg = msg
	return db.Get().Create(&channel).Error
}
