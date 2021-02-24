package serializers

// SendNotifications struct holds information about APIKey and Notifications
type SendNotifications struct {
	APIKey        string        `json:"api_key" binding:"required"`
	Notifications Notifications `json:"notifications" binding:"required"`
}

// Notifications serializer holds the information about notifications
type Notifications struct {
	Recipients []string `json:"recipients" binding:"required"`
	Priority   int      `json:"priority" binding:"required"`
	Title      string   `json:"title" binding:"required"`
	Body       string   `json:"body" binding:"required"`
}
