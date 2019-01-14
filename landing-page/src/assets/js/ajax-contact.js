$("#dynauth-alpha").submit(function(event){
    // cancels the form submission
    event.preventDefault();
    authenticate();
    if($("#email").val()){
        submitForm();
    } else {
        // Silence?
    }
});

function submitForm(){
// Initiate Variables With Form Content
var email = $("#email").val();

  $.ajax({
        type: "POST",
        url: "https://dynauth.io/dynauth.php",
        data: "email=" + email,
        success : function(text){
            if (text == "success"){
                formSuccess();
            }
        }
    });
}

function formSuccess(){
    $( "#failure-alert" ).hide();
    $( "#success-alert" ).slideDown();
}

function formFailure(){
    $( "#success-alert" ).hide();
    $( "#failure-alert" ).slideDown();
    generateLocksAndKeys();
}