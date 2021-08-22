console.log("Script linked properly")

//for reset password page
$('#password').keyup(function () {
    var pass = $('#password').val()
    var confirmPass = $('#confirmPassword').val()

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
    var pass = $('#password').val()
    var confirmPass = $('#confirmPassword').val()

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
    $('#passReset').on('submit', function () {
        if ($('#password').val() == $('#confirmPassword').val() && $('#password').val().length >= 5) {
            return true;  //if passwoed is valid then submit will be done
        }
        return false;
    });
});

//for reset common page
var h1text = $('form h1').text(), temp = "";
for (var i = 0; i < h1text.length; i++) {
    if (h1text[i] == '|') break;
    temp += h1text[i];
}
h1text = $('form h1').text(temp)

$('#email').keyup(function () {
    $('#errEmail').text("")
});

$(document).ready(function () {
    var testing = false;
    $('#resetFormEmail').on('submit', function () {
        $('.form').bind(); //prevent default submitting
        $.ajax({
            url: "/check?username=&email=" + $('#email').val().trim(),
            type: 'GET',
            async: false,
            success: function (data) {
                console.log(data)
                if (data.IsEmailExist == true) {   //email exist. proceed next
                    console.log(url)
                    testing = true;
                    $('form').attr('action');
                    $('form').unbind().submit();
                } else {  //email not found
                    $('#errEmail').text("No account found with this email. Enter correct one.")
                }
            },
            error: function () {
                console.log('Internal Server Error. Please try again after sometime or send us a feedback.');
            }
        });

        return testing;
    });
});