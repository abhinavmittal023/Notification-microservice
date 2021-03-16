package main

import (
	"fmt"
	"log"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	li "code.jtg.tools/ayush.singhal/notifications-microservice/shared/logwrapper"
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
		log.Println("Error Connecting to the Database", err.Error())
		return
	}

	dbG := db.Get()
	defer dbG.Close()

	f, err := li.OpenFile()
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	}
	defer f.Close()

	err = app.InitServer()
	if err != nil {
		log.Println("Error connecting to Server", err.Error())
		return
	}
}
