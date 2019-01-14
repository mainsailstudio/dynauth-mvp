/*
	Title:	JWT management for the API
	Author:	Connor Peters
	Date:	3/15/2018
	Desc:	Issues a valid JSON Web Token for authentication over the API
*/

package dynauthapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// issueJWT - takes in the user's email and returns a JWT with the email as a claim
// The default expiration is 1 week (168 hours), but this should be a dynamic value in the future
func issueJWT(email string) string {

	pkSecret, err := ioutil.ReadFile("../../../private.ppk") // in form of byte
	if err != nil {
		log.Fatal(err)
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "Dynauth Test",
		"email": email,
		//"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"iat": time.Now(),
		"exp": time.Now().Add(time.Hour * 168).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(pkSecret)
	if err != nil {
		fmt.Println("Error creating signed token")
		log.Fatal(err)
	}

	return tokenString
}
