package serializers

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// RecipientInfo serializer to bind request data
type RecipientInfo struct {
	ID                   uint64 `json:"recipient_id"`
	RecipientUUID        string `json:"recipient_uuid"`
	Email                string `json:"email,omitempty"`
	PushToken            string `json:"push_token,omitempty"`
	WebToken             string `json:"web_token,omitempty"`
	PreferredChannelType uint   `json:"preferred_channel_type,omitempty"`
	ChannelType          uint   `json:"channel_type,omitempty"`
}

// RecipientModelToRecipientInfo converts the Recipient model to RecipientInfo struct
func RecipientModelToRecipientInfo(info *RecipientInfo, recipient *models.Recipient) {
	info.ID = uint64(recipient.ID)
	info.RecipientUUID = recipient.RecipientUUID
	info.Email = recipient.Email
	info.PushToken = recipient.PushToken
	info.WebToken = recipient.WebToken
	info.PreferredChannelType = recipient.PreferredChannelType
}
