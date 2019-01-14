/*
	Title:  Dynauth MVP Starter
	Author:	Connor Peters
	Date:	1/11/2018
	Desc:	Starts the dynauth api and the pwm api on the correct ports to launch the entire local service
*/

package main

import (
	pwm "github.com/mainsailstudio/dynauth-mvp/pwm"
)

func main() {
	// Start the dynauth service on the default port
	// dynauth.Start()

	// Start the pwm service on the default port
	pwm.Start()
}
