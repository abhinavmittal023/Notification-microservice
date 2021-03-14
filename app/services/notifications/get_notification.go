package notifications

import (
	"time"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/jinzhu/gorm"
)

type notificationData struct {
	NotificationID uint64
}

// GetAllNotifications gets all the notifications from the database and returns []serializers.NotificationsInfo, err
func GetAllNotifications(pagination *serializers.Pagination, notificationFilter *filter.Notification) ([]serializers.NotificationsInfo, error) {
	dbG := db.Get()
	var recipientNotifications []serializers.NotificationsInfo
	var recipient models.Recipient

	type successFailedData struct {
		ChannelName string
		Status      int
		Count       uint64
	}

	type recipientInfo struct{
		RecipientID	string
		ChannelName	string
		Status		uint64
	}

	if notificationFilter.RecipientID != "" {
		err := dbG.Table("recipients").Select("id").Where("recipient_id = ?", notificationFilter.RecipientID).Find(&recipient).Error
		if err != nil {
			return nil, err
		}
	}
	var notificationID []notificationData
	tx := dbG.Model(&models.RecipientNotifications{})
	if notificationFilter.RecipientID != "" {
		tx = tx.Where("recipient_id = ?", recipient.ID)
	}
	if notificationFilter.ChannelName != "" {
		tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
	}
	if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
		tx = tx.Where("updated_at BETWEEN ? AND ?", notificationFilter.From, notificationFilter.To)
	}
	res := tx.Offset(pagination.Offset).Limit(pagination.Limit).Select("distinct(notification_id)").Order("notification_id").Scan(&notificationID)
	if res.Error != nil {
		return nil, res.Error
	}

	for i, value := range notificationID {

		tx = dbG.Model(&models.RecipientNotifications{})
		if notificationFilter.RecipientID != "" {
			tx = tx.Where("recipient_id = ?", recipient.ID)
		}
		if notificationFilter.ChannelName != "" {
			tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
		}
		if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
			tx = tx.Where("updated_at BETWEEN ? AND ?", notificationFilter.From, notificationFilter.To)
		}
		results := []successFailedData{}

		res = tx.Select("channel_name, status, count(*) as count").Where("notification_id = ?", value.NotificationID).Group("channel_name,status").Order("channel_name").Scan(&results)
		if res.Error != nil {
			return nil, res.Error
		}
		tx = dbG.Model(&models.Notification{})
		var notification models.Notification
		res = tx.Where("id = ?", value.NotificationID).Find(&notification)
		if res.Error != nil {
			return nil, res.Error
		}
		var prevVal successFailedData

		recipientNotifications = append(recipientNotifications, serializers.NotificationsInfo{
			Priority:  notification.Priority,
			Title:     notification.Title,
			Body:      notification.Body,
			CreatedAt: notification.CreatedAt,
		})
		var notificationChannels []serializers.NotificationChannels
		idx := 0
		for _, val := range results {
			if prevVal.ChannelName == val.ChannelName {
				if val.Status == constants.Success {
					notificationChannels[idx-1] = serializers.NotificationChannels{
						ChannelName: val.ChannelName,
						Successful:  val.Count,
						Failure:     notificationChannels[idx-1].Failure,
						Total:       val.Count + notificationChannels[idx-1].Failure,
					}
				} else if val.Status == constants.Failure {
					notificationChannels[idx-1] = serializers.NotificationChannels{
						ChannelName: val.ChannelName,
						Failure:     val.Count,
						Successful:  notificationChannels[idx-1].Successful,
						Total:       val.Count + notificationChannels[idx-1].Successful,
					}
				}
			} else {
				if val.Status == constants.Success {
					notificationChannels = append(notificationChannels, serializers.NotificationChannels{
						ChannelName: val.ChannelName,
						Successful:  val.Count,
						Total:       val.Count,
					})
				} else if val.Status == constants.Failure {
					notificationChannels = append(notificationChannels, serializers.NotificationChannels{
						ChannelName: val.ChannelName,
						Failure:     val.Count,
						Total:       val.Count,
					})
				}
				idx++
			}
			prevVal = val
		}

		var recipientsInfo []recipientInfo

		recipientNotifications[i] = serializers.NotificationsInfo{
			Priority:             notification.Priority,
			Title:                notification.Title,
			Body:                 notification.Body,
			CreatedAt:            notification.CreatedAt,
			NotificationChannels: notificationChannels,
		}
		res = dbG.Table("recipients").Select("recipients.recipient_id,channel_name,status").Joins("join recipient_notifications on recipients.id = recipient_notifications.recipient_id").Where("recipient_notifications.notification_id = ?",value.NotificationID).Order("recipients.recipient_id").Scan(&recipientsInfo)

		if res.Error == gorm.ErrRecordNotFound {
			continue;
		}else if res.Error != nil{
			return nil, res.Error
		}
		var prevRecipient recipientInfo
		var recipients []serializers.Recipients
		idx = 0
		for _,val := range recipientsInfo{
			if prevRecipient.RecipientID == val.RecipientID{
				recipients[idx-1].Channels = append(recipients[idx-1].Channels,serializers.Channels{
					ChannelName: val.ChannelName,
					Status: val.Status,
				})
			}else{
				recipients = append(recipients, serializers.Recipients{
					RecipientID: val.RecipientID,
					Channels: []serializers.Channels{
						{
							ChannelName: val.ChannelName,
							Status: val.Status,
						},
					},
				})
				idx++
			}
			prevRecipient = val
		}
		recipientNotifications[i].Recipients = recipients
	}
	return recipientNotifications, nil
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

	var count struct {
		Count int64
	}
	res := tx.Select("count(distinct(notification_id)) as count").Scan(&count)
	return count.Count, res.Error
}

