<template>
  <div class="row" id="main-content">
    <div class="col-lg-7">
        <h1 class="less-super-header">Dynamic Authentication</h1>
        <p class="subtitle">An iterative improvement on passwords</p>
        <h3 class="title">Introduction</h3>
	    <p>The reliability of passwords as a secure authentication scheme has been degrading rapidly since their digital inception in the very first multi-user operating systems. Due to Moore's Law and the general expansion of the Internet, passwords that were once secure are no longer, and the bar for what is an acceptable password is continually being pushed farther down a dark path.</p>
	
	    <p>The user experience when authenticating with a password has followed an identical route, becoming more and more unbearable with time. In short, as computers get faster, passwords must be more complex, and therefore user requirements more stringent. However the recall ability of humans has not changed at all with time, and most people's passwords are insecure regardless of the requirements.</p>
	
	    <p>This brings us to an interesting place in the world of cybersecurity. Users need to be authenticated of course, yet alternate schemes that would attempt to replace the password, such as biometric readers, have still not gained widespread prevalence across the Internet despite their growing ubiquity in end-user computing devices. Passwords are increasingly unfit for the job, yet there is no replacement.</p>
	
	    <p>But it's not all bad. Passwords <em>have</em> sorta worked, at least enough to protect the majority of Internet users from malevolent cyber threats. Let's give passwords a break for just a second to appreciate the benefits they provide.</p>
        <ul>
            <li>Passwords are platform agnostic. Any device with an input method can take advantage of them.</li>
            <li>Passwords do not require any additional hardware pieces for the users or the developers.</li>
            <li>Passwords are free to use.</li>
            <li>Passwords are relatively easy to implement for developers.</li>
            <li>Passwords can be reset.</li>
            <li>Passwords are marvelously familiar to users due to their prevalence on the Internet.</li>
        </ul>
	    <p>When framed like that, passwords actually seem to be quite competent at their job of authentication. Yet time and time again, the opposite is proven to be true. The reality is that the critical flaw to passwords are humans themselves, not the technology they are built on. But technology is here to suit humans; if it doesn't work for us, it just doesn't work.</p>
	
        <h3 class="title">An Improvement/Alternative</h3>
	    <p>The core concept behind dynamic authentication is to create an authentication scheme that attempts to retain the aspects of passwords that are beneficial, while eliminating the parts that aren't. Therefore, dynamic authentication should be thought of as a sort of iterative improvement on typical passwords, rather than a new thing entirely.</p>
	
        <h3 class="title">Dynamic Authentication: The Technical</h3>
        <p>Dynauth, just like a password-based authentication scheme, requires an identifying username as well as an authenticating secret. The difference is that the secret the user must remember for authentication is not a password. Instead, the secret is a list of words (typically plain English words), known as "keys", associated with numbers, known as "locks". Look at the list of keys on this page for an example of what they might look like.</p>

        <p>However, unlike password-based authentication, the user does not enter in the entirety of their secret to authenticate. Instead, after they enter in their username, the user is presented with 4 of their (numbered) locks in a random order, without repeat. The user then inputs the keys that are associated with the presented locks as one long string, in the same order as shown, without spaces or delimiters. Should the keys and locks match the ones on the list, the user will be authenticated.</p>

        <p><em>So why is this an improvement?</em> The storage design</p>        
        <p>This is where things get interesting. For all the same reasons as passwords, all user's locks and keys need to be hashed and stored in a safe and secure manner to prevent attackers from accessing cleartext entries of them if they can somehow get in the backend of the system. There also needs to be a way guarantee the security of each hash is greater than that of a typical password because otherwise, what's the point?</p>
        
        <p>This presents a problem: any normal dictionary word cannot be hashed by itself and referenced later as it would be far too insecure. Therefore, the locks and keys can't just be hashed and stuck in the database to be compared to the individually entered locks and keys later.</p>
        
        <p>The core difference that allows dynauth to operate more securely is the hashed storage of <em>all possible lock and key permutations</em>. This means that if a user configures 10 total keys, and are presented 4 total locks at the time of authentication (the base level configuration I chose), there will be a total of 10P4 (10 * 9 * 8 * 7 = 5040) permutations generated and stored.</p>
        
        <p><strong>Here is an example of a user's permutations being generated and hashed:</strong></p>
        <ol>
            <li>The first 4 keys of the user's configuration are concatenated: <strong>greatkingdecidedthat</strong></li>
            <li>That string is then hashed: <strong style="word-break: break-all;">1CF0B384D1D52133255970AE0B091D5BDFCB627FEA9048D1FBC265BBF00137B7T</strong></li>
            <li>The locks that those 4 keys are associated with are prepended to the hash string: <strong style="word-break: break-all;">12341CF0B384D1D52133255970AE0B091D5BDFCB627FEA9048D1FBC265BBF00137B7T</strong></li>
            <li>That entire string is then hashed again, and the result is what is stored in the database as a single permutation: <strong style="word-break: break-all;">0E60D213A1055A3F3D49BF4611D3307542615E53AT</strong></li>
            <li>This process would continue until all possible permutations of the user's configuration are generated and stored</li>
        </ol>
        
        <p><strong>Here is the process for a user to authenticate:</strong></p>
        <ol>
            <li>The user enters in their email on the client side and the email is sent to the server</li>
            <li>The server randomly selects the locks appropriate for the user and stores them in a database with an expiration date and time</li>
            <li>The server then sends the same locks it stored back to the client for the user to view</li>
            <li>The user enters in the keys associated with the locks</li>
            <li>The locks are hashed client side and the hash is then sent to the server</li>
            <li>The server then prepends the keys stored previously in step 2 to the hash received from the client and hashes that entire string again</li>
            <li>The resulting hash is then used to iterate over the user's database of lock and key permutations until a match is found. If any permutation matches, the user is authenticated</li>
        </ol>
        <h3 class="title">Main Benefits</h3>
        <h4>Crack time is increased</h4>
        <p>The largest benefit dynauth provides is how much longer it would take to successfully crack/bruteforce. The average password provides about 2^22 bits of entropy. Considering a worse case scenario, each key present in dynauth provides between 2^11 and 2^14 bits of entropy, depending on the words present in the dictionary used for a Dictionary Attack. With a 10x4 schema, that means the average dynauth setup provides between 2^44 and 2^56 bits of entropy. These bits of entropy would also provide more protection than a typical password since they are hashed twice, first on the client side, then on the server side with the locks. This would mean an attacker would need to perform twice as many operations per guess, doubling the average amount of computation time needed to crack a single hash.</p>
        <h4>Keylogging is harder</h4>
        <p>Due to the fact that the keylogging system won't know which keys they retrieved are associated with which locks are displayed on the screen, keylogging is much more intensive. It is still possible, but does not provide immediate access to the user's account.</p>
        <h4>It's easier to use</h4>
        <p>Since the system is designed to be secure even using normal dictionary words, people won't have to remember the symbols and capitals they used to be required to add to increase entropy. That means it's considerably easier to type on a mobile phone as well.</p>
        <h3 class="title">Notes</h3>
        <p>The implementation here is crude and doesn't actually use a server to query against - it's all browser based. This is simply to illustrate how things function as if there were a server. It's slow because browsers and JavaScript are just not well optimized compared to a server side language like Go. Also, the answer won't be hinted at you in the placeholder text of course, that is just for demonstration purposes.</p>
    </div>
    <div class="col-lg-5">
        <h2 class="mt-5">Watch it in action</h2>
        <p>We've already made a sample account for you to test. Start by typing in the email associated with the account</p>
        <label for="email">Email:</label>
        <input type="email" name="email" v-model.trim="email" placeholder="test@dynauth.io"><button class="next" v-on:click="getLocksFromEmail" title="Next!"><i class="fas fa-angle-double-right"></i></button>
        <label v-if="locksExist">Locks:&nbsp;&nbsp;<span class="locks">{{ lockString }}</span>&nbsp;&nbsp;&nbsp;<i class="fas fa-sync fa-xs" id="refresh-icon" v-on:click="generateLocksAndKeys"></i></label>
        <input v-if="locksExist" type="text" name="user-input" placeholder="loading..." id="user-input">
        <button v-if="locksExist" class="next" v-on:click="authenticate" title="Next!"><i class="fas fa-angle-double-right"></i></button>
        <div class="console">
            <h4>Console:</h4>
            <p>{{ consoleOutput }}</p>
        </div>
        <h2 class="mt-5 mb-5" data-toggle="collapse" data-target="#keys" aria-expanded="false" aria-controls="keys">Keys &nbsp;<i class="far fa-plus-square" id="plus-icon"></i></h2>
        <div class="collapse show" id="keys">
            <div class="card card-body mb-5">
            <ol class="split">
                <li>great</li>
                <li>king</li>
                <li>decided</li>
                <li>that</li>
                <li>kale</li>
                <li>was</li>
                <li>not</li>
                <li>healthy</li>
                <li>vegetable</li>
                <li>anymore</li>
            </ol>
            </div>
        </div>
    </div>
  </div>
