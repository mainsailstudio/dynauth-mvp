/**
Title:	Password  API
Desc:	Password manager secured with dynamic authentication and honey encryption
Notes:	THIS IS NOT CONSIDERED A SECURE PASSWORD MANAGER YET. This product is in it's pre-alpha stage and isn't really designed to be used by anyone yet, it's mostly for testing purposes
*/

// Package pwmapi creates all the routes and the corresponding functions for the password manager to be used "RESTfully"
package pwmapi

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mainsailstudio/dynauth-mvp/pwm/pwmdb"
	"github.com/rs/cors"
	respond "gopkg.in/matryer/respond.v1"
)

/**
=============================================
	API structures
=============================================
*/

// Password structure - contains the password and the associated URL and email/username
type Password struct {
	URL      string `json:"URL"`
	Email    string `json:"email"` // this is the field used for username as well since the function is identical
	Password string `json:"password"`
}

/**
=============================================
	Start API
=============================================
*/

// Start - start the password manager API and create all the HTTP routes using Mux
func Start() {
	pwmdb.InitDB()

	fmt.Println("Starting New Password Manager API Service")
	srv := initRouter()

	if _, err := os.Stat("./gitignore/cert.pem"); os.IsNotExist(err) {
		initTLSKeyAndCertificate()
	}
	if _, err := os.Stat("./gitignore/key.pem"); os.IsNotExist(err) {
		initTLSKeyAndCertificate()
	}

	log.Fatal(srv.ListenAndServeTLS("./gitignore/cert.pem", "./gitignore/key.pem"))
}

/**
=============================================
	Helper functions
=============================================
*/

func initRouter() *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Write([]byte("This is an example server.\n"))
	})

	// GET password entry based on the URL
	// Example call: /password?url=https://twitter.com
	r.HandleFunc("/password", getPassword).Queries("url", "{url}").Methods("GET")

	// POST new password entry with URL, email, and password
	// Example call: /password?url=https://twitter.com&email=test@test.com&password=password
	r.HandleFunc("/password", createPassword).Queries("url", "{url}", "email", "{email}", "password", "{password}").Methods("POST")

	// PUT password to update existing entry. Queries based on URL; either email or password are optional
	// Example call: : /password?url=https://twitter.com&password=newpassword
	r.HandleFunc("/password", updatePassword).Queries("url", "{url}", "email", "{email}", "password", "{password}").Methods("PUT")

	// DELETE password to delete existing entry. Queries based on URL
	// Example call: : /password?url=https://twitter.com
	r.HandleFunc("/password", deletePassword).Methods("DELETE")

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	// Add the default handler to allow CORS
	handler := cors.Default().Handler(r)

	// Configure the TLS server
	srv := &http.Server{
		Addr:         ":443",
		Handler:      handler,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	return srv
}

// connectPasswordDatabase - Helper function that connects to the correct database
func connectPasswordDatabase() {

}

// authenticateDatabase - Helper function that uses the dynamic authentication service to authenticate the user and gain access to the database
func authenticateDatabase() {

}

// encryptPassword - Helper function that encrypts a password
func encryptPassword() {

}

// encryptPassword - Helper function that honeys(?) a password before encryption
func honeyPassword() {

}

// decryptPassword - Helper function that decrypts a password
func decryptPassword() {

}

/**
=============================================
	API functions
=============================================
*/

// createPassword - POST request that adds a password to the password database along with an email/username and URL associated with the password
// PROTECTED function, accessibly only when authenticated
func createPassword(w http.ResponseWriter, r *http.Request) {

}

// getPassword - GET request that retrieves and decrypts a password from the password database based on the email/username and URL
// PROTECTED function, accessibly only when authenticated
func getPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["url"]
	if len(url) == 0 { // URL is empty
		fmt.Println("Url Param 'email' is missing")
		http.Error(w, "Parameter URL was empty in GET request", http.StatusBadRequest)
		return
	}

	fmt.Println("Url Param 'url' is: " + url)

	// database, err := sql.Open("sqlite3", "./private/database.db")
	// if err != nil {
	// 	fmt.Println("\n There was an issue readying the database \n")
	// 	log.Fatal(err)
	// }

	// row, err := database.Query("SELECT password FROM passwords WHERE email = ?", email)

	// var password string
	// if err := database.QueryRow("SELECT password FROM passwords WHERE email = ?", email).Scan(&password); err != nil {
	// 	fmt.Println("error oops")
	// 	log.Fatal(err)
	// }

	// fmt.Println("Password for url " + url + " is " + password)

	passwordData := Password{URL: url, Email: "example@example.com", Password: "testPassword"}

	// respond with OK, and the data
	respond.With(w, r, http.StatusOK, passwordData)
}

// updatePassword - PUT request that updates a password type, including ability to update email/username and URL associated with the password
// PROTECTED function, accessibly only when authenticated
func updatePassword(w http.ResponseWriter, r *http.Request) {

}

// deletePassword - DELETE request that deletes an entire password type, along with the username and emails associated with the password
// PROTECTED function, accessibly only when authenticated
func deletePassword(w http.ResponseWriter, r *http.Request) {

}
