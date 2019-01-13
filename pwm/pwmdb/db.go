// Package pwmdb manages all the configuration for the pwm database across the entire app
package pwmdb

import (
	"database/sql"
	"fmt"
	"log"

	// sqlite 3 import must be a blank import
	_ "github.com/mattn/go-sqlite3"
)

// db global variable to connect to database across api
var db *sql.DB

// InitDB - initialize the db global variable to make DB connection a breeze
func InitDB() {
	var dbPath = "./gitignore/pwm.db"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// fail-fast if can't connect to DB
	err = db.Ping()
	if err != nil {
		fmt.Println("Ping to database failed")
		log.Fatal(err)
	}
}

// CreateTables - create the correct tables for the password manager
func CreateTables() {
	_, err := db.Exec("CREATE TABLE example ( id integer, data varchar(32) )")
	if err != nil {
		log.Fatal(err)
	}
}
