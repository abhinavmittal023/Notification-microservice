package migrations

import (
	"database/sql"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	models "code.jtg.tools/ayush.singhal/notifications-microservice/db/base_models"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAlterTableRecipientNotifications, downAlterTableRecipientNotifications)
}

func upAlterTableRecipientNotifications(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	dbG := db.Get()
	err := dbG.Model(&models.RecipientNotifications{}).DropColumn("channel_id").Error
	if err != nil {
		return err
	}

	type recipientNotifications struct {
		ChannelName string `gorm:"default:'undefined';type:varchar(100);not null;index"`
	}

	err = dbG.AutoMigrate(&recipientNotifications{}).Error
	if err != nil {
		return err
	}
	return nil
}

func downAlterTableRecipientNotifications(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	dbG := db.Get()
	err := dbG.Model(&models.RecipientNotifications{}).DropColumn("channel_name").Error
	if err != nil {
		return err
	}

	type recipientNotifications struct {
		ChannelID uint64 `gorm:"default:0;not null;index"`
	}

	err = dbG.AutoMigrate(&recipientNotifications{}).Error
	if err != nil {
		return err
	}
	return nil
}
