package serializers

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// RecipientInfo serializer to bind request data
type RecipientInfo struct {
	ID                 uint64 `json:"recipient_id"`
	Email              string `json:"email,omitempty"`
	PushToken          string `json:"push_token,omitempty"`
	WebToken           string `json:"web_token,omitempty"`
	PreferredChannelID uint64 `json:"preferred_channel_id,omitempty"`
}

// RecipientModelToRecipientInfo converts the Recipient model to RecipientInfo struct
func RecipientModelToRecipientInfo(info *RecipientInfo, recipient *models.Recipient) {
	info.ID = uint64(recipient.ID)
	info.Email = recipient.Email
	info.PushToken = recipient.PushToken
	info.WebToken = recipient.WebToken
	info.PreferredChannelID = recipient.PreferredChannelID
}
