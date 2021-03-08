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

// GetLastChannel function gets the information of last record of the table
func GetLastChannel() (*models.Channel, error) {
	var channel models.Channel
	res := db.Get().Last(&channel)
	return &channel, res.Error
}

// GetFirstChannel function gets the information of first record of the table
func GetFirstChannel() (*models.Channel, error) {
	var channel models.Channel
	res := db.Get().First(&channel)
	return &channel, res.Error
}

// GetNextChannelfromID function gives the details of the next channel and returns record not found
// if the record is the last one
func GetNextChannelfromID(channelID uint64) (*models.Channel, error) {
	var channelDetails models.Channel
	res := db.Get().Model(&models.Channel{}).Where("id > ?", channelID).First(&channelDetails)
	return &channelDetails, res.Error
}

// GetPreviousChannelfromID function gives the details of the previous channel and returns record not found
// if the record is the first one
func GetPreviousChannelfromID(channelID uint64) (*models.Channel, error) {
	var channelDetails models.Channel
	res := db.Get().Model(&models.Channel{}).Where("id < ?", channelID).Last(&channelDetails)
	return &channelDetails, res.Error
}

// GetChannelWithType gets the channel of specified type from the database, and returns error/nil
func GetChannelWithType(channelType uint) (*models.Channel, error) {
	var channel models.Channel
	res := db.Get().Model(&models.Channel{}).Where("type = ?", channelType).First(&channel)
	return &channel, res.Error
}

// GetChannelsWithPriorityLessThan gets the channels with priority less than or equal to specified from the database, and returns error/nil
func GetChannelsWithPriorityLessThan(priority uint) (*[]models.Channel, error) {
	var channel []models.Channel
	res := db.Get().Model(&models.Channel{}).Where("priority > ?", priority-1).Find(&channel)
	return &channel, res.Error
}

// GetChannelWithName gets the channel of specified name from the database, and returns error/nil
func GetChannelWithName(channelName string) (*models.Channel, error) {
	var channel models.Channel
	res := db.Get().Model(&models.Channel{}).Where("name = ?", channelName).First(&channel)
	return &channel, res.Error
}

// GetAllChannels gets all the channels from the database and returns []models.Channel,err
func GetAllChannels(pagination *serializers.Pagination, channelFilter *filter.Channel) ([]models.Channel, error) {

	var channels []models.Channel
	dbG := db.Get()
	tx := dbG.Model(&models.Channel{})

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

// GetAllChannelsCount gets all the channels count from the database and returns records count,err
func GetAllChannelsCount(channelFilter *filter.Channel) (int64, error) {

	dbG := db.Get()
	tx := dbG.Model(&models.Channel{})

	if channelFilter.Name != "" {
		tx = tx.Where("name = ?", channelFilter.Name)
	}
	if channelFilter.Type != 0 {
		tx = tx.Where("type = ?", channelFilter.Type)
	}
	if channelFilter.Priority != 0 {
		tx = tx.Where("priority = ?", channelFilter.Priority)
	}

	var count int64
	res := tx.Count(&count)
	return count, res.Error
}
