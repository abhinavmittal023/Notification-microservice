package notifications

import (
	"time"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// GetAllNotifications gets all the notifications from the database and returns []models.RecipientNotifications, err
func GetAllNotifications(pagination *serializers.Pagination, notificationFilter *filter.Notification) ([]models.RecipientNotifications, error) {
	return nil, nil
	// dbG := db.Get()
	// var recipientNotifications []models.RecipientNotifications

	// var recipientNotification models.RecipientNotifications

	// var notificationsInfo []serializers.NotificationsInfo
	// data := make(map[uint64]serializers.NotificationsInfo)

	// if notificationFilter.RecipientID != "" {
	// 	err := db.Get().Table("recipients").Select("id").Where("recipient_id = ?", notificationFilter.RecipientID).Find(&recipientNotification).Error
	// 	if err != nil{
	// 		return nil,err
	// 	}
	// }
	// tx := dbG.Model(&models.RecipientNotifications{})
	// tx = tx.Select("distinct(notification_id)")
	// if notificationFilter.RecipientID != "" {
	// 	tx = tx.Where("recipient_id = ?", recipientNotification.ID)
	// }
	// if notificationFilter.ChannelName != "" {
	// 	tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
	// }
	// if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
	// 	tx = tx.Where("updated_at BETWEEN ? AND ?",notificationFilter.From, notificationFilter.To )
	// }

	// res := tx.Offset(pagination.Offset).Limit(pagination.Limit).Find(&recipientNotifications)

	// tx = dbG.Model(&models.Notification{})

	// for _,result := range(recipientNotifications){
	// 	var notification models.Notification
	// 	res = tx.Where("id = ?",result.NotificationID).Find(&notification)
	// 	if res!=nil{
	// 		return nil,res.Error
	// 	}
	// 	data[result.NotificationID] = serializers.NotificationsInfo{
	// 		Title: notification.Title,
	// 		Body: notification.Body,
	// 		Priority: notification.Priority,
	// 		NotificationChannels: []serializers.NotificationChannels{},
	// 	}
	// }

	// return recipientNotifications, res.Error
}

// GetAllNotificationsCount gets all the notifications count from the database and returns records count,err
func GetAllNotificationsCount(notificationFilter *filter.Notification) (int64, error) {

	dbG := db.Get()
	tx := dbG.Model(&models.RecipientNotifications{})
	if notificationFilter.RecipientID != "" {
		tx = tx.Where("recipient_id = ?", db.Get().Table("recipients").Select("id").Where("recipient_id = ?", notificationFilter.RecipientID).SubQuery())
	}
	if notificationFilter.ChannelName != "" {
		tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
	}
	if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
		tx = tx.Where("updated_at BETWEEN ? AND ?", notificationFilter.From, notificationFilter.To)
	}

	var count int64
	res := tx.Count(&count)
	return count, res.Error
}

// GetGraphData is used to get the required graph data
func GetGraphData(notificationFilter *filter.Notification) (*serializers.GraphData, error) {
	dbG := db.Get()
	var graphData serializers.GraphData
	graphKey := []time.Time{}
	type successFailedData struct {
		Successful uint64 `json:"successful"`
		Failed     uint64 `json:"failed"`
		Total      uint64 `json:"total"`
	}
	data := make(map[time.Time]successFailedData)
	var results []serializers.SuccessFailedData
	tx := dbG.Model(&models.RecipientNotifications{})
	if notificationFilter.RecipientID != "" {
		tx = tx.Where("recipient_id = ?", db.Get().Table("recipients").Select("id").Where("recipient_id = ?", notificationFilter.RecipientID).SubQuery())
	}
	if notificationFilter.ChannelName != "" {
		tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
	}
	if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
		tx = tx.Where("updated_at BETWEEN ? AND ?", notificationFilter.From, notificationFilter.To)
	}
	res := tx.Select("updated_at, count(*) as successful").Where("status = ?", constants.Success).Group("updated_at").Order("updated_at").Scan(&results)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, result := range results {
		data[result.UpdatedAt.Truncate(time.Minute)] = successFailedData{
			Successful: data[result.UpdatedAt.Truncate(time.Minute)].Successful + result.Successful,
			Total:      data[result.UpdatedAt.Truncate(time.Minute)].Successful + result.Successful,
		}
	}
	results = []serializers.SuccessFailedData{}
	tx = dbG.Model(&models.RecipientNotifications{})
	if notificationFilter.RecipientID != "" {
		tx = tx.Where("recipient_id = ?", db.Get().Table("recipients").Select("id").Where("recipient_id = ?", notificationFilter.RecipientID).SubQuery())
	}
	if notificationFilter.ChannelName != "" {
		tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
	}
	if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
		tx = tx.Where("updated_at BETWEEN ? AND ?", notificationFilter.From, notificationFilter.To)
	}
	res = tx.Select("updated_at, count(*) as failed").Where("status = ?", constants.Failure).Group("updated_at").Order("updated_at asc").Scan(&results)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, result := range results {
		graphKey = append(graphKey, result.UpdatedAt.Truncate(time.Minute))
		data[result.UpdatedAt.Truncate(time.Minute)] = successFailedData{
			Failed:     data[result.UpdatedAt.Truncate(time.Minute)].Failed + result.Failed,
			Successful: data[result.UpdatedAt.Truncate(time.Minute)].Successful,
			Total:      data[result.UpdatedAt.Truncate(time.Minute)].Failed + result.Failed + data[result.UpdatedAt.Truncate(time.Minute)].Successful,
		}
	}
	var prevVal time.Time = time.Time{}
	for _, val := range graphKey {
		if prevVal == val {
			continue
		}
		prevVal = val
		graphData.UpdatedAt = append(graphData.UpdatedAt, val)
		graphData.Successful = append(graphData.Successful, data[val].Successful)
		graphData.Failed = append(graphData.Failed, data[val].Failed)
		graphData.Total = append(graphData.Total, data[val].Total)
	}
	return &graphData, nil
}
