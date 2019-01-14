/*
	Title:	Test API methods
	Author:	Connor Peters
	Date:	2/26/2018
	Desc:	This file contains all the api calls for the test API
			It also contains the functions that supply the test API calls
*/

package dynauthapi

import (
	"bytes"
	"database/sql"

	// dynauthdb "github.com/mainsailstudio/dynauth-mvp/dynauth/dynauthdb"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mainsailstudio/dynauth-mvp/dynauth/dynauthconst"
	"golang.org/x/crypto/bcrypt"
)

// testAuth struct - for test users registering with dynauth standard
type testAuth struct {
	ID    string   `json:"id"`
	Locks []string `json:"locks"`
	Auths []string `json:"auths"`
}

// testKeys struct - for the user's plaintext keys
type testKeys struct {
	ID    string   `json:"id"`
	Keys  []string `json:"keys"`
	Locks []string `json:"locks"`
}

// testRegister - API call that checks if the test user is registered and then spits out the necessary information
func testRegister(w http.ResponseWriter, r *http.Request) {
	// create a testUser and decode the call
	var user testUser
	_ = json.NewDecoder(r.Body).Decode(&user)
	userExists, userID, userFname, userLname, userInit, userTestLevel, err := getTestUserInit(user.Email)
	if err != nil {
		fmt.Println("Error encountered when seeing if the user exists. Error is", err)
		return
	}

	// if the user is genuine
	if userExists {
		// issue them a new token using their email as part of the claims
		userToken := issueJWT(user.Email)

		// create a new testUser to return as JSON
		returnUser := testUser{ID: userID, Fname: userFname, Lname: userLname, Email: user.Email, Init: userInit, TestLevel: userTestLevel, Token: userToken}
		json.NewEncoder(w).Encode(returnUser)
	} else {
		http.Error(w, "This email is not pre-registered as a test user", 400)
	}
}

// testLoginLevel - API call to get the test users type of login and start the process
// NOTE - this stores the user's random locks to prepend for authentication
func testLoginLevel(w http.ResponseWriter, r *http.Request) {
	// create a testUser and decode the call
	var user testUser
	_ = json.NewDecoder(r.Body).Decode(&user)

	// is this necessary?
	if user.Email == "" {
		http.Error(w, "The email is empty", 400)
	}

	userExists, userID, _, _, userInit, userTestLevel, err := getTestUserInit(user.Email)
	if err != nil {
		fmt.Println("Error encountered when seeing if the user exists. Error is", err)
		return
	}

	if userInit == false {
		http.Error(w, "This account hasn't been initialized, please register first and then try logging in", 400)
		return
	}

	// if the user exists and if the user account has been initialized, serve up a random string of locks and store them temporarily
	if userExists == true && userInit == true {
		// test levels 2 and 3 are typical dynauth, the differentiator is the amount of auths
		if userTestLevel == 2 || userTestLevel == 3 {
			// this function also stores the served locks in the testTempLocks table
			lockSlice, err := serveTestLocks(userID, dynauthconst.DisplayLockNum)
			if err != nil {
				http.Error(w, "An error was encountered when getting the locks", 400)
			}

			// create a new testUser and return it
			newUserLogin := testUser{ID: user.ID, Email: user.Email, Init: userInit, Locks: lockSlice, TestLevel: userTestLevel}
			json.NewEncoder(w).Encode(newUserLogin)

			// The user isn't initialized so return a blank string for the locks
		} else {
			newUserLogin := testUser{ID: user.ID, Email: user.Email, Init: userInit, Locks: nil, TestLevel: userTestLevel}
			json.NewEncoder(w).Encode(newUserLogin)
		}
	} else {
		http.Error(w, "The email does not exist in our records", 400)
	}
}

