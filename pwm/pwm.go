/*
	Title:	Password Manager Service Manager
	Author:	Connor Peters
	Date:	12/28/2018
	Desc:	Just starts the local service for now! Eventually there will be more features in here like a CLI configuration UI or something
*/

package main

import pwmapi "github.com/mainsailstudio/dynauth-mvp/pwm/pwmapi"

func main() {
	pwmapi.Start()
	// pwmapi.Start("https")
} // end of main

// initialize - Starts a small CLI program to enable the user to create a  local dynauth account with locks and keys
// Can only be run if the user hasn't created a local account yet
func initialize() {

}
