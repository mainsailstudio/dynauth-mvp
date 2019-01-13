/*
	Title:	All constants
	Author:	Connor Peters
	Date:	2/12/2018
	Desc:	This is the package that defines all of the application constants. This should be modified eventually to pull in from a settings database
*/

package dynauthconst

// DatabaseName - the name of the database
const DatabaseName = ""

// DatabaseUser - the user of the database
const DatabaseUser = ""

// DatabasePass - the password of the database
const DatabasePass = ""

// BcryptIterations - the amount of iterations to perform when hashing a password or auth using bcrypt
const BcryptIterations = 1

// ScryptIterations - the amount of iterations to perform when hashing a password or auth using scrypt
const ScryptIterations = 16384

// KeyNum - the total amount of keys the user will want, be careful with this because the larger this number is, the factorially larger the amount of computations will be. Keep it > 30
const KeyNum = 10

// DisplayLockNum - the total amount of locks that will be displayed for dynamic authentication. Keep it small (> 7)
const DisplayLockNum = 4

// Timezone for date/time
// const Timezone = time.LoadLocation("EST")
