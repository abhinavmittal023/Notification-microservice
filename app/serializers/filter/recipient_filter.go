package filter

// Recipient struct is the serializer for recipient filter
type Recipient struct {
	RecipientUUID      string `form:"recipient_uuid"`
	Email              int    `form:"email"`
	PushToken          int    `form:"push_token"`
	WebToken           int    `form:"web_token"`
	PreferredChannelID uint64 `form:"preferred_channel_id"`
}
