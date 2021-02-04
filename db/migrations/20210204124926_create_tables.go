package migrations

import (
	"database/sql"
	"log"

	// "code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTables, downCreateTables)
}

func upCreateTables(tx *sql.Tx) error {
	log.Println("Up Migration")
	// This code is executed when the migration is applied.
	//dbG := db.Get()
	return nil
}

func downCreateTables(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	//dbG := db.Get()
	return nil
}
