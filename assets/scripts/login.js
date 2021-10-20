console.log("Login Script linked properly")

$(document).ready(function () {
    $('form').on('submit', function () {
        $('#submit').prop('disabled', true);
        $('#submit').val("Please wait...");

        if ($('#username').val().trim().length == 0) {
            $('#username').val("");
            $('#errUsername').text("username should no be empty!");
            return false;   // cancel submission
        }

        let formData = $('.loginForm').serialize();
        // console.log(formData);

        // sending ajax post request
        let request = $.ajax({
            async: true,
            type: "POST",
            url: "/login",
            data: formData,
        });
        request.done(function (response) {
            console.log(response)

            if (response.status == "error") {
                if (response.message == "username not found") {
                    $('#errUsername').text(response.message);
                } else {
                    $('#errPassword').text(response.message);
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
            $('#submit').val("Login");
        });

        return false;
    });

    $('#username').keyup(function () {
        $('#errUsername').text("")
    });

    $('#password').keyup(function () {
        $('#errPassword').text("")
    });
});