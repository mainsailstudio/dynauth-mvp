/**
Title:	Password Manager
Author:	Connor Peters
Date:	12/29/2018
Desc:	Password manager secured with dynamic authentication and honey encryption
Notes:	THIS IS NOT CONSIDERED A SECURE PASSWORD MANAGER YET. This product is in it's pre-alpha stage and isn't really designed to be used by anyone yet, it's mostly for testing purposes
*/

// Package pwmapi creates all the routes and the corresponding functions for the password manager to be used "RESTfully"
package pwmapi

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
func Start(protocol string) {
	// InitDB()
	// Init router
	fmt.Println("Starting New Password Manager API Service")
	r := mux.NewRouter()

	// Init the routes

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

	// Add the default handler to allow CORS
	handler := cors.Default().Handler(r)

	// Listen and serve the handler based on the desired protocol
	// Everything other than specified "http" will be run over TLS
	switch protocol {
	case "http":
		// Serve over normal HTTP, not HTTPS
		http.ListenAndServe(":8080", handler)
	default:
		// Serve over HTTPS (TLS)
		startSecureAPI(handler)
	}

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

/**
=============================================
	Helper functions
=============================================
*/

// startSecureAPI - determine if TLS is configured and run, otherwise initialize TLS
func startSecureAPI(handler http.Handler) {
	if _, err := os.Stat("ca.crt"); os.IsNotExist(err) {
		initTLS()
	}
	if _, err := os.Stat("ca.key"); os.IsNotExist(err) {
		initTLS()
	}

	signTLS()

	// Listen and serve the API over TLS (HTTPS)
	err := http.ListenAndServeTLS(":443", "bob.crt", "bob.key", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// initTLS - Initialize the certificate and private key components to start the API over TLS
func initTLS() {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"DYNAUTH_PASSWORD_MANAGER_LOCAL_ENCRYPT"},
			Country:       []string{"GLOBAL"},
			Province:      []string{"PROVINCE"},
			Locality:      []string{"CITY"},
			StreetAddress: []string{"ADDRESS"},
			PostalCode:    []string{"POSTAL_CODE"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey
	cab, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		log.Println("create ca failed", err)
		return
	}

	// Public key
	certOut, err := os.Create("ca.crt")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cab})
	certOut.Close()
	log.Print("written initial cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("ca.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Print("written initial key.pem\n")
}

func signTLS() {

	// Load CA
	catls, err := tls.LoadX509KeyPair("ca.crt", "ca.key")
	if err != nil {
		panic(err)
	}
	ca, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		panic(err)
	}

	// Prepare certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"DYNAUTH_PASSWORD_MANAGER_LOCAL_ENCRYPT"},
			Country:       []string{"GLOBAL"},
			Province:      []string{"PROVINCE"},
			Locality:      []string{"CITY"},
			StreetAddress: []string{"ADDRESS"},
			PostalCode:    []string{"POSTAL_CODE"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey

	// Sign the certificate
	cert_b, err := x509.CreateCertificate(rand.Reader, cert, ca, pub, catls.PrivateKey)

	// Public key
	certOut, err := os.Create("bob.crt")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert_b})
	certOut.Close()
	log.Print("written SIGNED cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("bob.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Print("written SIGNED key.pem\n")

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
