package recipients

import (
	"fmt"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/jinzhu/gorm"
)

// AddUpdateRecipientWithID creates or updates the recipient given recipient details
func AddUpdateRecipientWithID(recipientInfo *serializers.RecipientInfo) (int, error) {
	dbG := db.Get()

	var lastRecipient models.Recipient
	var lastID uint
	err := dbG.Last(&lastRecipient).Error
	if err == gorm.ErrRecordNotFound {
		lastID = 0
	} else if err != nil {
		return http.StatusInternalServerError, err
	} else {
		lastID = lastRecipient.ID
	}
	if lastID+1 < uint(recipientInfo.ID) {
		return http.StatusBadRequest, (fmt.Errorf("There cannot be any gap in ids as in %v can have max value as %v", recipientInfo.ID, lastID+1))
	}
	var recipient models.Recipient
	err = dbG.Model(&models.Recipient{}).Where(recipientInfo.ID).FirstOrCreate(&recipient).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}
	recipient.PreferredChannelID = recipientInfo.PreferredChannelID
	recipient.PushToken = recipientInfo.PushToken
	recipient.WebToken = recipientInfo.WebToken
	recipient.Email = strings.ToLower(recipientInfo.Email)

	return http.StatusInternalServerError, dbG.Save(&recipient).Error
}
