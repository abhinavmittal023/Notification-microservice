package main

import (
	"log"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
)

func main() {
	err := configuration.Init()

	if err != nil {
		log.Println(err.Error())
		return
	}

	dbstring := "user=" + configuration.GetResp().Database.User + " password=" +
		configuration.GetResp().Database.Password + " dbname=" +
		configuration.GetResp().Database.DbName +
		" sslmode=" + configuration.GetResp().Database.SSLMode

	err = db.Init(dbstring)
	if err != nil {
		panic("Error Connecting to the Database")
	}

	dbG := db.Get()
	defer dbG.Close()
}
