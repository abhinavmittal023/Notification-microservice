package filter

import "strings"

// Channel struct is the serializer for channel filter
type Channel struct {
	Name     string `form:"name"`
	Type     int    `form:"type"`
	Priority int    `form:"priority"`
}

// ConvertChannelStringToLower converts string values to lower case
func ConvertChannelStringToLower(channel *Channel) {
	channel.Name = strings.ToLower(channel.Name)
}
