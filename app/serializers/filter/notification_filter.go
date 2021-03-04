package filter

import (
	"strings"
	"time"
)

// Notification struct is the serializer for notification filter
type Notification struct {
	RecipientID string    `form:"recipient_id"`
	ChannelName string    `form:"channel_name"`
	From        time.Time `form:"from"`
	To          time.Time `form:"to"`
}

// ConvertChannelNameStringToLower converts string values to lower case
func ConvertChannelNameStringToLower(notification *Notification) {
	notification.ChannelName = strings.ToLower(notification.ChannelName)
}
