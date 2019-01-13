/*
	Title:	Authenticate User
	Author:	Connor Peters
	Date:	ORIGINALLY 3/17/2018, modified 12/23/2018
	Desc:	Simply authenticates a user based on their email and hashes of the keys/locks
*/

package dynauthapi

import (
	"dynauthcore"
	"encoding/json"
	"fmt"
	"net/http"
)

/**
* This function
 */
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var user UserLoginCheck
	_ = json.NewDecoder(r.Body).Decode(&user)
	if user.Email == "" {
		http.Error(w, "The email is empty", 400)
	}

	userExists, userID, err := checkUserExists(user.Email)
	if err != nil {
		http.Error(w, "Error encountered when checking if the user exists", 500)
	}

	if userExists == true {
		fmt.Println("User login state is", user.LoginState)
		if user.LoginState == "3" {
			// if the user is logging in via dynauth
			authCorrect := dynauthcore.AuthenticateBcryptAPI(userID, user.Secret)
			if authCorrect {
				fmt.Println("Correctly authenticated via auths")
				token := issueJWT(user.Email) // sending the user's email to be a part of the jwt claim
				userSuccess := UserLoginSuccess{ID: userID, Email: user.Email, LoginState: user.LoginState, Token: token}
				json.NewEncoder(w).Encode(userSuccess)
			} else {
				fmt.Println("NOT authenticated via auths")
			}
		}

		if user.LoginState == "1" {
			// if the user is loggin in via a temp pass
			passwordCorrect, err := dynauthcore.TempPassAuth(userID, user.Secret)
			if err != nil {
				http.Error(w, "There was an error encountered when authenticating the user with their password", 500)
			}

			if passwordCorrect {
				fmt.Println("Correctly authenticated via a temp pass")
				token := issueJWT(user.Email) // sending the user's email to be a part of the jwt claim
				userSuccess := UserLoginSuccess{ID: userID, Email: user.Email, LoginState: user.LoginState, Token: token}
				json.NewEncoder(w).Encode(userSuccess)
			} else {
				fmt.Println("NOT authenticated via a temp pass")
			}
		}
	}
}
