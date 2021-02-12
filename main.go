package main

import (
	"log"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
)

func main() {
	err := configuration.Init()
	if err != nil {
		log.Println(err.Error())
		return
	}

	dbstring := configuration.GetResp().Database.Dbstring
	err = db.Init(dbstring)
	if err != nil {
		log.Println("Error Connecting to the Database")
		return
	}

	dbG := db.Get()
	defer dbG.Close()
	err = api.InitServer()
	if err != nil {
		log.Println("Error connecting to Server")
		return
	}
}
