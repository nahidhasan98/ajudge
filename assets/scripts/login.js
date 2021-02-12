console.log("Login Script linked properly")

var username = $('#username');
var password = $('#password');
var errUsername = $('#errUsername');
var errPassword = $('#errPassword');

username.keyup(function () {
    errUsername.text("")
});
password.keyup(function () {
    errPassword.text("")
});

$(document).ready(function () {
    var testing = false;
    $('form').on('submit', function () {
        $('form').bind(); //prevent default submitting
        $.ajax({
            url: "/check?username=" + username.val().trim(),
            type: 'GET',
            async: false,
            success: function (data) {
                if (data.IsUsernameExist == true) {   //username exist. go for login
                    testing = true;
                    $('form').attr('action');
                    $('form').unbind().submit();
                } else {  //username not found
                    errUsername.text("Username not found.")
                }
            },
            error: function () {
                alert('Internal Server Error. Please try again after sometime or send us a feedback.');
            }
        });

        return testing;
    });
});