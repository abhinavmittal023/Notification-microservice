package serializers

import (
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// ChannelInfo serializer to get and show channel information
type ChannelInfo struct {
	ID               uint   `json:"id"`
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description,omitempty"`
	Type             uint   `json:"type" binding:"required"`
	Priority         int    `json:"priority" binding:"required"`
	Configuration    string `json:"configuration,omitempty"`
}

// ChannelInfoToChannelModel function copies the data from channel serializer to channel model
func ChannelInfoToChannelModel(channelInfo *ChannelInfo, channelModel *models.Channel) {
	channelModel.Name = strings.ToLower(channelInfo.Name)
	channelModel.ShortDescription = strings.ToLower(channelInfo.ShortDescription)
	channelModel.Type = int(channelInfo.Type)
	channelModel.Priority = channelInfo.Priority
	channelModel.Configuration = strings.ToLower(channelInfo.Configuration)
}

// ChannelModelToChannelInfo function copies the data from channel model to channel serializer
func ChannelModelToChannelInfo(channelInfo *ChannelInfo, channelModel *models.Channel) {
	channelInfo.ID = channelModel.ID
	channelInfo.Name = channelModel.Name
	channelInfo.ShortDescription = channelModel.ShortDescription
	channelInfo.Type = uint(channelModel.Type)
	channelInfo.Priority = channelModel.Priority
	channelInfo.Configuration = channelModel.Configuration
}
