package recipients

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
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

	for index, recipientRecord := range *recipientRecords {
		invalid := false

		if recipientRecord.RecipientUUID == "" {
			errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("Recipient ID should not be empty at %v", index+2)})
			invalid = true
		}

		if recipientRecord.Email != "" {
			status, err := serializers.EmailRegexCheck(recipientRecord.Email)

			if err != nil {
				if status == http.StatusInternalServerError {
					log.Println("Error Due to Regex", err.Error())
					errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("Internal Server Error at %v", index+2)})
					tx.Rollback()
					return http.StatusInternalServerError, &errors
				}
				if status == http.StatusBadRequest {
					errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("Email at %v is invalid", index+2)})
					continue
				}
			}
		}

		if recipientRecord.PreferredChannelID != 0 {
			channel, err := channels.GetChannelWithID(uint(recipientRecord.PreferredChannelID))
			if err == gorm.ErrRecordNotFound {
				errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("PreferredChannelID at %v is not in the database", index+2)})
				invalid = true
			} else if err != nil {
				log.Println(err.Error())
				errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("Internal Server Error at %v", index+2)})
				return http.StatusInternalServerError, &errors
			}
			channelType := constants.ChannelType(uint(channel.Type))

			if channelType == "Email" && recipientRecord.Email == "" {
				errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("PreferredChannel %s cannot be empty at %v", channelType, index+2)})
				invalid = true
			} else if channelType == "Web" && recipientRecord.WebToken == "" {
				errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("PreferredChannel %s cannot be empty at %v", channelType, index+2)})
				invalid = true
			} else if channelType == "Push" && recipientRecord.PushToken == "" {
				errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("PreferredChannel %s cannot be empty at %v", channelType, index+2)})
				invalid = true
			} else if channelType == "" {
				errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("PreferredChannel at %v is invalid", index+2)})
				invalid = true
			}
		}

		if invalid {
			continue
		}

		status, err := AddUpdateRecipientWithID(&recipientRecord, tx)
		if err != nil {
			if status == http.StatusInternalServerError {
				errors = append(errors, serializers.ErrorInfo{Error: fmt.Sprintf("Internal Server Error id %v", index+2)})
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

	var recipient models.Recipient
	err := tx.Model(&models.Recipient{}).Where("recipient_uuid = ?", recipientInfo.RecipientUUID).FirstOrCreate(&recipient).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}
	recipient.RecipientUUID = recipientInfo.RecipientUUID
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
