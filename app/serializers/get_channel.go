package serializers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// ChannelInfo serializer to get and show channel information
type ChannelInfo struct {
	ID               uint   `json:"id"`
	Name             string `json:"name" binding:"required,max=100"`
	ShortDescription string `json:"short_description,omitempty" binding:"max=500"`
	Type             int    `json:"type" binding:"required"`
	Priority         int    `json:"priority" binding:"required"`
	Configuration    string `json:"configuration," binding:"required"`
	RecipientsCount  uint64 `json:"recipients"`
}

// WebConfig serializer to get config info for web channel
type WebConfig struct {
	ServerKey string `json:"server_key" binding:"required"`
}

// PushConfig serializer to get config info for push channel
type PushConfig struct {
	ServerKey string `json:"server_key" binding:"required"`
}

// EmailConfig serializer to get config info for email channel
type EmailConfig struct {
	Email    string `json:"email" binding:"required"`
	From     string `json:"from"`
	Password string `json:"password" binding:"required"`
	SMTPHost string `json:"smtp_host" binding:"required"`
	SMTPPort string `json:"smtp_port" binding:"required"`
}

// ChannelListResponse serializer for channel list response
type ChannelListResponse struct {
	RecordsAffected int64         `json:"records_count"`
	ChannelInfo     []ChannelInfo `json:"channels"`
}

// ChannelConfigValidation function checks if configuration details are ok and can be deserialized
func ChannelConfigValidation(channelInfo *ChannelInfo) error {
	if channelInfo.Type == (constants.ChannelIntType()[0]) {
		var config EmailConfig
		err := json.Unmarshal([]byte(channelInfo.Configuration), &config)
		if err != nil || config.Email == "" || config.SMTPHost == "" || config.SMTPPort == "" {
			return fmt.Errorf(constants.Errors().InvalidJSON)
		}
		_, err = EmailRegexCheck(config.Email)
		if err != nil {
			return err
		}
		_, err = HostRegexCheck(config.SMTPHost)
		if err != nil {
			return err
		}
		_, err = PortRegexCheck(config.SMTPPort)
		if err != nil {
			return err
		}
		return nil
	} else if channelInfo.Type == (constants.ChannelIntType()[1]) {
		var config PushConfig
		err := json.Unmarshal([]byte(channelInfo.Configuration), &config)
		if err != nil || config.ServerKey == "" {
			return fmt.Errorf(constants.Errors().InvalidJSON)
		}
		return nil
	}
	var config WebConfig

	err := json.Unmarshal([]byte(channelInfo.Configuration), &config)
	if err != nil || config.ServerKey == "" {
		return fmt.Errorf(constants.Errors().InvalidJSON)
	}
	return nil
}

// HostRegexCheck checks for Host in valid format
func HostRegexCheck(host string) (int, error) {
	match, err := regexp.MatchString(constants.HostRegex, host)
	if err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, fmt.Errorf(constants.Errors().InternalError)
	}
	if !match {
		return http.StatusBadRequest, fmt.Errorf(constants.Errors().InvalidHost)
	}
	return http.StatusOK, nil
}

// PortRegexCheck checks for Port in valid format
func PortRegexCheck(port string) (int, error) {
	match, err := regexp.MatchString(constants.PortRegex, port)
	if err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, fmt.Errorf(constants.Errors().InternalError)
	}
	if !match {
		return http.StatusBadRequest, fmt.Errorf(constants.Errors().InvalidPort)
	}
	return http.StatusOK, nil
}

// ChannelInfoToChannelModel function copies the data from channel serializer to channel model
func ChannelInfoToChannelModel(channelInfo *ChannelInfo, channelModel *models.Channel) {
	channelModel.Name = strings.ToLower(channelInfo.Name)
	channelModel.ShortDescription = strings.ToLower(channelInfo.ShortDescription)
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

// ChannelModelToChannelInfoWithRecipientCount function copies the data from channel model to channel serializer
func ChannelModelToChannelInfoWithRecipientCount(channelInfo *ChannelInfo, channelModel *models.Channel, count uint64) {
	channelInfo.ID = channelModel.ID
	channelInfo.Name = channelModel.Name
	channelInfo.ShortDescription = channelModel.ShortDescription
	channelInfo.Type = channelModel.Type
	channelInfo.Priority = channelModel.Priority
	channelInfo.Configuration = channelModel.Configuration
	channelInfo.RecipientsCount = count
}
