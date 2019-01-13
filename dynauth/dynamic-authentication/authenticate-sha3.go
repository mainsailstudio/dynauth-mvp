/*
	Title:	Authentication package
	Author:	Connor Peters
	Date:	3/12/2018
	Desc:	Authenticates a user using the SHA3 hashing algorithm
*/

package dynauthcore

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	sha3 "golang.org/x/crypto/sha3"
)

// AuthenticateSHA3 - to perform command line SHA3 hash authentication for testing
func AuthenticateSHA3(locks string, otp string, userid string, iterations int) {
	// first prep auth for comparison
	toHash := locks + otp

	authenticated, err := compareAuthsWithSaltSHA3(toHash, userid)
	if err != nil {
		fmt.Println("Error encountered:", err)
	}

	if authenticated == true {
		fmt.Println("AUTHENTICATED")
	} else {
		fmt.Println("NO MATCH FOUND")
	}
}

// AuthenticateSHA3API - to perform authentication for the restful API
func AuthenticateSHA3API(userid string, auth string) (bool, error) {
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
