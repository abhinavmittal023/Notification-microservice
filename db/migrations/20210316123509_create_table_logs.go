package migrations

import (
	"database/sql"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableLogs, downCreateTableLogs)
}

func upCreateTableLogs(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	dbG := db.Get()
	err := dbG.AutoMigrate(&models.Logs{}).Error
	if err != nil {
		return errors.Wrap(err, "Unable to create log model")
	}
	return nil
}

func downCreateTableLogs(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	dbG := db.Get()
	err := dbG.DropTable(&models.Logs{}).Error
	if err != nil {
		return errors.Wrap(err, "Unable to remove log model")
	}
	return nil
}
