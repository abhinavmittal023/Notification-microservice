package serializers

import "code.jtg.tools/ayush.singhal/notifications-microservice/db/models"

// ChannelInfo serializer to get and show channel information
type ChannelInfo struct {
	ID               uint   `json:"id"`
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description,omitempty"`
	Type             int    `json:"type" binding:"required"`
	Priority         int    `json:"priority" binding:"required"`
	Configuration    string `json:"configuration,omitempty"`
}

// ChannelInfoToChannelModel function copies the data from channel serializer to channel model
func ChannelInfoToChannelModel(channelInfo *ChannelInfo, channelModel *models.Channel) {
	channelModel.Name = channelInfo.Name
	channelModel.ShortDescription = channelInfo.ShortDescription
	channelModel.Type = channelInfo.Type
	channelModel.Priority = channelInfo.Priority
	channelModel.Configuration = channelInfo.Configuration
}

// ChannelModelToChannelInfo function copies the data from channel model to channel serializer
func ChannelModelToChannelInfo(channelInfo *ChannelInfo, channelModel *models.Channel) {
	channelInfo.ID = channelModel.ID
	channelInfo.Name = channelModel.Name
	channelInfo.ShortDescription = channelModel.ShortDescription
	channelInfo.Type = channelModel.Type
	channelInfo.Priority = channelModel.Priority
	channelInfo.Configuration = channelModel.Configuration
}
