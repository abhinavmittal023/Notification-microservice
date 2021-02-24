package recipients

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// GetRecipientWithID gets the recipient with specified ID from the database, and returns error/nil
func GetRecipientWithID(recipientID uint64) (*models.Recipient, error) {
	var recipient models.Recipient
	res := db.Get().First(&recipient, recipientID)
	return &recipient, res.Error
}

// GetAllRecipients gets all Recipients from the database and returns []models.Recipient,err
func GetAllRecipients(pagination *serializers.Pagination, recipientFilter *filter.Recipient) ([]models.Recipient, error) {

	var recipients []models.Recipient
	dbg := db.Get()
	tx := dbg.Model(&models.Recipient{})

	if recipientFilter.RecipientID != "" {
		tx = tx.Where("recipient_id = ?", recipientFilter.RecipientID)
	}
	if recipientFilter.PreferredChannelType != 0 {
		tx = tx.Where("preferred_channel_type = ?", recipientFilter.PreferredChannelType)
	}
	if recipientFilter.Email > 0 {
		tx = tx.Not("email", "")
	} else if recipientFilter.Email < 0 {
		tx = tx.Where("email = ?", "")
	}
	if recipientFilter.PushToken > 0 {
		tx = tx.Not("push_token", "")
	} else if recipientFilter.PushToken < 0 {
		tx = tx.Where("push_token = ?", "")
	}
	if recipientFilter.WebToken > 0 {
		tx = tx.Not("web_token", "")
	} else if recipientFilter.WebToken < 0 {
		tx = tx.Where("web_token = ?", "")
	}

	res := tx.Offset(pagination.Offset).Limit(pagination.Limit).Find(&recipients)
	return recipients, res.Error
}

// GetRecipientWithRecipientID gets the recipient with specified ID from the database, and returns error/nil
func GetRecipientWithRecipientID(recipientID string) (*models.Recipient, error) {
	var recipient models.Recipient
	res := db.Get().Model(&recipient).Where("recipient_id = ?", recipientID).First(&recipient)
	return &recipient, res.Error
}
