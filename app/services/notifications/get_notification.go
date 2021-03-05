package notifications

import (
	"time"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// GetAllNotifications gets all the notifications from the database and returns []serializers.NotificationsInfo, err
func GetAllNotifications(pagination *serializers.Pagination, notificationFilter *filter.Notification) ([]serializers.NotificationsInfo, error) {
	dbG := db.Get()
    var recipientNotifications []serializers.NotificationsInfo

    var recipientNotification models.RecipientNotifications

	type successFailedData struct{
		Successful uint64 
        Failed     uint64 
        Total      uint64 
	}

	type notificationData struct {
		NotificationID uint64
	}

    type structResponse struct{
        Priority             int                    
        Title                string                 
        Body                 string                 
        NotificationChannels map[string]successFailedData 
	}

    data := make(map[uint64]structResponse)

    if notificationFilter.RecipientID != "" {
        err := db.Get().Table("recipients").Select("id").Where("recipient_id = ?", notificationFilter.RecipientID).Find(&recipientNotification).Error
        if err != nil{
            return nil,err
        }
	}
	var notificationID []notificationData
	tx := dbG.Model(&models.RecipientNotifications{})
	if notificationFilter.RecipientID != "" {
		tx = tx.Where("recipient_id = ?", recipientNotification.ID)
	}
	if notificationFilter.ChannelName != "" {
		tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
	}
	if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
		tx = tx.Where("updated_at BETWEEN ? AND ?",notificationFilter.From, notificationFilter.To )
	}
	res := tx.Offset(pagination.Offset).Limit(pagination.Limit).Select("distinct(notification_id)").Scan(&notificationID)
	if res.Error != nil {
		return nil, res.Error
	}

	for _,value := range notificationID{

		if _,found := data[uint64(value.NotificationID)]; !found{
			tx = dbG.Model(&models.Notification{})
			var notification models.Notification
			res = tx.Where("id = ?",value.NotificationID).Find(&notification)
			data[uint64(value.NotificationID)] = structResponse{
				Body: notification.Body,
				Title: notification.Title,
				Priority: notification.Priority,
				NotificationChannels: make(map[string]successFailedData),
			}
		}
		
		tx = dbG.Model(&models.RecipientNotifications{})
		if notificationFilter.RecipientID != "" {
			tx = tx.Where("recipient_id = ?", recipientNotification.ID)
		}
		if notificationFilter.ChannelName != "" {
			tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
		}
		if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
			tx = tx.Where("updated_at BETWEEN ? AND ?",notificationFilter.From, notificationFilter.To )
		}
		results := []serializers.NotificationChannels{}

		res = tx.Select("channel_name, count(*) as successful").Where("notification_id = ? and status = ?", value.NotificationID , constants.Success).Group("channel_name").Scan(&results)
		if res.Error != nil {
			return nil, res.Error
		}
		for _,val := range results {
			data[uint64(value.NotificationID)].NotificationChannels[val.ChannelName] = successFailedData{
				Successful: val.Successful,
				Total: val.Successful,
			}
		}
	}

	for _,value := range notificationID{

		if _,found := data[uint64(value.NotificationID)]; !found{
			tx = dbG.Model(&models.Notification{})
			var notification models.Notification
			res = tx.Where("id = ?",value.NotificationID).Find(&notification)
			data[uint64(value.NotificationID)] = structResponse{
				Body: notification.Body,
				Title: notification.Title,
				Priority: notification.Priority,
				NotificationChannels: make(map[string]successFailedData),
			}
		}
		tx = dbG.Model(&models.RecipientNotifications{})
		if notificationFilter.RecipientID != "" {
			tx = tx.Where("recipient_id = ?", recipientNotification.ID)
		}
		if notificationFilter.ChannelName != "" {
			tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
		}
		if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
			tx = tx.Where("updated_at BETWEEN ? AND ?",notificationFilter.From, notificationFilter.To )
		}
		results := []serializers.NotificationChannels{}

		res = tx.Select("channel_name, count(*) as failure").Where("notification_id = ? and status = ?", value.NotificationID , constants.Failure).Group("channel_name").Scan(&results)
		if res.Error != nil {
			return nil, res.Error
		}
		for _,val := range results {
			data[uint64(value.NotificationID)].NotificationChannels[val.ChannelName] = successFailedData{
				Successful: data[uint64(value.NotificationID)].NotificationChannels[val.ChannelName].Successful,
            	Failed: val.Failure,
            	Total: data[uint64(value.NotificationID)].NotificationChannels[val.ChannelName].Successful + val.Failure,
			}
		}
	}

    var channelSlice []serializers.NotificationChannels
    for _,val := range data{
        channelSlice = []serializers.NotificationChannels{}
        for name,channel := range val.NotificationChannels{
            channelSlice = append(channelSlice, serializers.NotificationChannels{
                ChannelName: name,
                Successful: channel.Successful,
                Failure: channel.Failed,
                Total: channel.Total,
            })
        }
        recipientNotifications = append(recipientNotifications, serializers.NotificationsInfo{
            Priority: val.Priority,
            Title: val.Title,
            Body: val.Body,
            NotificationChannels: channelSlice,
        })
    }

    return recipientNotifications, res.Error
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

	var count struct{
		Count int64
	}
	res := tx.Select("count(distinct(notification_id)) as count").Scan(&count)
	return count.Count, res.Error
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
