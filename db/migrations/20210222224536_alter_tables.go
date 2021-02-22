package migrations

import (
	"database/sql"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	models "code.jtg.tools/ayush.singhal/notifications-microservice/db/base_models"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAlterTables, downAlterTables)
}

func upAlterTables(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	dbG := db.Get()

	err := dbG.Model(&models.Channel{}).AddUniqueIndex("type_date", "type", "deleted_at").Error
	if err != nil {
		return errors.Wrap(err, "Unable to add unique index to channel model")
	}
	err = dbG.Exec("ALTER TABLE recipients RENAME COLUMN recipient_uuid TO recipient_id;").Error
	return err
}

func downAlterTables(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	dbG := db.Get()

	err := dbG.Model(&models.Channel{}).RemoveIndex("type_date").Error
	if err != nil {
		return errors.Wrap(err, "Unable to remove unique index from channel model")
	}
	err = dbG.Exec("ALTER TABLE recipients RENAME COLUMN recipient_id TO recipient_uuid;").Error
	return err
}
