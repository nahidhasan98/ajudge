console.log("Script linked properly")

$(document).ready(function () {
    $('form').on('submit', function () {
        let res = validateForm();
        if (!res) {
            return false
        }

        $('#submit').prop('disabled', true);
        $('#submit').val("Signing up...");

        let formData = $('form').serialize();
        // console.log(formData);

        // sending ajax post request
        let request = $.ajax({
            async: true,
            type: "POST",
            url: "/register",
            data: formData,
        });
        request.done(function (response) {
            // console.log(response)

            if (response.status == "error") {
                for (i = 0; i < response.errors.length; i++) {
                    if (response.errors[i].type == "fullName") {
                        $('#errFullName').text(response.errors[i].message);
                    } else if (response.errors[i].type == "email") {
                        $('#errEmail').text(response.errors[i].message);
                    } else if (response.errors[i].type == "username") {
                        $('#errUsername').text(response.errors[i].message);
                    } else if (response.errors[i].type == "password") {
                        $('#errPassword').text(response.errors[i].message);
                    } else { //captcha error or any other errors
                        $('#errCaptcha').text(response.errors[i].message);
                    }
                }
            } else {
                window.location.replace(response.redirectURL);
            }
        });
        request.fail(function (response) {
            console.log(response)
        });
        request.always(function () {
            $('#submit').prop('disabled', false);
            $('#submit').val("Sign Up");
        });

        return false;
    });

    $('#fullName').keyup(function () {
        $('#errFullName').text("")
    });

    $('#username').keyup(function () {
        $('#errUsername').text("")
    });

    $('#email').keyup(function () {
        $('#errEmail').text("")
    });

    $('#password').keyup(function () {
        // for length
        if ($('#password').val().length < 5) {
            $('#errPassword').text("password length should be at least 5 characters")
        } else {
            // for matching
            if ($('#confirmPassword').val().length == 0) {
                $('#errPassword').text("");
            } else {
                if ($('#password').val() !== $('#confirmPassword').val()) {
                    $('#errPassword').text("password mismatched, put cautiously")
                } else {
                    $('#errPassword').text("")
                }
            }
        }
    });

    $('#confirmPassword').keyup(function () {
        // for matching
        if ($('#password').val() !== $('#confirmPassword').val()) {
            $('#errPassword').text("password mismatched, put cautiously")
        } else {
            // for length
            if ($('#password').val().length < 5) {
                $('#errPassword').text("password length should be at least 5 characters")
            } else if ($('#password').val().length >= 5) {
                $('#errPassword').text("")
            }
        }
    });
});

function validateForm() {
    // taking care of fullName
    if ($('#fullName').val().trim().length == 0) {
        $('#fullName').val("");
        $('#errFullName').text("name should no be empty");
        return false;   // cancel submission
    }

    // taking care of username
    // PART 1: length
    if ($('#username').val().trim().length == 0) {
        $('#username').val("");
        $('#errUsername').text("username should no be empty");
        return false;   // cancel submission
    }

    // PART 2: space
    let index = $('#username').val().trim().indexOf(" ");
    if (index > -1) {   // contains spaces
        $('#errUsername').text("username can't contains space");
        return false;   // cancel submission
    }

    // taking care of password
    // PART 1: length
    if ($('#password').val().length < 5) {
        $('#errPassword').text("password length should be at least 5 characters")
        return false;   // cancel submission
    }

    // PART 2: matching
    if ($('#password').val() !== $('#confirmPassword').val()) {
        $('#errPassword').text("password mismatched, put cautiously")
        return false;   // cancel submission
    }

    // taking care of recaptcha
    if ($('[name="g-recaptcha-response"]').val() == "") {
        $('#errCaptcha').text("captcha error, please verify you are not a robot")
        return false;   // cancel submission
    } else {
        $('#errCaptcha').text("")
    }

    return true
}