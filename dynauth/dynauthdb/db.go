package dynauthdb

import (
	"database/sql"
	"log"

	// sqlite 3 import must be a blank import
	_ "github.com/mattn/go-sqlite3"
)

// DB global variable to connect to database across api
var DB *sql.DB

// CreateTables initializes the database
func CreateTables() {
	_, err := DB.Exec(`CREATE TABLE passwords (
		id	INTEGER PRIMARY KEY AUTOINCREMENT,
		url	TEXT,
		email	TEXT,
		password	TEXT
	);`)
	if err != nil {
		log.Panic(err)
	}
}
