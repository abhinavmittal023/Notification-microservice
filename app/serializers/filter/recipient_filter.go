package filter

// Recipient struct is the serializer for recipient filter
type Recipient struct {
	RecipientID          string `form:"recipient_id"`
	Email                int    `form:"email"`
	PushToken            int    `form:"push_token"`
	WebToken             int    `form:"web_token"`
	PreferredChannelType uint64 `form:"preferred_channel_type"`
}
