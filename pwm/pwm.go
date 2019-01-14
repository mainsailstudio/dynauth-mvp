/*
	Title:	Password Manager Service Manager
	Author:	Connor Peters
	Date:	12/28/2018
	Desc:	Just starts the local service for now! Eventually there will be more features in here like a CLI configuration UI or something
*/

package pwm

import (
	"database/sql"
	"fmt"
	"log"

	pwmapi "github.com/mainsailstudio/dynauth-mvp/pwm/pwmapi"
	pwmdb "github.com/mainsailstudio/dynauth-mvp/pwm/pwmdb"

	// sqlite 3 import must be a blank import
	_ "github.com/mattn/go-sqlite3"
)

// Start the pwm local service
// Seems redundant considering the the mvp package calls essentially the same thing, but it's necessary for future extensibility. More functions can be easily added into the package without disrupting the application initialization
func Start() {
	// initialize the global DB connection
	initPwmDb()

	fmt.Println("Starting pwm local API")
	pwmapi.Start()
}

func initPwmDb() {
	var dbPath = "pwm/gitignore/pwm.db"
	var err error
	pwmdb.DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Test if the database is already configured
	_, err = pwmdb.DB.Query("SELECT url FROM passwords LIMIT 1;")
	if err != nil {
		fmt.Println("Database is new, create the tables")
		pwmdb.CreateTables()
	} else {
		fmt.Println("The database is good to go!")
	}
	// Keep the connection open!!!
	// defer pwmdb.DB.Close()
}
