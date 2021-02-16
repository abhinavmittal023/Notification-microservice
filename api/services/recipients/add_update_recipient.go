package recipients

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// AddUpdateRecipientWithID creates or updates the recipient given recipient details
func AddUpdateRecipientWithID(recipientInfo *serializers.RecipientInfo) error {
	dbG := db.Get()

	var recipient models.Recipient
	err := dbG.Model(&models.Recipient{}).Where(recipientInfo.ID).FirstOrCreate(&recipient).Error
	if err != nil {
		return err
	}
	recipient.ID = uint(recipientInfo.ID)
	recipient.PreferredChannelID = recipientInfo.PreferredChannelID
	recipient.PushToken = recipientInfo.PushToken
	recipient.WebToken = recipientInfo.WebToken
	recipient.Email = recipientInfo.Email

	return dbG.Save(&recipient).Error
}
