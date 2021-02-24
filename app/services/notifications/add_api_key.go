package notifications

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// PostAPIKey creates a new API Key
func PostAPIKey() (string, error) {
	var organisation models.Organisation
	err := db.Get().First(&organisation).Error
	if err != gorm.ErrRecordNotFound {
		er := db.Get().Delete(&organisation).Error
		if er != nil {
			return "", errors.Wrap(err, "Delete Previous API Key Error")
		}
	}
	apiKey := hash.GenerateSecureToken(constants.APIKeyLength)
	apiLast := apiKey[len(apiKey)-8:]

	organisation = models.Organisation{}

	organisation.APIKey = hash.Message(apiKey, configuration.GetResp().APIHash)
	organisation.APILast = apiLast
	err = db.Get().Create(&organisation).Error
	if err != nil {
		return "", errors.Wrap(err, "Creating new API Key error")
	}
	return apiKey, nil
}