// testLogin - API call finishes the user login process
func testLogin(w http.ResponseWriter, r *http.Request) {
	var user testUser
	_ = json.NewDecoder(r.Body).Decode(&user)

	// is this necessary?
	if user.Email == "" {
		http.Error(w, "The email is empty", 400)
	}

	userExists, userID, userFname, userLname, userInit, userTestLevel, err := getTestUserInit(user.Email)
	if err != nil {
		fmt.Println("Error encountered when seeing if the user exists. Error is", err)
		return
	}

	// if the user exists and if the user account has been initialized, serve up a random string of locks and store them temporarily
	if userExists == true {
		// test levels 2 and 3 are typical dynauth, the differentiator is the amount of auths
		if userTestLevel == 2 || userTestLevel == 3 {

			// this is where the authentication happens
			authCorrect, err := dynauthcore.AuthenticateSHA3WithLocksAPI(userID, user.Secret) // returns a boolean
			if err != nil {
				http.Error(w, "There was an error authenticating the user", 500)
			}

			// if authenticated
			if authCorrect {
				userToken := issueJWT(user.Email) // create and issue a new token using their email as a claim
				returnUser := testUser{ID: userID, Fname: userFname, Lname: userLname, Email: user.Email, Init: userInit, TestLevel: userTestLevel, Token: userToken}
				json.NewEncoder(w).Encode(returnUser)
			} else {
				fmt.Println("NOT authenticated via auths")
			}
		}

		if userTestLevel == 1 {
			// if the user is loggin in via a  pass
			passwordCorrect, err := dynauthcore.TestPassAuth(userID, user.Secret)
			if err != nil {
				http.Error(w, "An error was encountered when authentication a user with their auths", 500)
			}
			if passwordCorrect {
				fmt.Println("Correctly authenticated via a temp pass")
				userToken := issueJWT(user.Email)
				returnUser := testUser{ID: userID, Fname: userFname, Lname: userLname, Email: user.Email, Init: userInit, TestLevel: userTestLevel, Token: userToken}
				json.NewEncoder(w).Encode(returnUser)
			} else {
				http.Error(w, "Incorrect, not authenticated", 400)
			}
		}
	} else {
		http.Error(w, "The user doesn't exist", 400)
	}
}

// logConfigActivity - API call that logs the user's front-end activity
func logConfigActivity(w http.ResponseWriter, r *http.Request) {
	var log ConfigActivity
	_ = json.NewDecoder(r.Body).Decode(&log)

	err := insertConfigActivity(log)
	if err != nil {
		message := []string{"There was an issue logging the user's activity", err.Error()}
		errorString := strings.Join(message, " ")
		http.Error(w, errorString, 500)
	}
}

// logLoginActivity - API call that logs the user's front-end activity
func logLoginActivity(w http.ResponseWriter, r *http.Request) {
	var log LoginActivity
	_ = json.NewDecoder(r.Body).Decode(&log)

	err := insertLoginActivity(log)
	if err != nil {
		message := []string{"There was an issue logging the user's activity", err.Error()}
		errorString := strings.Join(message, " ")
		http.Error(w, errorString, 500)
	}
}

// logPracticeActivity - API call that logs the user's practice front-end activity
func logPracticeActivity(w http.ResponseWriter, r *http.Request) {
	var log LoginActivity
	_ = json.NewDecoder(r.Body).Decode(&log)

	err := insertPracticeActivity(log)
	if err != nil {
		message := []string{"There was an issue logging the user's activity", err.Error()}
		errorString := strings.Join(message, " ")
		http.Error(w, errorString, 500)
	}
}

// PROTECTED
// testRegisterPass - API call that posts the user's pass into the database
func testRegisterPass(w http.ResponseWriter, r *http.Request) {
	var user testPass
	_ = json.NewDecoder(r.Body).Decode(&user)

	// assume everything is set and store the hashed password
	err := storeTestPass(user.ID, user.HashedPassword)
	if err != nil {
		http.Error(w, "There was an error storing the test user's test password", 500)
	}

	// assume everything is set and store the plaintext password
	err = storeTestDisplayPass(user.ID, user.Password)
	if err != nil {
		http.Error(w, "There was an error storing the test user's display password", 500)
	}

	// flip the user's init flag
	err = initTestUser(user.ID)
	if err != nil {
		fmt.Println("Error encountered when initializing the user. Error is", err)
		return
	}

}

