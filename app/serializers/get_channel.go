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

// WebConfig serializer to get config info for web channel
type WebConfig struct {
	ServerKey	string	`json:"server_key" binding:"required"`
}

// PushConfig serializer to get config info for push channel
type PushConfig struct {
	ServerKey string `json:"server_key" binding:"required"`
}

// EmailConfig serializer to get config info for email channel
type EmailConfig struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	SMTPHost string `json:"smtp_host" binding:"required"`
	SMTPPort string `json:"smtp_port" binding:"required"`
}

// ChannelListResponse serializer for channel list response
type ChannelListResponse struct {
	RecordsAffected		int64	`json:"records_count"`
	ChannelInfo			[]ChannelInfo 	`json:"channels"`
}

// ChannelInfoToChannelModel function copies the data from channel serializer to channel model
func ChannelInfoToChannelModel(channelInfo *ChannelInfo, channelModel *models.Channel) {
	channelModel.Name = strings.ToLower(channelInfo.Name)
	channelModel.ShortDescription = strings.ToLower(channelInfo.ShortDescription)
	channelModel.Type = int(channelInfo.Type)
	channelModel.Priority = channelInfo.Priority
	channelModel.Configuration = channelInfo.Configuration
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


