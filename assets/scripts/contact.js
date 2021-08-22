$(document).ready(function () {
    var testing = false;
    $('form').on('submit', function () {
        $('form').bind(); //prevent default submitting
        $.ajax({
            url: "/captcha/" + grecaptcha.getResponse(),
            type: 'GET',
            async: false,
            success: function (data) {
                if (data.success == true) {
                    testing = true;
                    $("#errCaptcha").text("");
                } else {
                    $("#errCaptcha").text("Captcha Error. Please fix this.");
                }
            },
            error: function () {
                console.log('Internal Server Error. Please try again after sometime or send us a feedback.');
            }
        });

        return testing;
    });

    $("#mailName").click(function () {
        $("#errCaptcha").text("");
    });
    $("#mailEmail").click(function () {
        $("#errCaptcha").text("");
    });
    $("#mailMessage").click(function () {
        $("#errCaptcha").text("");
    });
});