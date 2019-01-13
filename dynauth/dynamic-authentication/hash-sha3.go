/*
	Title:	Hash SHA3 (with salts)
	Author:	Connor Peters
	Date:	2/24/2018
	Desc:	Takes in a "toHash" slice and returns the corresponding hashes
			Hash 2 options, 1 with salts and 1 without. It is recommended to use the one with salts wherever possible.
	NOTE:	This was implemented for sake of efficiency, Bcrypt is a better password hasher but this isn't a password baby
*/

package dynauthcore

import (
	"fmt"
	"io/ioutil"
	"log"

	sha3 "golang.org/x/crypto/sha3"
)

// HashPermsSHA3 - takes in the slice to hash for SHA3 and returns a completely hashed slice of strings.
func HashPermsSHA3(toHash []string) []string {
	hashed := []string{}

	// get the private key from file
	pkSecret, err := ioutil.ReadFile("../../../private.ppk") // in form of byte
	if err != nil {
		log.Fatal(err)
	}

	// iterates through toHash and hashes them all
	for i := 0; i < len(toHash); i++ {
		h := make([]byte, 32)
		d := sha3.NewShake256()
		// Write the key into the hash.
		d.Write(pkSecret)
		// Now write the data.
		d.Write([]byte(toHash[i]))
		// Read 32 bytes of output from the hash into h.
		d.Read(h)
		fmt.Printf("%x\n", h)

		hashString := fmt.Sprintf("%x\n", h)
		fmt.Println("Hash casted to string is", hashString)

		// add the new hash to the slice
		hashed = append(hashed, hashString)
	}
	return hashed
}

// HashPermsWithSaltSHA3 - takes in the slice to hash for SHA3 and returns a 2d slice of completely hashed strings and their corresponding salts.
func HashPermsWithSaltSHA3(toHash []string) [][]string {
	hashed := [][]string{}

	// get the private key from file
	pkSecret, err := ioutil.ReadFile("../../../private.ppk") // in form of byte
	if err != nil {
		log.Fatal(err)
	}

	// iterates through toHash and hashes them all
	for i := 0; i < len(toHash); i++ {
		h := make([]byte, 64)

		// generate salt
		s, err := getSaltString(32)
		if err != nil {
			log.Print(err)
		}

		// add salt to hash
		toHashWithSalt := toHash[i] + s
		// create a new SHA3 256
		d := sha3.NewShake256()
		// Write the key into the hash.
		d.Write(pkSecret)
		// Now write the data.
		d.Write([]byte(toHashWithSalt))

		// Read 32 bytes of output from the hash into h.
		d.Read(h)

		hashString := fmt.Sprintf("%x\n", h)
		//fmt.Println("Hash casted to string is", hashString)

		// add the new hash to the slice
		keyPair := []string{}
		keyPair = append(keyPair, hashString)
		keyPair = append(keyPair, s)
		hashed = append(hashed, keyPair)
	}
	return hashed
}
