package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAlterOrganisation, downAlterOrganisation)
}

func upAlterOrganisation(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func downAlterOrganisation(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
