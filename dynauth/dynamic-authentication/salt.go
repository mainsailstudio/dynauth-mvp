/*
	Title:	Simple salting baby
	Author:	Connor Peters
	Date:	2/24/2018
	Desc:	2 simple functions to generate salts, one is a byte one is just a string
*/

package dynauthcore

import (
	"crypto/rand"
	"errors"
	"fmt"
)

// getSalt - generates a random byte salt with a saltLen, minimum length is 32
func getSalt(saltLen int) ([]byte, error) {
	// check to make sure the salt length is AT LEAST 32 bytes
	if saltLen < 32 {
		return nil, errors.New("Salt: salt length has to be at least 32")
	}

	salt := make([]byte, saltLen)
	if _, err := rand.Reader.Read(salt); err != nil {
		return nil, errors.New("Salt: salt creation failed for some reason, sorry I can't be more specific")
	}

	return salt, nil
}

// getSaltString - generates a random byte salt with a saltLen, minimum length is 32
// casts the byte into a string for convenience sake
func getSaltString(saltLen int) (string, error) {
	// check to make sure the salt length is AT LEAST 32
	if saltLen < 32 {
		return "", errors.New("Salt: salt length has to be at least 32")
	}

	salt := make([]byte, saltLen)
	if _, err := rand.Reader.Read(salt); err != nil {
		return "", errors.New("Salt: salt creation failed for some reason, sorry I can't be more specific")
	}

	saltString := fmt.Sprintf("%x", salt)
	return saltString, nil
}
