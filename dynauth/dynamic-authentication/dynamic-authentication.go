/*
	Title:	Authenticate User
	Author:	Connor Peters
	Date:	ORIGINALLY 3/17/2018, modified 12/23/2018
	Desc:	Simply authenticates a user based on their email and hashes of the keys/locks
*/

package dynamicAuthentication

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// start() - start the password manager API and create all the HTTP routes using Mux
func Start(protocol string) {
	// Init router
	fmt.Println("Starting New Dynamic Authentication API")
	r := mux.NewRouter()

	// Init the routes
	/*
	- getLocks(): GET that retrieves a user's locks and temporarily stores requested locks on the server side
		- GET: /auth/users/locks?userid=123456789
	- authenticate(): GET that compares the submitted hash to storage to determine if the user is authenticated. Returns a token
		- GET: /auth/users/keys?userid=123456789&secret=hashedsecrethere
	- register(): POST that registers a user using email and dynamic authentication
		- POST: /auth/users?email=example@example.com&secret=allhashedlocksandkeyshere
	- logAuthenticationAttempt()
		- POST /auth/users/logs?userid=123456789&logVar=example
	*/
	r.HandleFunc("/users/locks", getUserLocks).Methods("GET")
	r.HandleFunc("/users/keys", authenticateUser).Methods("GET")
	r.HandleFunc("/users", registerUser).Methods("POST")
	r.HandleFunc("/users/logs", logUserAuthentication).Methods("POST")

	// Add the default handler to allow CORS
	handler := cors.Default().Handler(r)

	if protocol == "http" {
		// Serve over normal HTTP, not HTTPS
		http.ListenAndServe(":8080/auth", handler)
	} else {
		// // restricted API call to register a user with a password
		// // requires a proper JWT token to access
		// r.Handle("/test/register-pass", negroni.New(
		// 	negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		// 	negroni.Wrap(http.HandlerFunc(testRegisterPass)),
		// ))

		// // restricted API call to register a user's auths
		// // requires a proper JWT token to access
		// r.Handle("/test/register-auth", negroni.New(
		// 	negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		// 	negroni.Wrap(http.HandlerFunc(testRegisterAuth)),
		// ))

		// Listen and serve the API over TLS (HTTPS)
		err := http.ListenAndServeTLS(":443", "../private/certificate.crt", "../private/private.key", handler)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}


/**
=============================================
	API structures
=============================================
*/

// user struct
type user struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Locks     []string `json:"locks"`
	Token     string   `json:"token"`
	Secret    string   `json:"secret"`
}

/**
=============================================
	API functions
=============================================
*/

// getUserLocks() - GET request that retrieves a user's locks and temporarily stores requested locks on the server side
func getUserLocks(w http.ResponseWriter, r *http.Request) {

}

// authenticateUser() - GET request that compares the submitted hash to storage to determine if the user is authenticated. Returns a JWT token
func authenticateUser(w http.ResponseWriter, r *http.Request) {
	var user user
	_ = json.NewDecoder(r.Body).Decode(&user)

	// is this necessary?
	if user.Email == "" {
		http.Error(w, "The email is empty", 400)
	}

	userID, err := getTestUserInit(user.Email)
	if err != nil {
		fmt.Println("Error encountered when seeing if the user exists. Error is", err)
		return
	}

	// if the user exists and if the user account has been initialized, serve up a random string of locks and store them temporarily
	if len(userID) > 0 { // the user exists
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
	} else {
		http.Error(w, "The user doesn't exist", 400)
	}
}

// registerUser() - POST request that registers a user using email and dynamic authentication
func registerUser(w http.ResponseWriter, r *http.Request) {

}

// logUserAuthentication() - POST request that logs a user's authentication activity for later analysis
// PROTECTED function, accessibly only when authenticated
func logUserAuthentication(w http.ResponseWriter, r *http.Request) {

}

/**
=============================================
	Helper functions
=============================================
*/

// AuthenticateSHA3WithLocksAPI - to perform authentication using the locks that were stored when the user was served, this is NOT a RESTful API anymore, but who cares?
func AuthenticateSHA3WithLocksAPI(userid string, auth string) (bool, error) {
	tempLocks, err := getTestTempLocks(userid)
	if err != nil {
		fmt.Println("Error enountered:", err)
	}
	auth = tempLocks + auth
	authenticated, err := compareAuthsWithSaltSHA3(auth, userid)
	if err != nil {
		return false, errors.New("There was an issue comparing the user's auths. Chances are this is because they don't have an auth table setup")
	}

	if authenticated == true {
		fmt.Println("AUTHENTICATED")
		return true, nil
	}
	fmt.Println("NO MATCH FOUND")
	return false, nil

}

// compareAuthsWithSaltsSHA3
func compareAuthsWithSaltSHA3(toCompare string, userid string) (bool, error) {
	// initialize a false authentication return
	authenticated := false

	salts, err := getSalts(userid)
	if err != nil {
		fmt.Println("Error encountered:", err)
	}

	authSlice, err := getAuths(userid) // get all of the auths into a slice
	if err != nil {
		return false, err
	}

	// ready the private key
	pkSecret, err := ioutil.ReadFile("../../../private.ppk") // in form of byte
	if err != nil {
		log.Fatal(err)
	}

	for i := range authSlice {

		// DEBUG fmt.Println("Compare number", i)

		h := make([]byte, 64)
		d := sha3.NewShake256()
		// Write the key into the hash.
		d.Write(pkSecret)
		// Now write the data.
		d.Write([]byte(toCompare + salts[i]))
		d.Read(h)

		hashString := fmt.Sprintf("%x", h)
		if strings.Compare(hashString, authSlice[i]) == 0 {
			authenticated = true
			break
		}
	}

	return authenticated, nil
}

// getAuths()
func getAuths(userid string) ([]string, error) {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return nil, errors.New("Opening the database connection for getAuths went wrong")
	}

	defer db.Close()
	authSlice := []string{}
	query := "SELECT auth FROM auth" + userid
	auths, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer auths.Close()
	for auths.Next() {
		var auth string
		err := auths.Scan(&auth)
		if err != nil {
			log.Fatal(err)
		}
		authSlice = append(authSlice, auth)
	}
	err = auths.Err()
	if err != nil {
		log.Fatal(err)
	}
	return authSlice, nil
}

// getSalts()
func getSalts(userid string) ([]string, error) {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return nil, errors.New("Opening the database connection for getSalts went wrong")
	}

	defer db.Close()
	lockSlice := []string{}
	query := "SELECT salt FROM auth" + userid
	locks, err := db.Query(query)
	if err != nil {
		return nil, errors.New("Getting the user's salts caused an error")
	}

	defer locks.Close()
	for locks.Next() {
		var lock string
		err := locks.Scan(&lock)
		if err != nil {
			log.Fatal(err)
		}
		lockSlice = append(lockSlice, lock)
	}
	err = locks.Err()
	if err != nil {
		log.Fatal(err)
	}
	return lockSlice, nil
}


// this gets the user's temp locks that were stored when served to enable the authentication
func getTestTempLocks(userid string) (string, error) {
	dbinfo := dbinfo.Db()
	db, err := sql.Open(dbinfo[0], dbinfo[1]) // gets the database information from the dbinfo package and enters the returned slice values as arguments
	if err != nil {
		return "", errors.New("Opening the database connection for getTestTempLocks went wrong")
	}
	defer db.Close()

	// select the temp pass of the user
	var lockString string
	err = db.QueryRow("SELECT locks FROM tempTestLocks WHERE userid = ?", userid).Scan(&lockString)
	if err != nil {
		log.Fatal(err)
	}
	return lockString, nil
}
