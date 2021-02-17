package recipients

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/jinzhu/gorm"
)

// AddUpdateRecipients reads each recipient info and updates the database after validation
func AddUpdateRecipients(recipientRecords *[]serializers.RecipientInfo) (int, *[]serializers.ErrorInfo) {

	var errors []serializers.ErrorInfo
	tx := db.Get().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, recipientRecord := range *recipientRecords {

		if recipientRecord.Email != "" {
			er := serializers.EmailRegexCheck(recipientRecord.Email)

			if er == "internal_server_error" {
				log.Println("Error Due to Regex")
				errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("Internal Server Error at id %v", recipientRecord.ID)})
				tx.Rollback()
				return http.StatusInternalServerError, &errors
			}
			if er == "bad_request" {
				errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("Email of ID %v is invalid", recipientRecord.ID)})
				continue
			}
		}
		status, err := AddUpdateRecipientWithID(&recipientRecord, tx)
		if err != nil {
			if status == http.StatusInternalServerError {
				errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("Internal Server Error at id %v", recipientRecord.ID)})
				tx.Rollback()
				return http.StatusInternalServerError, &errors
			}
			errors = append(errors, serializers.ErrorInfo{Error: err.Error()})
		}
	}
	if len(errors) > 0 {
		tx.Rollback()
		return http.StatusBadRequest, &errors
	}
	err := tx.Commit().Error
	if err != nil {
		errors = append(errors, serializers.ErrorInfo{Error: err.Error()})
		return http.StatusInternalServerError, &errors
	}
	return http.StatusOK, nil
}

// AddUpdateRecipientWithID creates or updates the recipient given recipient details
func AddUpdateRecipientWithID(recipientInfo *serializers.RecipientInfo, tx *gorm.DB) (int, error) {

	var lastRecipient models.Recipient
	var lastID uint
	err := tx.Last(&lastRecipient).Error
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
	err = tx.Model(&models.Recipient{}).Where(recipientInfo.ID).FirstOrCreate(&recipient).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}
	recipient.PreferredChannelID = recipientInfo.PreferredChannelID
	recipient.PushToken = recipientInfo.PushToken
	recipient.WebToken = recipientInfo.WebToken
	recipient.Email = strings.ToLower(recipientInfo.Email)

	err = tx.Save(&recipient).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
