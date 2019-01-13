package dynauthcore

import (
	"database/sql"
	dbinfo "dbinfo"
	"errors"
	"fmt"
)

func checkUserExists(email string) (bool, string, error) {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return false, "", errors.New("Opening the database connection for checkUserExists went wrong")
	}
	defer db.Close()

	exists := false
	// search to make sure this email doesn't already exist
	var userID string
	row := db.QueryRow("SELECT id FROM users where email = ?", email).Scan(&userID)
	switch row {
	case sql.ErrNoRows:
		fmt.Println("No rows selected")
		exists = false
	default:
		exists = true
	}
	return exists, userID, nil
}
