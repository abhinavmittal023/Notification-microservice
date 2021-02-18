package channels

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// GetChannelWithID gets the channel with specified ID from the database, and returns error/nil
func GetChannelWithID(id uint) (*models.Channel, error) {
	var channel models.Channel
	res := db.Get().First(&channel, id)
	return &channel, res.Error
}

// GetAllChannels gets all the channels from the database and returns []models.Channel,err
func GetAllChannels(pagination serializers.Pagination) ([]models.Channel, error) {
	var channels []models.Channel
	res := db.Get().Offset(pagination.Offset).Limit(pagination.Limit).Find(&channels)
	return channels, res.Error
}
