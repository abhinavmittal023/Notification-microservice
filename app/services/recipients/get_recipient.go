package recipients

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
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
func GetAllRecipients(pagination serializers.Pagination) ([]models.Recipient, error) {
	var recipients []models.Recipient
	dbg := db.Get()
	res := dbg.Offset(pagination.Offset).Limit(pagination.Limit).Find(&recipients)
	return recipients, res.Error
}