// PROTECTED
// testRegisterAuth - API call that posts the user's auths into the database
func testRegisterAuth(w http.ResponseWriter, r *http.Request) {
	var user testAuth
	_ = json.NewDecoder(r.Body).Decode(&user)

	// store the user's locks in the form of a received slice
	err := storeTestLocks(user.Locks, user.ID, "1") // using 1 as the locktype since it does nothing currently
	if err != nil {
		http.Error(w, "There was an error storing the test user's locks", 500)
	}

	// generate the permutations of the slice, assuming that FOR NOW 4 is the default
	lockPermArray := dynauthcore.LimPerms(user.Locks, 4) // 4 is the default

	// combine the permutated locks with the received array of hashed auths delivered by the frontend
	// NOTE: it is assumed the frontend is using the EXACT same permutation algorithms (find all subsets then permute all subsets) to generate them, otherwise it will not work
	combineArray := dynauthcore.CombinePerms(lockPermArray, user.Auths)

	// hash each perm again with a salt
	hashedPermsWithSalt := dynauthcore.HashPermsWithSaltSHA3(combineArray)

	// store the final slice of auths
	err = dynauthcore.StoreAuthsWithSalts(hashedPermsWithSalt, user.ID)
	if err != nil {
		http.Error(w, "Error encountered when storing the user's auth with salts", 500)
		return
	}

	// flip the user's init flag
	err = initTestUser(user.ID)
	if err != nil {
		http.Error(w, "Error encountered when initializing the user. Error is", 500)
	}
}

// PROTECTED
// testRegisterKeys - API call that posts the user's PLAINTEXT keys into the database
// NOTE: This is for testing only, the user's entire key should not be kept in plaintext EVER
func testRegisterKeys(w http.ResponseWriter, r *http.Request) {
	var user testKeys
	_ = json.NewDecoder(r.Body).Decode(&user)

	// store the user's locks in the form of a received slice
	err := storeTestKeys(user.Keys, user.Locks, user.ID) // using 1 as the locktype since it does nothing currently
	if err != nil {
		http.Error(w, "There was an issue storing the test user's plaintext keys", 500)
	}
}

// PROTECTED
// testGetUserKeys - API call that get's the user's keys in plaintext for testing
func testGetUserKeys(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	keyArray, err := getTestKeys(params["userID"][0])
	if err != nil {
		http.Error(w, "There was an issue getting your keys", 500)
	}

	json.NewEncoder(w).Encode(keyArray)
}

// PROTECTED
// testGetUserDisplayPass - API call that get's the user's password in plaintext for testing
func testGetUserDisplayPass(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	password, err := getTestDisplayPass(params["userID"][0])
	if err != nil {
		http.Error(w, "There was an issue getting your password to display", 500)
	}

	json.NewEncoder(w).Encode(password)
}

// getTestUserInit - get the test users information needed for the front-end
// NOTE: the name might be a bit misleading, but it's original purpose was just to fetch the user's init flag but not it just grabs all the user's information for efficiency
func getTestUserInit(email string) (bool, string, string, string, bool, int, error) {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return false, "", "", "", false, 0, errors.New("Opening the database connection for getTestUserInit went wrong")
	}
	defer db.Close()

	// assume the user does NOT exist
	exists := false

	// query the database to make sure this email doesn't already exist
	var userID string
	var fname string
	var lname string
	var init bool
	var testLevel int
	row := db.QueryRow("SELECT id, fname, lname, init, testLevel FROM testUsers where email = ?", email).Scan(&userID, &fname, &lname, &init, &testLevel)
	switch row {
	case sql.ErrNoRows: // if no rows were returned, the user does not exist
		exists = false
	}

	if userID != "" || testLevel != 0 {
		exists = true
	}
	return exists, userID, fname, lname, init, testLevel, nil
}

// initTestUser - flip the flag to initialize a user after the are registered
func initTestUser(userid string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("opening the database connection for initTestUser went wrong")
	}
	defer db.Close()

	// This is where each unique user password is created
	updateInit := "UPDATE testUsers SET init = true WHERE id = ?"
	smtUpdateInit, err := db.Prepare(updateInit)
	if err != nil {
		return errors.New("preparing to update the user's init flag went wrong")
	}
	defer smtUpdateInit.Close()

	_, err = smtUpdateInit.Exec(userid)
	if err != nil {
		return errors.New("executing the query to update the user's init flag went wrong")
	}
	fmt.Println("The user was initialized")
	return nil
}

// storeTestLocks - stores the test user' locks
// Needs the user's locks in a slice of strings, the user's id as a string, and the lockType as a string
func storeTestLocks(locks []string, userid string, lockType string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("opening the database connection for storeTestLocks went wrong")
	}
	defer db.Close()

	// Prepare statement for inserting the user's auth into the new table
	prepareStatement := "INSERT INTO testLocks VALUES("
	// for loop adds all perms into prepared statement
	for i := 1; i < len(locks); i++ {
		prepareStatement += "DEFAULT, ?, ?, ?), ("
	}
	prepareStatement += "DEFAULT, ?, ?, ?)"

	stmtIns, err := db.Prepare(prepareStatement)
	if err != nil {
		return errors.New("error encountered when preparing to insert the test user's locks")
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
		return errors.New("error encountered when executing the query to insert the test user's locks")
	}

	return nil // no errors, yay!
}

