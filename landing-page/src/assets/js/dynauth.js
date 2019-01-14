/**
 * Author: Connor Peters
 * 
 * Description: This simple JavaScript implementation of Dynauth is JUST FOR PRESENTATION PURPOSES
 *  It does not utilize the actual Dynauth API and is not intended to actually provide any sort of security benefits
 * 
 * It is JUST FOR PRESENTATION PURPOSES
 * 
 * Don't hate please.
 * 
 */
(function($) {
    "use strict"; // Start of use strict
    

    var keyArray = [
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
    ]

    var lockNum = 4; // The number of locks displayed to users

    var answerKeyArray = []; // The correct answer array

    function generateLocksAndKeys(){
        answerKeyArray = []; // Reset the correct answer array to being empty
        var allLocks = []; // All of the locks, derived from the index of the answer array
        var chosenLocks = []; // The "chosen" locks after they are shuffled and cut down to the right lockNum size
        var displayLocks = []; // The chosenLocks with 1 added to each index


        // Generate the list of locks from the length of the key array
        for(i = 0; i < keyArray.length; i++) {
            allLocks.push(i);
        }

        // Then shuffle that array randomly
        shuffle(allLocks);

        // Then take only the locks from index 0 to lockNum (typically 4)
        // This process removes the possiblity of random repeats when generating the locks
        var chosenLocks = allLocks.slice(keyArray.length - lockNum);

        for(i = 0; i < lockNum; i++) {
            var key = keyArray[chosenLocks[i]];
            answerKeyArray.push(key);

            // Add 1 to the index to display the correct locks
            displayLocks.push(chosenLocks[i]+1);
        }
    
        // Display the locks to the user
        document.getElementById("locks").innerHTML = displayLocks.join(" - ");

        // Reset the user input box to be empty on refresh
        var userInput = document.getElementById("user-input");
        userInput.value = "";

        // Give the user a hint by setting the placeholder for the input to be the correct answer
        userInput.placeholder = answerKeyArray.join("");
    }

    function authenticate(){

        var userInput = document.getElementById("user-input").value;
        var correctAnswer = "";

        // convert answer key array into a singular string to compare to the userInput string
        for(i = 0; i < answerKeyArray.length; i++) {
            correctAnswer += answerKeyArray[i];
        }

        if(userInput == correctAnswer){
            console.log("\n\n+++ User input, " + userInput + " is correct when compared to " + correctAnswer + " +++\n\n");
            formSuccess();
        } else {
            console.log("\n\n--- User input, " + userInput + " is INCORRECT when compared to " + correctAnswer + " ---\n\n");
            formFailure();
        }

    }

    // Taken off the Internet, thanks guys!
    function shuffle(array) {
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
    }

    window.onload = function(){ 
        // Initiate the locks and keys
        generateLocksAndKeys();
    };

});