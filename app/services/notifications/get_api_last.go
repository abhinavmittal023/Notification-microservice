package notifications

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// GetAPILast gets the APIKey's last 8 character from the database, and returns error/nil
func GetAPILast() (string, error) {
	var organisation models.Organisation
	res := db.Get().First(&organisation)
	return organisation.APILast, res.Error
}