// serveTestLocks - to query the database and return the test user's locks in a string
func serveTestLocks(userid string, lockNum int) ([]string, error) {

	// this is a copy of the bode from GetLocks in the serve dynauthcorem package, it is here for AGILE development boy
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return nil, errors.New("error encountered when opening the database connection for serveTestLocks")
	}
	defer db.Close()

	lockSlice := []string{}
	locks, err := db.Query("SELECT locksAre FROM testLocks WHERE userid = ?", userid)
	if err != nil {
		return nil, errors.New("no locks were receieved from the database, user must not have initialized them")
	}
	defer locks.Close()
	for locks.Next() {
		var lockInfo string
		err := locks.Scan(&lockInfo)
		if err != nil {
			return nil, errors.New("locks weren't added to the slice properly for unknown reasons")
		}
		lockSlice = append(lockSlice, lockInfo)
	}

	if len(lockSlice) > 0 {
		fmt.Println("Locks look good")
		dynauthcore.Shuffle(lockSlice) // from internet code
		lockSlice = lockSlice[:lockNum]

		// to make into string
		var lockString string

		for i := range lockSlice {
			lockString += lockSlice[i]
		}
		err := storeTestTempLocks(userid, lockString)
		if err != nil {
			return nil, errors.New("there was an issue storing the user's temporary locks")
		}
		// return the slice NOT string
		return lockSlice, nil
	}

	return nil, errors.New("No locks received")

}

// storeTestPass - stores the test users password
func storeTestPass(userid string, password string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("error encountered when opening the database connection in storeTestPass")
	}
	defer db.Close()

	// This is where each unique user password is created
	initUserPass := "INSERT INTO testPass (userid, pass) VALUES (?, ?)"
	stmtInsPass, err := db.Prepare(initUserPass)
	if err != nil {
		return errors.New("error encountered when preparing to insert the test user's testPass")
	}
	defer stmtInsPass.Close()

	// create hashed password
	hashedPasswordBcrypt, err := bcrypt.GenerateFromPassword([]byte(password), dynauthconst.BcryptIterations)
	if err != nil {
		return errors.New("error encountered when hashing the user's password using Bcrypt")
	}
	password = bytes.NewBuffer(hashedPasswordBcrypt).String()

	_, err = stmtInsPass.Exec(userid, password)
	if err != nil {
		return errors.New("error encountered when executing the query to insert the test user's testPass")
	}

	return nil
}

// storeTestTempLocks - stores the test user' locks TEMPORARILY for session based authentication
// important note, if the user already has an entry for locks (which is likely), this will UPDATE that entry
func storeTestTempLocks(userid string, locks string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("error encountered when opening the database connection for storeTestTempLocks")
	}
	defer db.Close()

	// first check if the user already has temp locks
	// NOTE: we don't care about the output at all, only if it exists
	var userLocksExists bool
	var whoCares string // I sure don't
	row := db.QueryRow("SELECT userid FROM tempTestLocks where userid = ?", userid).Scan(&whoCares)
	switch row {
	case sql.ErrNoRows:
		userLocksExists = false
	default:
		userLocksExists = true
	}

	// if the user already has up to date locks
	// update the existing lock entry rather than inserting
	if userLocksExists {
		// This is where each unique user password is created
		updateInit := "UPDATE tempTestLocks SET locks = ? WHERE userid = ?"
		smtUpdateInit, err := db.Prepare(updateInit)
		if err != nil {
			return errors.New("error encountered when preparing to update the test user's temp locks used for the final stage of dynamic authentication")
		}
		defer smtUpdateInit.Close()

		_, err = smtUpdateInit.Exec(locks, userid)
		if err != nil {
			return errors.New("error encountered when executing the query to update the test user's temp locks used for the final stage of dynamic authentication")
		}

		// if the user is logging in for the first time (either at all or in awhile, if that makes any sense)
		// insert a new lock entry rather than updating an existing
	} else {
		// Prepare statement for inserting the user's auth into the new table
		prepareStatement := "INSERT INTO tempTestLocks VALUES(?, ?, ?)"
		stmtInsLocks, err := db.Prepare(prepareStatement)
		if err != nil {
			return errors.New("error encountered when preparing to insert the test user's temp locks used for the final stage of dynamic authentication.\nServed error was")
		}
		defer stmtInsLocks.Close()

		// have the locks be relevent for 1 day
		// NOTE: currently, there is no check here but it is likely this will be implemented in the future
		expireDate := time.Now().Local().AddDate(0, 0, 1)
		expireDate.Format("2006-01-02 15:04:05")

		_, err = stmtInsLocks.Exec(userid, locks, expireDate)
		if err != nil {
			return errors.New("error encountered when executing the query to insert the test user's temp locks used for the final stage of dynamic authentication")
		}
	}

	return nil
}

