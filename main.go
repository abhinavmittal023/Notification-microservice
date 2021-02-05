package main

import (
	"log"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features"
)

func main() {

	err := configuration.Init()
	if err != nil {
		log.Println(err.Error())
		return
	}

	dbstring := configuration.GetDBString()
	err = db.Init(dbstring)
	if err != nil {
		panic("Error Connecting to the Database")
	}

	err = constants.Init()
	if err != nil {
		log.Println(err.Error())
		return
	}

	dbG := db.Get()
	defer dbG.Close()
	err = features.InitServer()
	if err != nil{
		log.Println("Error connecting to Server")
	}
}
