package apimessage

import (
	"strings"
)

// OpenAPI is the serializer for showing status messages of all the notifications
type OpenAPI struct {
	NotificationStatus          map[string]OpenAPIChannel `json:"notification_status,omitempty"`
	RecipientIDIncorrect        []string                  `json:"recipient_id_incorrect,omitempty"`
	PreferredChannelTypeDeleted []string                  `json:"Preferred_channel_deleted,omitempty"`
}

// OpenAPIChannel is the serializer to show success and failure recipient IDs per channel
type OpenAPIChannel struct {
	Success []string `json:"success,omitempty"`
	Failure []string `json:"failure,omitempty"`
}

// AddRecipientID Adds the recipient ID to the success or failure list of the struct
func (openAPI *OpenAPI) AddRecipientID(ID string, channelName string, success bool) {
	channelName = strings.ToLower(channelName)
	channelStatus, present := openAPI.NotificationStatus[channelName]
	if !present {
		var openAPIChannel OpenAPIChannel
		openAPIChannel.Success = []string{}
		openAPIChannel.Failure = []string{}
		openAPI.NotificationStatus[channelName] = openAPIChannel
	}
	if success {
		channelStatus.Success = append(channelStatus.Success, ID)
	} else {
		channelStatus.Failure = append(channelStatus.Failure, ID)
	}
	openAPI.NotificationStatus[channelName] = channelStatus
}
