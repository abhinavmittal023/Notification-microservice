package main

import (
	"flag"
	"log"
	"os"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	_ "code.jtg.tools/ayush.singhal/notifications-microservice/db/migrations"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	_ = flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 2 {
		flags.Usage()
		return
	}

	dbstring, command := args[1], args[2]

	err := db.Init(dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}
	defer db.Get().Close()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db.GetDbHandle(), *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
