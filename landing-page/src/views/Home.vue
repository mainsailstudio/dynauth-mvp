<template>
  <div class="row" id="main-content">
    <div class="col-lg-7">
      <h1 class="super-header">Dynauth</h1>
      <p class="subtitle">Identity, authentication, and authorization</p>
      <h3 class="title">Introduction</h3>
      <p>We want to provide users with a highly secure, self-owned, identity.</p>
      <p>By “identity”, we mean a single online master account that connects all your other accounts into one that automatically logs you on and keeps you safe.</p>
      <p>By “self-owned”, we mean that Dynauth does not own the data you provide, it just securely manages it. Data harvesting to serve targeted ads is not part of the business model.</p>
      <p>To keep things "highly secure", the online identity will be built on the idea of dynamic authentication and continual authentication.</p>
      <h3 class="title">Dynauth Authentication</h3>
      <p>Passwords suck. While not devoid of attractive features, passwords are a relic from the past that no longer provide a secure authentication session to users. Dynamic authentication is a complete alternative to passwords designed to be very similar in function, but more convenient and vastly more secure.</p>
      <p>Rather than describing how dynamic authentication works, why don't you give it a try? It's pretty simple! You have a table of "keys" that correlate to numbered "locks". Everytime you log in, there is a different random pattern of locks and you just type in the associated keys. We've already made a table of keys, why don't you try logging in to see how it works?</p>
      <p>We are currently working on writing a white paper describing dynamic authentication in detail, please check back here in the future for more info.</p>
      <h3 class="title">Continual Authentication</h3>
      <p>While dynamic authentication is (theoretically) a huge step-up from passwords, it still has a large problem. If a user logs into, say, a web app on their laptop and then leaves to get more coffee, the user is no longer the point of authentication -- the laptop is. Any bad actor could potentially impersonate the user simply by gaining physical or remote access to that laptop.</p>
      <p>The idea behind continual authentication is that by continually monitoring an authentication session through various usage metrics like keyboard typing speed, the user will always remain the point of authentication rather than the device.</p>
      <p>This is still in early stages of development, so we don't have anything to show you this quite yet! But trust us, we're working on it.</p>
      <h3 class="title">In the meantime...</h3>
      <p>Give dynamic authentication a try!</p>
    </div>
    <div class="col-lg-5">
      <h2 class="mt-5">Give it a test run:</h2>
      <label>Email:</label>
      <input type="email" name="email" id="email" placeholder="example@example.com"><br />
      <label>Locks:&nbsp;&nbsp;<span class="locks">{{ lockString }}</span>&nbsp;&nbsp;&nbsp;<i class="fas fa-sync fa-xs" id="refresh-icon" v-on:click="generateLocksAndKeys"></i></label>
      <input type="text" name="user-input" placeholder="answer from vue here" id="user-input">
      <p class="subtitle">Your email will be added to a database of "alpha" users.<br />
      Just leave it empty if you don't want to be added.</p>
      <button v-on:click="authenticate">Login!</button>
      <div class="mt-4 alert alert-success" role="alert" id="success-alert">
        <h4 class="alert-heading">Success!</h4>
        <p>You got it! Questions? Concerns? Email: <a href="mailto:design@mainsailstudio.com">design@mainsailstudio.com</a></p>
        <hr>
      </div>
      <div class="mt-4 alert alert-danger" role="alert" id="failure-alert">
          <h4 class="alert-heading">The keys you entered are not correct!</h4>
          <p>Go ahead and give it another try. If you're confused as to what this is, please read about it and then re-try.</p>
          <hr>
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
export default {
  data: ()=> ({
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
    lockNum: 4,
    answerKeyArray: [],
    displayLocks: [],
    lockString: "",
  }),
  mounted(){
    // Initiate the locks and keys
    this.generateLocksAndKeys();

  },
  methods: {
    // Generate the locks and keys and display them to the user
    generateLocksAndKeys: function(){
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
        this.lockString = this.displayLocks.join(" - ");

        // Reset the user input box to be empty on refresh
        var userInput = document.getElementById("user-input");
        userInput.value = "";

        // Give the user a hint by setting the placeholder for the input to be the correct answer
        userInput.placeholder = this.answerKeyArray.join("");
    },

    // Actually authenticate
    authenticate: function(){

        var userInput = document.getElementById("user-input").value;
        var correctAnswer = "";

        // convert answer key array into a singular string to compare to the userInput string
        for(var i = 0; i < this.answerKeyArray.length; i++) {
            correctAnswer += this.answerKeyArray[i];
        }

        if(userInput == correctAnswer){
            console.log("\n\n+++ User input, " + userInput + " is correct when compared to " + correctAnswer + " +++\n\n");
            this.formSuccess();
        } else {
            console.log("\n\n--- User input, " + userInput + " is INCORRECT when compared to " + correctAnswer + " ---\n\n");
            this.formFailure();
        }

        // email
        if($("#email").val()){
            this.submitForm();
        } else {
            // Silence?
        }

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

    submitForm: function(){
      // Initiate Variables With Form Content
      var email = $("#email").val();

        $.ajax({
              type: "POST",
              url: "https://dynauth.io/dynauth.php",
              data: "email=" + email,
              success : function(text){
                  if (text == "success"){
                      this.formSuccess();
                  }
              }
        });
    },

    formSuccess: function(){
        $( "#failure-alert" ).hide();
        $( "#success-alert" ).slideDown();
    },

    formFailure: function(){
        $( "#success-alert" ).hide();
        $( "#failure-alert" ).slideDown();
        this.generateLocksAndKeys();
    }
  } 
}
</script>