// storeTestKeys - stores the test user's plaintext keys
func storeTestKeys(keys []string, locks []string, userid string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("opening the database connection for storeTestKeys went wrong")
	}
	defer db.Close()

	// Prepare statement for inserting the user's auth into the new table
	prepareStatement := "INSERT INTO testKeysDisplay VALUES("
	// for loop adds all perms into prepared statement
	for i := 1; i < len(keys); i++ {
		prepareStatement += "DEFAULT, ?, ?, ?, ?), ("
	}
	prepareStatement += "DEFAULT, ?, ?, ?, ?)"

	stmtIns, err := db.Prepare(prepareStatement)
	if err != nil {
		return errors.New("error encountered when preparing to insert the test user's keys")
	}
	defer stmtIns.Close()

	// casts the data to insert into a slice interface for variadic function inclusion below, quite elegant
	dataPrepared := []interface{}{}
	for i := 0; i < len(keys); i++ {
		dataPrepared = append(dataPrepared, userid)
		dataPrepared = append(dataPrepared, keys[i])
		dataPrepared = append(dataPrepared, locks[i])
		dataPrepared = append(dataPrepared, 1) // the 1 is the default key type, not relevent for now
	}
	_, err = stmtIns.Exec(dataPrepared...) // adds all data in the slice as a separate argument (variadic) BEAUTIFUL
	if err != nil {
		fmt.Println("Err was", err)
		return errors.New("error encountered when executing the query to insert the test user's keys")
	}

	return nil // no errors, yay!
}

// storeTestDisplayPass - stores the test user's plaintext password
func storeTestDisplayPass(userid string, password string) error {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return errors.New("opening the database connection for storeTestPass went wrong")
	}
	defer db.Close()

	// Prepare statement for inserting the user's auth into the new table
	prepareStatement := "INSERT INTO testPassDisplay VALUES(DEFAULT, ?, ?)"

	stmtIns, err := db.Prepare(prepareStatement)
	if err != nil {
		return errors.New("error encountered when preparing to insert the test user's display password")
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(userid, password)
	if err != nil {
		fmt.Println("Err was", err)
		return errors.New("error encountered when executing the query to insert the test user's display password")
	}

	return nil // no errors, yay!
}

// getTestKeys queries the database and returns all of the user's locks into a slice
func getTestKeys(userid string) ([][]string, error) {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return nil, errors.New("Unable to connect to the database in the getTestKeys function in serve.go")
	}
	defer db.Close()

	comboSlice := [][]string{}
	keys, err := db.Query("SELECT keysAre, keyLocks FROM testKeysDisplay WHERE userid = ?", userid)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("No locks were receieved from the database, user must not have initialized them")
	}
	defer keys.Close()
	for keys.Next() {
		keyLockSlice := []string{}
		var keyInfo string
		var lockInfo string
		err := keys.Scan(&keyInfo, &lockInfo)
		if err != nil {
			return nil, errors.New("Locks weren't added to the slice properly for unknown reasons")
		}
		keyLockSlice = append(keyLockSlice, lockInfo)
		keyLockSlice = append(keyLockSlice, keyInfo)
		comboSlice = append(comboSlice, keyLockSlice)
	}

	return comboSlice, nil
}

// getTestDisplayPass queries the database and returns the user's display password as a string
func getTestDisplayPass(userid string) (string, error) {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return "", errors.New("Unable to connect to the database in the getTestDisplayPass function in serve.go")
	}
	defer db.Close()

	var password string
	row := db.QueryRow("SELECT pass FROM testPassDisplay WHERE userid = ?", userid).Scan(&password)
	switch row {
	case sql.ErrNoRows: // if no rows were returned, the user does not exist
		return "", errors.New("No password was returned")
	}

	return password, nil
}
