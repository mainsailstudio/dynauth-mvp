/*
	Title:	Serve a REST-like API for a usability test using Mux
	Author:	Connor Peters
	Date:	3/17/2018
	Desc:	This is the test controller of the application that presents an api over TLS
	NOTES:	For testing purposes, the negroni middleware has to be subbed out for a normal handlefunc to prevent CORS errors
*/

package dynauthapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// StartAPI - Start the api for mux api and serve it over TLS
func StartAPI() {

	// Init router
	fmt.Println("Starting Dynauth Alpha API")
	r := mux.NewRouter()

	// Default hello world landing page for API
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	// Basic handles for API
	r.HandleFunc("/auth", AuthenticateUser).Methods("POST")

	// Run the handler through CORS for testing
	handler := cors.Default().Handler(r)

	// Listen and serve the API over HTTP
	// http.ListenAndServe(":8080", handler)

	// Listen and serve the API over TLS (HTTPS)
	err := http.ListenAndServeTLS(":443", "gitignore/certificate.crt", "gitignore/private.key", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
