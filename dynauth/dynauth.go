/*
	Title:	Dynauth Service Start
	Author:	Connor Peters
	Date:	12/28/2018
	Desc:	Just starts the local service for now! Eventually there will be more features in here like a CLI configuration UI or something
*/

package dynauth

import (
	"database/sql"
	"fmt"
	"log"

	dynauthapi "github.com/mainsailstudio/dynauth-mvp/dynauth/dynauthapi"
	dynauthdb "github.com/mainsailstudio/dynauth-mvp/dynauth/dynauthdb"

	// sqlite 3 import must be a blank import
	_ "github.com/mattn/go-sqlite3"
)

// Start the dynauth local service
// Seems redundant considering the the mvp package calls essentially the same thing, but it's necessary for future extensibility. More functions can be easily added into the package without disrupting the application initialization
func Start() {
	// initialize the global DB connection
	initDynauthDb()

	fmt.Println("Starting dynauth local API")
	dynauthapi.Start()
}

func initDynauthDb() {
	var dbPath = "dynauth/gitignore/dynauth.db"
	var err error
	dynauthdb.DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Test if the database is already configured
	_, err = dynauthdb.DB.Query("SELECT url FROM passwords LIMIT 1;")
	if err != nil {
		fmt.Println("Database is new, create the tables")
		dynauthdb.CreateTables()
	} else {
		fmt.Println("The database is good to go!")
	}
	// Keep the connection open!!!
	// defer dynauthdb.DB.Close()
}
