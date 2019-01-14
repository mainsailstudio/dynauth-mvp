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
	pwmdb "github.com/mainsailstudio/dynauth-mvp/pwm/pwmdb"
	"github.com/rs/cors"
	respond "gopkg.in/matryer/respond.v1"
)

/**
=============================================
	API structures
=============================================
*/

// EncryptedPassword structure - contains the encrypted password and the associated locks to decrypt
type EncryptedPassword struct {
	URL      string `json:"URL"`
	Email    string `json:"email"` // this is the field used for username as well since the function is identical
	Password string `json:"password"`
	Locks    string `json:"locks"`
}

// Password structure - contains the decrypted (plaintext) password and the associated URL and email/username
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
	srv := initRouter()

	if _, err := os.Stat("pwm/gitignore/cert.pem"); os.IsNotExist(err) {
		initTLSKeyAndCertificate()
	}
	if _, err := os.Stat("pwm/gitignore/key.pem"); os.IsNotExist(err) {
		initTLSKeyAndCertificate()
	}

	log.Fatal(srv.ListenAndServeTLS("pwm/gitignore/cert.pem", "pwm/gitignore/key.pem"))
}

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

	// GET locks based on based on the URL
	// Example call: /lock?url=https://twitter.com
	r.HandleFunc("/locks", getLocks).Queries("url", "{url}").Methods("GET")

	// POST new temp lock entry for encryption password based on url and email
	// Example call: /locks?url=https://twitter.com&email=test@test.com
	// r.HandleFunc("/locks", storeTempLocks).Queries("url", "{url}", "email", "{email}").Methods("POST")

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

/**
=============================================
	API functions
=============================================
*/

// createPassword - POST request that adds a password to the password database along with an email/username and URL associated with the password
// PROTECTED function, accessibly only when authenticated
func createPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["url"]
	if len(url) == 0 { // URL is empty
		fmt.Println("Url Param 'url' is missing")
		respond.With(w, r, http.StatusBadGateway, "Url not understood")
		return
	}

	email := vars["email"]
	if len(url) == 0 { // Email is empty
		fmt.Println("Url Param 'email' is missing")
		respond.With(w, r, http.StatusBadGateway, "Email not understood")
		return
	}

	password := vars["password"]
	if len(url) == 0 { // Password is empty
		fmt.Println("Url Param 'password' is missing")
		respond.With(w, r, http.StatusBadGateway, "Password not understood")
		return
	}

	err := encryptAndStorePassword(url, email, password)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

}

// getPassword - GET request that retrieves and decrypts a password from the password database based on the email/username and URL
// PROTECTED function, accessibly only when authenticated
func getPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["url"]
	if len(url) == 0 { // URL is empty
		fmt.Println("Url Param 'url' is missing")
		respond.With(w, r, http.StatusBadGateway, "Url not understood")
		return
	}

	encryptedPass, err := getEncryptedPassword(url)
	if err != nil {
		fmt.Println("Password for url " + url + " doesn't exist")
		respond.With(w, r, http.StatusNotFound, err)
		return
	}

	decryptedPass, err := decryptPassword(encryptedPass)
	if err != nil {
		fmt.Println("Can't decrypt url:  " + url)
		respond.With(w, r, http.StatusNotFound, err)
		return
	}

	// respond with OK, and the data
	respond.With(w, r, http.StatusOK, decryptedPass)
}

// updatePassword - PUT request that updates a password type, including ability to update email/username and URL associated with the password
// PROTECTED function, accessibly only when authenticated
func updatePassword(w http.ResponseWriter, r *http.Request) {

}

// deletePassword - DELETE request that deletes an entire password type, along with the username and emails associated with the password
// PROTECTED function, accessibly only when authenticated
func deletePassword(w http.ResponseWriter, r *http.Request) {

}

// getLocks - GET request to get random locks to serve to the user
func getLocks(w http.ResponseWriter, r *http.Request) {

}

// storeTempLocks - POST request
// PROTECTED function, accessibly only when authenticated
// func storeTempLocks(w http.ResponseWriter, r *http.Request) {

// }

/**
=============================================
	Helper functions
=============================================
*/

// encryptAndStorePassword
func encryptAndStorePassword(url string, email string, password string) (err error) {
	preppedStatement, err := pwmdb.DB.Prepare("INSERT INTO passwords (url, email, password) values (?, ?, ?)")
	if err != nil {
		log.Panic(err)
	}

	_, err = preppedStatement.Exec(url, email, password)
	if err != nil {
		log.Panic(err)
		return err
	}

	return nil
}

func getEncryptedPassword(url string) (encryptedPassword EncryptedPassword, err error) {
	var email string
	var password string
	var locks string
	err = pwmdb.DB.QueryRow("SELECT email, password, locks FROM passwords WHERE url = ?", url).Scan(&email, &password, &locks)
	if err != nil {
		log.Panic(err)
		return EncryptedPassword{}, err
	}

	return EncryptedPassword{URL: url, Email: email, Password: password, Locks: locks}, nil
}

// connectPasswordDatabase - Helper function that connects to the correct database
func connectPasswordDatabase() {

}

// authenticateDatabase - Helper function that uses the dynamic authentication service to authenticate the user and gain access to the database
func authenticateDatabase() {

}

// encryptPassword - Helper function that encrypts a password
func encryptPassword(password string) {

}

// encryptPassword - Helper function that honeys(?) a password before encryption
func honeyPassword() {

}

// decryptPassword - Helper function that decrypts a password
func decryptPassword(encryptedPassword EncryptedPassword) (decryptedPassword Password, err error) {
	return Password{URL: encryptedPassword.URL, Email: encryptedPassword.Email, Password: encryptedPassword.Password}, nil
}
