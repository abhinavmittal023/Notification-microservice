package serializers

// NotificationsResponse serializer to send notifications list
type NotificationsResponse struct {
	RecordsAffected int64         `json:"records_count"`
	NotificationsInfo	[]NotificationsInfo `json:"notifications"`
}

// NotificationsInfo serializer to get and show notifications information
type NotificationsInfo struct {
	Priority             int                    `json:"priority"`
	Title                string                 `json:"title"`
	Body                 string                 `json:"body"`
	NotificationChannels []NotificationChannels `json:"channels"`
}

// NotificationChannels serializer to get and show channels information
type NotificationChannels struct {
	ChannelName string `json:"name"`
	Successful  uint64 `json:"successful"`
	Failure     uint64 `json:"failure"`
	Total       uint64 `json:"total"`
}