// GetGraphData is used to get the required graph data
func GetGraphData(notificationFilter *filter.Notification) (*serializers.GraphData, error) {
	dbG := db.Get()
	var recipient models.Recipient
	var graphData serializers.GraphData
	type successFailedData struct {
		Status int
		Count  uint64
	}

	if notificationFilter.RecipientID != "" {
		err := dbG.Table("recipients").Select("id").Where("recipient_id = ?", notificationFilter.RecipientID).Find(&recipient).Error
		if err != nil {
			return nil, err
		}
	}

	var notificationID []notificationData
	tx := dbG.Model(&models.RecipientNotifications{})
	if notificationFilter.RecipientID != "" {
		tx = tx.Where("recipient_id = ?", recipient.ID)
	}
	if notificationFilter.ChannelName != "" {
		tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
	}
	if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
		tx = tx.Where("updated_at BETWEEN ? AND ?", notificationFilter.From, notificationFilter.To)
	}
	res := tx.Select("distinct(notification_id)").Order("notification_id").Scan(&notificationID)
	if res.Error != nil {
		return nil, res.Error
	}

	for i, value := range notificationID {

		tx = dbG.Model(&models.Notification{})
		var notification models.Notification
		res = tx.Where("id = ?", value.NotificationID).Find(&notification)
		if res.Error != nil {
			return nil, res.Error
		}
		graphData.UpdatedAt = append(graphData.UpdatedAt, notification.UpdatedAt.Truncate(time.Minute))

		tx = dbG.Model(&models.RecipientNotifications{})
		if notificationFilter.RecipientID != "" {
			tx = tx.Where("recipient_id = ?", recipient.ID)
		}
		if notificationFilter.ChannelName != "" {
			tx = tx.Where("channel_name = ?", notificationFilter.ChannelName)
		}
		if !notificationFilter.From.IsZero() && !notificationFilter.To.IsZero() {
			tx = tx.Where("updated_at BETWEEN ? AND ?", notificationFilter.From, notificationFilter.To)
		}
		results := []successFailedData{}

		res = tx.Select("status, count(*) as count").Where("notification_id = ?", value.NotificationID).Group("status").Scan(&results)
		if res.Error != nil {
			return nil, res.Error
		}

		idx := 0
		for _, val := range results {
			if idx == 1 {
				if val.Status == constants.Success {
					graphData.Successful[i] = graphData.Successful[i] + val.Count
					graphData.Total[i] = graphData.Total[i] + val.Count
				} else if val.Status == constants.Failure {
					graphData.Failed[i] = graphData.Failed[i] + val.Count
					graphData.Total[i] = graphData.Total[i] + val.Count
				}
			} else {
				graphData.Total = append(graphData.Total, val.Count)
				if val.Status == constants.Success {
					graphData.Successful = append(graphData.Successful, val.Count)
					graphData.Failed = append(graphData.Failed, 0)
				} else if val.Status == constants.Failure {
					graphData.Successful = append(graphData.Successful, 0)
					graphData.Failed = append(graphData.Failed, val.Count)
				}
				idx++
			}
		}
	}
	return &graphData, nil
}
