package apimessage

import (
	"strings"
	"sync"
)

// OpenAPI is the serializer for showing status messages of all the notifications
type OpenAPI struct {
	NotificationStatus          map[string]OpenAPIChannel `json:"notification_status,omitempty"`
	RecipientIDIncorrect        []string                  `json:"recipient_id_incorrect,omitempty"`
	PreferredChannelTypeDeleted []string                  `json:"Preferred_channel_deleted,omitempty"`
	RepeatedID                  []string                  `json:"recipient_id_repeated,omitempty"`
}

// OpenAPIChannel is the serializer to show success and failure recipient IDs per channel
type OpenAPIChannel struct {
	Success []string `json:"success,omitempty"`
	Failure []string `json:"failure,omitempty"`
}

const (
	// PreferredChannelTypeDeleted is the first enum member for message options
	PreferredChannelTypeDeleted = iota + 1
	// Success is the second enum member for message options
	Success
	// Failure is the second enum member for message options
	Failure
)

// Message is the struct for OpenAPI message channel
type Message struct {
	ID          string
	Option      int
	ChannelName string
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

// AddStatus is the function used to add some message to the openAPI based on messageOption
func (openAPI *OpenAPI) AddStatus(stopChan chan bool, messageChan chan Message, mainWaitGroup *sync.WaitGroup) {

	defer mainWaitGroup.Done()

	var message Message
	for ok := true; ok; {
		select {
		case <-stopChan:
			return
		case message, ok = <-messageChan:
			if !ok {
				break
			}
			if message.Option == PreferredChannelTypeDeleted {
				openAPI.PreferredChannelTypeDeleted = append(openAPI.PreferredChannelTypeDeleted, message.ID)
			} else if message.Option == Success {
				openAPI.AddRecipientID(message.ID, message.ChannelName, true)
			} else if message.Option == Failure {
				openAPI.AddRecipientID(message.ID, message.ChannelName, false)
			}
		default:
		}
	}
}
