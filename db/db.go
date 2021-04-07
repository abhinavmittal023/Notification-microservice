package db

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // import for psql handler
	"github.com/pkg/errors"
)

// Db gorm type database handler
var dbG *gorm.DB

// DbHandle sqltype database handler
var dbHandle *sql.DB

// Get returns the gorm handle
func Get() *gorm.DB {
	return dbG
}

// GetDbHandle returns the sql handle
func GetDbHandle() *sql.DB {
	return dbHandle
}

// Init is the initialisation function
func Init(url string) error {
	var err error
	dbG, err = gorm.Open("postgres", url)
	if err != nil {
		return errors.Wrap(err, "Unable to initialize the database connection")
	}
	// get db connection handle
	dbHandle = dbG.DB()
	dbHandle.SetMaxOpenConns(99)
	dbHandle.SetMaxIdleConns(10)
	return nil
}
