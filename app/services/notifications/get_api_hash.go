package notifications

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// GetAPIHash returns the API Hash from the database
func GetAPIHash() (string, error) {
	var organisation models.Organisation
	err := db.Get().First(&organisation).Error
	if err == gorm.ErrRecordNotFound {
		return "", errors.Wrap(err, "No API Key found")
	} else if err != nil {
		return "", errors.Wrap(err, "Get API Hash error")
	}
	return organisation.APIKey, nil
}
