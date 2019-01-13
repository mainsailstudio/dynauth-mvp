package dynauthdev

import (
	dynauthconst "dynauthconst"
	dynauthcore "dynauthcore"
	"encoding/json"
	"fmt"
)

func createUser() {
	// ask for num of lock/key combos (testing purposes)
	fmt.Println("The default number of lock/key combos as defined by the keyNum constant are", dynauthconst.KeyNum)
	fmt.Println("For testing purposes please enter in the number of lock/key combos you want: ")
	var numHash int
	fmt.Scan(&numHash)

	// get user data
	var userid string
	var fname string
	var lname string
	var email string
	var phone string
	var securityLv string
	fmt.Println("For testing purposes please enter in the userid (randomize please): ")
	fmt.Scan(&userid)
	fmt.Println("For testing purposes please enter in the firstname: ")
	fmt.Scan(&fname)
	fmt.Println("For testing purposes please enter in the lastname: ")
	fmt.Scan(&lname)
	fmt.Println("For testing purposes please enter in the email: ")
	fmt.Scan(&email)
	fmt.Println("For testing purposes please enter in the phone: ")
	fmt.Scan(&phone)
	fmt.Println("For testing purposes please enter in the security level: ")
	fmt.Scan(&securityLv)

	// store user data
	dynauthcore.StoreUserInfo(userid, fname, lname, email, phone, securityLv)
	fmt.Println("User info was stored")

	// initialize the 2d lock-key combo slice
	lockSlice := make([]string, numHash)
	keySlice := make([]string, numHash)

	// for loop that asks for locks and keys (testing purposes)
	for i := 0; i < numHash; i++ {
		// intialize the slice of this particular iteration of lockKeySlice
		var lock string
		var key string

		// start by getting the lock and putting it into the slice
		fmt.Print("Enter in lock number ", i, ": ")
		fmt.Scan(&lock)
		fmt.Println("Lock is: " + lock) // print lock
		lockSlice[i] = lock

		// next get the key and put it into the slice
		fmt.Print("Enter in key correlating to lock number ", i, ": ")
		fmt.Scan(&key)
		fmt.Println("Key is: " + key) // print lock
		keySlice[i] = key
	}

	// store locks
	dynauthcore.StoreLocks(lockSlice, userid, "1") // using 1 as the locktype since it does nothing currently
	fmt.Println("User's locks were stored")

	// create and store auths
	lockPerms := dynauthcore.LimPerms(lockSlice, dynauthconst.DisplayLockNum) // create the limited permutations for the locks from the dynauthcore permutations.go package
	keyPerms := dynauthcore.LimPerms(keySlice, dynauthconst.DisplayLockNum)
	keyPermsJSON, _ := json.Marshal(keyPerms)
	fmt.Println("Lim perms are", string(keyPermsJSON))
	permsToHash := dynauthcore.CombinePerms(lockPerms, keyPerms) // create the perms to hash (should most likely be in a package eventually)
	fmt.Println("Perms to hash is", permsToHash)
	fmt.Println("Total number of permutations is", len(permsToHash))
	hashedPermsWithSalt := dynauthcore.HashPermsWithSaltSHA3(permsToHash)
	//hashPerms := dynauthcore.HashPermsSHA3(permsToHash)
	//hashPerms := dynauthcore.HashPermsScrypt(permsToHash)
	//hashPerms := dynauthcore.HashPermsBcrypt(permsToHash)
	//fmt.Println("Hashed perms is", hashedPermsWithSalt)
	//fmt.Println("Let's try to store them!")
	dynauthcore.StoreAuthsWithSalts(hashedPermsWithSalt, userid)
	//dynauthcore.StoreAuthsPlain(hashPerms, userid)
	fmt.Println("User's auths were stored")
	fmt.Println("Total number of permutations is", len(permsToHash))

	//hash(lockKeySlice)
}
