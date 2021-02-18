package filter

// Channel struct is the serializer for channel filter
type Channel struct {
	ID       uint   `form:"channel_id"`
	Name     string `form:"name"`
	Type     int    `form:"type"`
	Priority int    `form:"priority"`
}
