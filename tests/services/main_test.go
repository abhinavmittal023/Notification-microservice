package services

import (
	"log"
	"os"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
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

	os.Exit(m.Run())
}
