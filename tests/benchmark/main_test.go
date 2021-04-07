package benchmark

import (
	"log"
	"os"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

func TestMain(m *testing.M) {
	err := configuration.Init()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	dbstring := configuration.GetResp().Database.TestDbstring
	err = db.Init(dbstring)
	if err != nil {
		log.Fatalln("Error Connecting to the Database")
		return
	}

	dbG := db.Get()
	defer dbG.Close()
	if err := RefreshAllTables(); err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// RefreshAllTables function clears the records in the database
func RefreshAllTables() error {

	dbG := db.Get()

	err := dbG.DropTableIfExists(&models.User{}, &models.Recipient{}, &models.Notification{}, &models.Organisation{}, &models.Channel{}, &models.RecipientNotifications{}, &models.Logs{}).Error
	if err != nil {
		return err
	}

	err = dbG.AutoMigrate(&models.User{}, &models.Recipient{}, &models.Notification{}, &models.Organisation{}, &models.Channel{}, &models.RecipientNotifications{}, &models.Logs{}).Error
	if err != nil {
		return err
	}

	err = dbG.Model(&models.User{}).AddUniqueIndex("email_date", "email", "deleted_at").Error
	if err != nil {
		return err
	}

	err = dbG.Model(&models.Channel{}).AddUniqueIndex("type_date", "type", "deleted_at").Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil

}
