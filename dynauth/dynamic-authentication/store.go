/*
	Title:	Storing package
	Author:	Connor Peters
	Date:	2/12/2018
	Desc:
*/

package dynauthcore

import (
	"database/sql"
	dbinfo "dbinfo"
	"errors"

	_ "github.com/go-sql-driver/mysql" // mysql driver helper
)

// StoreAuthsWithSalts - to store a slice of hashed permutations into a MySQL database.
func StoreAuthsWithSalts(authsWithSalts [][]string, userid string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("Unable to connect to the database in the StoreAuthsWithSalts function")
	}
	defer db.Close()

	// This is where each unique user auth table is created
	createTable := "CREATE TABLE auth" + userid + " (auth char(128) binary, salt char(64) binary)"
	_, err = db.Exec(createTable) // like lean cuisine no preperation needed
	if err != nil {
		return errors.New("Issue creating the user's auths and salts table")
	}

	// Prepare statement for inserting the user's auth into the new table
	prepareStatement := "INSERT INTO auth" + userid + " VALUES("
	// for loop adds all perms into prepared statement
	for i := 1; i < len(authsWithSalts); i++ {
		prepareStatement += "?, ?), ("
	}
	prepareStatement += "?, ?)"

	stmtIns, err := db.Prepare(prepareStatement)
	if err != nil {
		return errors.New("Issue preparing the statement to insert the user's auths and salts into the user's auth table")
	}
	defer stmtIns.Close()

	// casts the data to insert into a slice interface for variadic function inclusion below, quite elegant
	dataPrepared := []interface{}{}
	for i := 0; i < len(authsWithSalts); i++ {
		dataPrepared = append(dataPrepared, authsWithSalts[i][0], authsWithSalts[i][1])
	}
	_, err = stmtIns.Exec(dataPrepared...) // adds all data in the slice as a separate argument (variadic) BEAUTIFUL
	if err != nil {
		return errors.New("Issue executing the query to insert the user's auths and salts into the user's auth table")
	}

	return nil
}

// StoreAuthsPlain - to store a slice of hashed permutations into a MySQL database.
func StoreAuthsPlain(auths []string, userid string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("Unable to connect to the database in the StoreAuthsPlain function")
	}
	defer db.Close()

	// This is where each unique user auth table is created
	createTable := "CREATE TABLE auth" + userid + " (auth char(64) binary)"
	_, err = db.Exec(createTable) // like lean cuisine no preperation needed
	if err != nil {
		return errors.New("Issue creating the user's plain auth table")
	}

	// Prepare statement for inserting the user's auth into the new table
	prepareStatement := "INSERT INTO auth" + userid + " VALUES("
	// for loop adds all perms into prepared statement
	for i := 1; i < len(auths); i++ {
		prepareStatement += "?), ("
	}
	prepareStatement += "?)"

	stmtIns, err := db.Prepare(prepareStatement)
	if err != nil {
		return errors.New("Issue preparing the statement to insert the user's plain auths into the user's auth table")
	}
	defer stmtIns.Close()

	// casts the data to insert into a slice interface for variadic function inclusion below, quite elegant
	dataPrepared := []interface{}{}
	for i := 0; i < len(auths); i++ {
		dataPrepared = append(dataPrepared, auths[i])
	}
	_, err = stmtIns.Exec(dataPrepared...) // adds all data in the slice as a separate argument (variadic) BEAUTIFUL
	if err != nil {
		return errors.New("Issue executing the query to insert the user's plain auths into the user's auth table")
	}

	return nil
}

// StoreLocks function stores the user' locks
// Needs the user's locks in a slice of strings, the user's id as a string, and the lockType as a string
func StoreLocks(locks []string, userid string, lockType string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("Unable to connect to the database in the StoreLocks function")
	}
	defer db.Close()

	// Prepare statement for inserting the user's auth into the new table
	prepareStatement := "INSERT INTO locks VALUES("
	// for loop adds all perms into prepared statement
	for i := 1; i < len(locks); i++ {
		prepareStatement += "DEFAULT, ?, ?, ?), ("
	}
	prepareStatement += "DEFAULT, ?, ?, ?)"

	stmtIns, err := db.Prepare(prepareStatement)
	if err != nil {
		return errors.New("Issue preparing the statement to insert the user's locks")
	}
	defer stmtIns.Close()

	// casts the data to insert into a slice interface for variadic function inclusion below, quite elegant
	dataPrepared := []interface{}{}
	for i := 0; i < len(locks); i++ {
		dataPrepared = append(dataPrepared, userid)
		dataPrepared = append(dataPrepared, locks[i])
		dataPrepared = append(dataPrepared, lockType)
	}
	_, err = stmtIns.Exec(dataPrepared...) // adds all data in the slice as a separate argument (variadic) BEAUTIFUL
	if err != nil {
		return errors.New("Issue executing the query to insert the user's locks")
	}

	return nil
}

// StoreUserInfo takes in userinfo and stores it
func StoreUserInfo(userid string, fname string, lname string, email string, phone string, securityLv string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("Unable to connect to the database in the StoreUserInfo function")
	}
	defer db.Close()

	// This is where each unique user is created
	prepareStatement := "INSERT INTO users (id, fname, lname, email, phone, security) VALUES (?, ?, ?, ?, ?, ?)"
	// _, err = db.Exec(createTable) // like lean cuisine no preperation needed
	stmtIns, err := db.Prepare(prepareStatement)
	if err != nil {
		return errors.New("Issue preparing the statement to insert the user's information")
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(userid, fname, lname, email, phone, securityLv)
	if err != nil {
		return errors.New("Issue executing the query to insert the user's information")
	}

	return nil
}

// StoreTempLocks - stores the user' locks temporarily for *GASP* session based authentication
// Needs the user's locks as a string
func StoreTempLocks(locks []string, userid string, lockType string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("Unable to connect to the database in the StoreTempLocks function")
	}
	defer db.Close()

	// Prepare statement for inserting the user's auth into the new table
	prepareStatement := "INSERT INTO locks VALUES("
	// for loop adds all perms into prepared statement
	for i := 1; i < len(locks); i++ {
		prepareStatement += "DEFAULT, ?, ?, ?), ("
	}
	prepareStatement += "DEFAULT, ?, ?, ?)"

	stmtIns, err := db.Prepare(prepareStatement)
	if err != nil {
		return errors.New("Issue preparing the statement to insert the user's locks used for *GASP* session based authentication")
	}
	defer stmtIns.Close()

	// casts the data to insert into a slice interface for variadic function inclusion below, quite elegant
	dataPrepared := []interface{}{}
	for i := 0; i < len(locks); i++ {
		dataPrepared = append(dataPrepared, userid)
		dataPrepared = append(dataPrepared, locks[i])
		dataPrepared = append(dataPrepared, lockType)
	}
	_, err = stmtIns.Exec(dataPrepared...) // adds all data in the slice as a separate argument (variadic) BEAUTIFUL
	if err != nil {
		return errors.New("Issue executing the query to insert the user's locks used for *GASP* session based authentication")
	}

	return nil
}
