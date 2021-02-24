package migrations

import (
	"database/sql"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	models "code.jtg.tools/ayush.singhal/notifications-microservice/db/base_models"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAlterOrganisation, downAlterOrganisation)
}

func upAlterOrganisation(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	dbG := db.Get()
	err := dbG.Model(&models.Organisation{}).DropColumn("name").Error
	if err != nil {
		return errors.Wrap(err, "Unable to drop column name from organisation model")
	}
	err = dbG.Exec("ALTER TABLE organisations ADD COLUMN api_last varchar(10) not null;").Error
	return err
}

func downAlterOrganisation(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	dbG := db.Get()
	err := dbG.Exec("ALTER TABLE organisations ADD COLUMN name varchar(255) not null;").Error
	if err != nil {
		return errors.Wrap(err, "Unable to add column name to organisation model")
	}
	err = dbG.Model(&models.Organisation{}).DropColumn("api_last").Error
	return err
}
