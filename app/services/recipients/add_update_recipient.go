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
func AddUpdateRecipients(recipientRecords *[]serializers.RecipientInfo) (int, *serializers.ErrorInfo) {

	var errors serializers.ErrorInfo
	errors.Error = make(map[int][]string)
	var index int
	errorFlag := false
	tx := db.Get().Begin()
	defer func() {
		if r := recover(); r != nil {
			var errorMap []string
			tx.Rollback()
			log.Println(r)
			errorMap = append(errorMap, constants.Errors().InternalError)
			errors.Error[index+2] = errorMap
		}
	}()

	for index, recipientRecord := range *recipientRecords {
		invalid := false
		var errorMap []string

		if recipientRecord.RecipientID == "" {
			errorMap = append(errorMap, "Recipient ID should not be empty")
			invalid = true
		}

		if len(recipientRecord.PushToken) > constants.MaxPushToken {
			errorMap = append(errorMap, fmt.Sprintf("Push Token cannot be more than %v characters long", constants.MaxPushToken))
			invalid = true
		}

		if len(recipientRecord.WebToken) > constants.MaxWebToken {
			errorMap = append(errorMap, fmt.Sprintf("web Token cannot be more than %v characters long", constants.MaxWebToken))
			invalid = true
		}

		if recipientRecord.Email != "" {
			if len(recipientRecord.Email) > constants.MaxEmail {
				errorMap = append(errorMap, fmt.Sprintf("Email cannot be more than %v characters long", constants.MaxEmail))
				invalid = true
			}
			status, err := serializers.EmailRegexCheck(recipientRecord.Email)

			if err != nil {
				if status == http.StatusInternalServerError {
					log.Println("Error Due to Regex", err.Error())
					errorMap = append(errorMap, constants.Errors().InternalError)
					errors.Error[index+2] = errorMap
					tx.Rollback()
					return http.StatusInternalServerError, &errors
				}
				if status == http.StatusBadRequest {
					errorMap = append(errorMap, constants.Errors().InvalidEmail)
					invalid = true
				}
			}
		}

		if recipientRecord.ChannelType != 0 {
			_, err := channels.GetChannelWithType(uint(recipientRecord.ChannelType))
			if err == gorm.ErrRecordNotFound {
				errorMap = append(errorMap, fmt.Sprintf("Preferred Channel %s is not in the database", constants.ChannelType(recipientRecord.ChannelType)))
				invalid = true
			} else if err != nil {
				log.Println(err)
				errorMap = append(errorMap, constants.Errors().InternalError)
				errors.Error[index+2] = errorMap
				return http.StatusInternalServerError, &errors
			}

			channelType := constants.ChannelType(uint(recipientRecord.ChannelType))

			if channelType == "Email" && recipientRecord.Email == "" {
				errorMap = append(errorMap, fmt.Sprintf("PreferredChannel %s cannot be empty", channelType))
				invalid = true
			} else if channelType == "Web" && recipientRecord.WebToken == "" {
				errorMap = append(errorMap, fmt.Sprintf("PreferredChannel %s cannot be empty", channelType))
				invalid = true
			} else if channelType == "Push" && recipientRecord.PushToken == "" {
				errorMap = append(errorMap, fmt.Sprintf("PreferredChannel %s cannot be empty", channelType))
				invalid = true
			}
		}

		if invalid {
			errors.Error[index+2] = errorMap
			errorFlag = true
			continue
		}

		status, err := AddUpdateRecipientWithID(&recipientRecord, tx)
		if err != nil {
			if status == http.StatusInternalServerError {
				errorMap = append(errorMap, constants.Errors().InternalError)
				errors.Error[index+2] = errorMap
				tx.Rollback()
				return http.StatusInternalServerError, &errors
			}
			errorMap = append(errorMap, err.Error())
		}
	}

	if errorFlag {
		tx.Rollback()
		return http.StatusBadRequest, &errors
	}
	err := tx.Commit().Error
	if err != nil {
		var errorMap []string
		errorMap = append(errorMap, err.Error())
		errors.Error[index+2] = errorMap
		return http.StatusInternalServerError, &errors
	}
	return http.StatusOK, nil
}

// AddUpdateRecipientWithID creates or updates the recipient given recipient details
func AddUpdateRecipientWithID(recipientInfo *serializers.RecipientInfo, tx *gorm.DB) (int, error) {

	var recipient models.Recipient
	err := tx.Model(&models.Recipient{}).Where("recipient_id = ?", recipientInfo.RecipientID).FirstOrCreate(&recipient).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}
	recipient.RecipientID = recipientInfo.RecipientID
	recipient.PreferredChannelType = recipientInfo.ChannelType
	recipient.PushToken = recipientInfo.PushToken
	recipient.WebToken = recipientInfo.WebToken
	recipient.Email = strings.ToLower(recipientInfo.Email)

	err = tx.Save(&recipient).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
