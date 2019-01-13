/**
Title:	Password Manager
Author:	Connor Peters
Date:	12/29/2018
Desc:	DB connection handler for pwmapi
*/

// Package pwmapi creates all the routes and the corresponding functions for the password manager to be used "RESTfully"
package pwmapi

import (
	"database/sql"
	"fmt"
	"log"
	// _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB - initialize the db global variable to make DB connection a breeze
func InitDB() {
	var err error
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		fmt.Println("There was an issue readying the database")
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}
