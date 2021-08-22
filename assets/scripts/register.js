console.log("Script linked properly")

$('#username').keyup(function () {
    $('#errUsername').text("")
});

$('#email').keyup(function () {
    $('#errEmail').text("")
});

$('#password').keyup(function () {
    var pass = $('#password').val();
    var confirmPass = $('#confirmPassword').val();

    //for length
    if (pass.length < 5) {
        $('#errPassword').text("Password length should be at least 5 characters!")
    } else {
        //for matching
        if (confirmPass.length == 0) {
            $('#errPassword').text("");
        } else {
            if (pass !== confirmPass) {
                $('#errPassword').text("Password mismatched. Put cautiously.")
            } else {
                $('#errPassword').text("")
            }
        }
    }
});
$('#confirmPassword').keyup(function () {
    var pass = $('#password').val();
    var confirmPass = $('#confirmPassword').val();

    //for matching
    if (pass !== confirmPass) {
        $('#errPassword').text("Password mismatched. Put cautiously.")
    } else {
        //for length
        if (pass.length < 5) {
            $('#errPassword').text("Password length should be at least 5 characters!")
        } else if (pass.length >= 5) {
            $('#errPassword').text("")
        }
    }
});

$(document).ready(function () {
    var testing = false;
    $('form').on('submit', function () {
        $('form').bind(); //prevent default submitting
        $.ajax({
            url: "/check?username=" + $('#username').val().trim() + "&email=" + $('#email').val().trim(),
            type: 'GET',
            async: false,
            success: function (data) {
                if (data.IsUsernameExist == true) {   //username exist. display error
                    $('#errUsername').text("Username already taken. Choose another one.")
                }
                if (data.IsEmailExist == true) {   //email exist. display error
                    $('#errEmail').text("Email already registered. Choose another one.")
                }

                if (!data.IsUsernameExist && !data.IsEmailExist) { //username & email are ok
                    //now checking for password
                    if ($('#password').val() == $('#confirmPassword').val() && $('#password').val().length >= 5) {
                        //password section is ok. now checking for captcha
                        $.ajax({
                            url: "/captcha/" + grecaptcha.getResponse(),
                            type: 'GET',
                            async: false,
                            success: function (data2) {
                                if (data2.success == true) {
                                    $("#errCaptcha").text("");

                                    testing = true;
                                    $('form').attr('action');
                                    $('form').unbind().submit();
                                } else {
                                    $("#errCaptcha").text("Captcha Error. Please fix this.");
                                }
                            },
                            error: function () {
                                console.log('Internal Server Error. Please try again after sometime or send us a feedback.');
                            }
                        });
                    }
                }
            },
            error: function () {
                console.log('Internal Server Error. Please try again after sometime or send us a feedback.');
            }
        });

        return testing;
    });
});