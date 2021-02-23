package channels

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// GetChannelWithID gets the channel with specified ID from the database, and returns error/nil
func GetChannelWithID(id uint) (*models.Channel, error) {
	var channel models.Channel
	res := db.Get().First(&channel, id)
	return &channel, res.Error
}

// GetChannelWithType gets the channel of specified type from the database, and returns error/nil
func GetChannelWithType(channelType uint) (*models.Channel, error) {
	var channel models.Channel
	res := db.Get().Model(&models.Channel{}).Where("type = ?", channelType).First(&channel)
	return &channel, res.Error
}

// GetAllChannels gets all the channels from the database and returns []models.Channel,err
func GetAllChannels(pagination *serializers.Pagination, channelFilter *filter.Channel) ([]models.Channel, error) {

	var channels []models.Channel
	dbG := db.Get()
	tx := dbG.Model(&models.Channel{})

	if channelFilter.ID != 0 {
		tx = tx.Where("id = ?", channelFilter.ID)
	}
	if channelFilter.Name != "" {
		tx = tx.Where("name = ?", channelFilter.Name)
	}
	if channelFilter.Type != 0 {
		tx = tx.Where("type = ?", channelFilter.Type)
	}
	if channelFilter.Priority != 0 {
		tx = tx.Where("priority = ?", channelFilter.Priority)
	}

	res := tx.Offset(pagination.Offset).Limit(pagination.Limit).Find(&channels)
	return channels, res.Error
}