</template>
<script>
import axios from 'axios';
import crypto from "crypto";

export default {
  data: ()=> ({
    // API boilerplate variables
    loading: true,
    consoleOutput: "> Waiting for input", // this will be used to depict what's going on with the API during the process

    // Specific variables
    keyArray: [
        'great',
        'king',
        'decided',
        'that',
        'kale',
        'was',
        'not',
        'healthy',
        'vegetable',
        'anymore'
    ],
    lockArray: [
        1,
        2,
        3,
        4,
        5,
        6,
        7,
        8,
        9,
        10
    ],
    lockNum: 4,
    answerKeyArray: [],
    displayLocks: [],
    email: "",
    locksExist: false,
    lockString: null,
    keys: null,
    secret: null,
    success: null,
  }),
  mounted(){
  },
  methods: {
    // Start by getting the server to display the user's locks
    getLocksFromEmail: function(){
        if(this.email !== "test@dynauth.io"){
            this.consoleOutput = "> The account associated with this email does not exist"
            this.email = "";
        } else {
            this.generateLocksAndKeys();
        }

        // axios
        //     .get('/auth/users/locks', {
        //         params: {
        //             email: this.email
        //         }
        //     })
        //     .then(response => {
        //         this.locks =  response.data.locks
        //         this.consoleOutput += "\n> response: The user's locks were loaded as: "
        //         this.consoleOutput += response.data
        //         this.consoleOutput += "\n> Those locks have been temporarily stored on the server side for authentication"
        //     })
        //     .catch(error => {
        //         console.log(error)
        //         this.consoleOutput += "\n> response: ERROR"
        //         this.consoleOutput += error.response.status
        //         this.consoleOutput += error.response.data
        //         this.locks =  "for the win"
        //     })
        //     .finally(() => this.loading = false)
    },

    // Generate the locks and keys and display them to the user
    generateLocksAndKeys: function(){
        $(".console").css("background-color", "#333")
        this.consoleOutput = "> Loading user's locks from the API..."
        this.answerKeyArray = []; // Reset the correct answer array to being empty
        var allLocks = []; // All of the locks, derived from the index of the answer array
        var chosenLocks = []; // The "chosen" locks after they are shuffled and cut down to the right lockNum size
        this.displayLocks = []; // The chosenLocks with 1 added to each index

        // Generate the list of locks from the length of the key array
        for(var i = 0; i < this.keyArray.length; i++) {
            allLocks.push(i);
        }

        // Then shuffle that array randomly
        this.shuffle(allLocks);

        // Then take only the locks from index 0 to lockNum (typically 4)
        // This process removes the possiblity of random repeats when generating the locks
        var chosenLocks = allLocks.slice(this.keyArray.length - this.lockNum);

        for(i = 0; i < this.lockNum; i++) {
            var key = this.keyArray[chosenLocks[i]];
            this.answerKeyArray.push(key);

            // Add 1 to the index to display the correct locks
            this.displayLocks.push(chosenLocks[i]+1);
        }

        // Display the locks to the user
        this.locksExist = true;

        var self = this;
        setTimeout(function(){
            self.lockString = self.displayLocks.join(" - ");
            // Reset the user input box to be empty on refresh
            var userInput = document.getElementById("user-input");
            userInput.value = "";

            // Give the user a hint by setting the placeholder for the input to be the correct answer
            userInput.placeholder = self.answerKeyArray.join("");
            self.consoleOutput = "> Locks were loaded from the server API and stored in the server temporarily"
        }, 400);
    },

    // Taken off the Internet, thanks guys!
    shuffle: function(array) {
        var i = array.length,
            j = 0,
            temp;

        while (i--) {

            j = Math.floor(Math.random() * (i+1));

            // swap randomly chosen element with current element
            temp = array[i];
            array[i] = array[j];
            array[j] = temp;

        }

        return array;
    },

     // Next, try to authenticate
    authenticate: function(){
        this.updateConsole("\n> Hashing the user's keys on the client side...")
        // this.consoleOutput = "> Hashing the user's keys on the client side..."
        var userInput = document.getElementById("user-input").value;

        // Hash secret on client side
        const hash = crypto.createHash('sha256')
        hash.update(userInput)

        this.secret = hash.digest('hex')
        this.updateConsole("\n> The user's keys were hashed using SHA256 and are now: " + this.secret)
        // this.consoleOutput += "\n> The user's keys were hashed using SHA256 and are now: " + this.secret

        // Now simulate the server side authentication for the benefit of the technical user
        this.simulateServerSideAuthentication()

        // axios
        //     .get(this.testApiURL + '/auth/users/keys', {
        //         params: {
        //             secret: this.secret
        //         }
        //     })
        //     .then(response => {
        //         this.locks = response.data.bpi
        //     })
        //     .catch(error => {
        //         console.log(error)
        //     })
        //     .finally(() => this.loading = false)
    },

    generateAnswerKeyArray: function(){
        for(i = 0; i < this.lockNum; i++) {
            var key = this.keyArray[chosenLocks[i]];
            this.answerKeyArray.push(key);

            // Add 1 to the index to display the correct locks
            this.displayLocks.push(chosenLocks[i]+1);
        }
    },

    simulateServerSideAuthentication: function(){
        // permutate the locks and keys
        const keyPermArray = this.generateLimPerms(this.keyArray, this.lockNum);

        // hash the lock and key combos
        const hashArray = this.hashPermsSHA256(keyPermArray);
    },

    // combinePerms - to correctly concat 2 slices of permutations into 1.
    // Needs 2 slices, 1 of locks and 1 of keys. It is assumed that they match up perfectly and will result in logic errors if they do not.
    combinePerms: function(locks, keys){
        const combined = []; // assumes locks and keys are at the same index (SHOULD ALWAYS BE)
        for (let i = 0; i < locks.length; i++) {
            const combineString = locks[i] + keys[i];
            combined.push(combineString);
        }
        return combined;
    },

    // takes in the 'toHash' array, hashes each item and then returns a new array with the hashes
    hashPermsSHA256: function(toHash) {
        const currentLocks = this.displayLocks.join("")
        this.updateConsole("\n> The user input has been sent to the server")
        // this.consoleOutput += "\n> The user input has been sent to the server"
        this.updateConsole("\n> The locks '" + currentLocks + "' were stored on the server last")
        // this.consoleOutput += "\n> The locks: " + currentLocks + " were stored on the server last.\n"
        const hashedKeys = [];
        for (let i = 0; i < toHash.length; i++) {
            const hash = crypto.createHash('sha256')
            hash.update(toHash[i])
            const hashString = hash.digest('hex')

            hashedKeys.push(hashString);
        }
        const lockPerms = this.generateLimPerms(this.lockArray, this.lockNum)
        const combinedPerms = this.combinePerms(lockPerms, hashedKeys)

        const finalClientString = currentLocks.concat(this.secret)
        this.updateConsole("\n> The locks have been concatenated to the user input hash and is now: " + finalClientString)
        // this.consoleOutput += "\n> The locks have been concatenated to the user input hash and is now: " + finalClientString

        const hash2 = crypto.createHash('sha256')
        hash2.update(finalClientString)
        const finalClientHash = hash2.digest('hex')

        this.updateConsole("\n> The user input hash string has been hashed again and is now in it's final form: " + finalClientHash)
        // this.consoleOutput += "\n> The user input hash string has been hashed again and is now in it's final form: " + finalClientHash
        this.updateConsole("\n\n> The string is now being compared to all the permutations in the database... ")
        // this.consoleOutput += "\n> The string is now being compared to all the permutations in the database... "

        for (let i = 0; i < combinedPerms.length; i++) {
            if(finalClientString === combinedPerms[i]){
                this.success = true;
                this.authenticateResult(this.success);
                return
            }
        }
        this.success = false;
        this.authenticateResult(this.success);

        return;
    },

    // Generates the subsets of the passed array, and then permutes each subset
    generateLimPerms: function (toPerm, num) {
        const subsets = this.getSubsets(toPerm, num);
        const perms = this.getPerms(subsets);
        return perms;
    },

    // Generates each unique subset of the array that is passed in, with the num being the limiting factor
    getSubsets: function(locks, num){
        const res = [];
        const findSubset = [];

        // helper 'fat arrow' lambda function
        const helper = (arr, subset, n) => {
            if (subset.length === num) {
                const tmp = subset.slice();
                res.push(tmp);
                return;
            }

            if (n === arr.length) {
                return;
            }

            subset.push(arr[n]);
            // recursion starts here
            helper(locks, subset, n + 1);
            subset.pop();
            helper(locks, subset, n + 1);
        };

        helper(locks, findSubset, findSubset.length);
        return res;
    },

    // A javascript implementation of Heap's algorithm
    heapPermutation: function(arr){
        const res = [];

        const helper = (array, n) => {
            if (n === 1) {
                const tmp = array.join('').slice();
                res.push(tmp);
            } else {
                for (let i = 0; i < n; i++) {
                    // recursion people
                    helper(array, n - 1);
                    if (n % 2 === 1) {
                        const tmp = array[i];
                        array[i] = array[n - 1];
                        array[n - 1] = tmp;
                    } else {
                        const tmp = array[0];
                        array[0] = array[n - 1];
                        array[n - 1] = tmp;
                    }
                }
            }
        };
        helper(arr, arr.length);
        return res;
    },

    // This is the organizer for heap's algo
    // It takes a 2d array of each subset (generated above) and returns a single array
    // With each of heap's permutations as a single string with no delimiter (joined)
    getPerms: function(perms){
        const res = [];
        for (let i = 0; i < perms.length; i++) {
            const list = perms[i].slice();
            const tmp = this.heapPermutation(list).slice();
            for (let j = 0; j < tmp.length; j++) {
                res.push(tmp[j]);
            }
        }
        return res;
    },

    authenticateResult: function(success){
        var self = this;
        setTimeout(function(){ 
            if(success){
                self.updateConsole("\n\n> Success!\nThe user input hash above matched one of the permutations in the database, just like designed!")
                $(".console").css("background-color", "green")
            } else {
                self.updateConsole("\n\n> Failure :( \nNone of the hashes matched!")
                $(".console").css("background-color", "#d30202")
            }
            $(".console").animate({ scrollTop: $(document).height() }, "slow");
        }, 200);
    },
    // updateConsole: function(string){
    //     var self = this;
    //     setTimeout(function(){ 
    //         self.consoleOutput += string
    //         $('.console').scrollTop($('.console')[0].scrollHeight);
    //     }, 2500);
    // },
    updateConsole: function(string){
        this.$nextTick(()=>{
            this.consoleOutput += string
            $(".console").animate({ scrollTop: $(document).height() }, "slow");
        })
    },
  } 
}
</script>