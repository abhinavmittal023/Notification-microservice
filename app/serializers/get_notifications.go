package serializers

import "time"

// NotificationsResponse serializer to send notifications list
type NotificationsResponse struct {
	RecordsAffected   int64               `json:"records_count"`
	NotificationsInfo []NotificationsInfo `json:"notifications"`
}

// NotificationsInfo serializer to get and show notifications information
type NotificationsInfo struct {
	Priority             int                    `json:"priority"`
	Title                string                 `json:"title"`
	Body                 string                 `json:"body"`
	CreatedAt            time.Time              `json:"created_at"`
	NotificationChannels []NotificationChannels `json:"channels"`
	Recipients				[]Recipients				`json:"recipients"`
}

// Recipients serializer to get and show recipients information
type Recipients struct {
	RecipientID		string	`json:"recipient_id"`
	Channels []Channels		`json:"channels"`
}

// Channels serializer to get and show channels information
type Channels struct {
	ChannelName string `json:"name"`
	Status		uint64	`json:"status"`
}

// NotificationChannels serializer to get and show channels information
type NotificationChannels struct {
	ChannelName string `json:"name"`
	Successful  uint64 `json:"successful"`
	Failure     uint64 `json:"failure"`
	Total       uint64 `json:"total"`
}
