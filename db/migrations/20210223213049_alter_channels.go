package migrations

import (
	"database/sql"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	models "code.jtg.tools/ayush.singhal/notifications-microservice/db/base_models"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAlterChannels, downAlterChannels)
}

func upAlterChannels(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	dbG := db.Get()

	err := dbG.Model(&models.Channel{}).AddUniqueIndex("name_date", "name", "deleted_at").Error
	if err != nil {
		return err
	}

	return nil
}

func downAlterChannels(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	dbG := db.Get()

	err := dbG.Model(&models.Channel{}).RemoveIndex("name_date").Error
	if err != nil {
		return err
	}

	return nil
}
