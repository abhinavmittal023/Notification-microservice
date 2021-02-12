package migrations

import (
	"database/sql"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTables, downCreateTables)
}

func upCreateTables(tx *sql.Tx) error {
	dbG := db.Get()

	err := dbG.AutoMigrate(&models.User{}, &models.Recipient{}, &models.Notification{}, &models.Organisation{}, &models.Channel{}, &models.RecipientNotifications{}).Error
	if err != nil {
		return errors.Wrap(err, "Unable to create models")
	}

	err = dbG.Model(&models.User{}).AddUniqueIndex("email_date", "email", "deleted_at").Error
	if err != nil {
		return errors.Wrap(err, "Unable to add unique index to user model")
	}

	return nil
}

func downCreateTables(tx *sql.Tx) error {
	dbG := db.Get()

	err := dbG.Model(&models.User{}).RemoveIndex("email_date").Error
	if err != nil {
		return errors.Wrap(err, "Unable to remove unique index from user model")
	}

	err = dbG.DropTable(&models.User{}, &models.Recipient{}, &models.Notification{}, &models.Organisation{}, &models.Channel{}, &models.RecipientNotifications{}).Error
	if err != nil {
		return errors.Wrap(err, "Unable to remove models")
	}

	return nil
}
