var currentURL = window.location.href;

var xmlHttp = new XMLHttpRequest();
xmlHttp.open("GET", "https://localhost:443/password?url=" + currentURL, true); // true for asynchronous 
xmlHttp.onreadystatechange = function() { 
    if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
        // callback(xmlHttp.responseText);
        var jsonResponse = JSON.parse(xmlHttp.responseText);
        document.querySelectorAll('input[type=text]')[0].value = jsonResponse.email;
        document.querySelectorAll('input[type=password]')[0].value = jsonResponse.password;
        document.querySelectorAll('form')[0].submit();
    }
}
xmlHttp.send();


// var formAction = document.querySelectorAll('form')[0].action;
// console.log(formAction);

// var formData = new FormData();
// formData.append("email", jsonResponse.email);
// formData.append("password", jsonResponse.password);

// var request = new XMLHttpRequest();
// request.open("POST", formAction);
// request.send(formData);