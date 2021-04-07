package migrations

import (
	"database/sql"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	models "code.jtg.tools/ayush.singhal/notifications-microservice/db/base_models"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddResetTokenField, downAddResetTokenField)
}

type user struct {
	ResetToken string
}

func upAddResetTokenField(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	dbG := db.Get()

	err := dbG.AutoMigrate(&user{}).Error
	if err != nil {
		return err
	}
	return nil
}

func downAddResetTokenField(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	dbG := db.Get()

	err := dbG.Model(&models.User{}).DropColumn("reset_token").Error
	if err != nil {
		return err
	}
	return nil
}
