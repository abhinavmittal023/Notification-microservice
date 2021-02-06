package migrations

import (
	"database/sql"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTables, downCreateTables)
}

func upCreateTables(tx *sql.Tx) error {
	dbG := db.Get()

	err := dbG.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}

	err = dbG.Model(&models.User{}).AddUniqueIndex("email_date", "email", "deleted_at").Error
	if err != nil {
		return err
	}

	err = dbG.AutoMigrate(&models.Recipient{}).Error
	if err != nil {
		return err
	}

	err = dbG.AutoMigrate(&models.Notification{}).Error
	if err != nil {
		return err
	}

	err = dbG.AutoMigrate(&models.Organisation{}).Error
	if err != nil {
		return err
	}


	err = dbG.AutoMigrate(&models.Channel{}).Error
	if err != nil {
		return err
	}

	err = dbG.AutoMigrate(&models.RecipientNotifications{}).Error

	return err
}

func downCreateTables(tx *sql.Tx) error {
	dbG := db.Get()

	err := dbG.Model(&models.User{}).RemoveIndex("email_date").Error
	if err != nil {
		return err
	}
	
	err = dbG.DropTable(&models.User{}).Error
	if err != nil {
		return err
	}

	err = dbG.DropTable(&models.Recipient{}).Error
	if err != nil {
		return err
	}

	err = dbG.DropTable(&models.Notification{}).Error
	if err != nil {
		return err
	}

	err = dbG.DropTable(&models.Organisation{}).Error
	if err != nil {
		return err
	}


	err = dbG.DropTable(&models.Channel{}).Error
	if err != nil {
		return err
	}

	err = dbG.DropTable(&models.RecipientNotifications{}).Error

	return err
}
