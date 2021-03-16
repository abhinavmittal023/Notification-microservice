package logs

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// GetAllLogs gets all the logs from the database and returns []models.Logs,err
func GetAllLogs(pagination *serializers.Pagination, logFilter *filter.Logs) ([]models.Logs, error) {

	var logs []models.Logs
	dbG := db.Get()
	tx := dbG.Model(&models.Logs{})

	if logFilter.Level != 0 {
		tx = tx.Where("level = ?", logFilter.Level)
	}
	if pagination.Limit != 0 {
		tx = tx.Offset(pagination.Offset).Limit(pagination.Limit)
	}
	res := tx.Order("created_at").Find(&logs)
	return logs, res.Error
}

// GetAllLogsCount gets all the logs count from the database and returns records count,err
func GetAllLogsCount(logFilter *filter.Logs) (int64, error) {

	dbG := db.Get()
	tx := dbG.Model(&models.Logs{})

	if logFilter.Level != 0 {
		tx = tx.Where("level = ?", logFilter.Level)
	}

	var count int64
	res := tx.Count(&count)
	return count, res.Error
}