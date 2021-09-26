$(document).ready(function () {
    let submit = false;
    $('form').on('submit', function () {
        $('#send-btn').prop('disabled', true);
        $('#send-btn').val("Please wait...");

        let request = $.ajax({
            url: "/captcha/" + grecaptcha.getResponse(),
            type: 'GET',
            async: false,
        });
        request.done(function (response) {
            if (response.success == true) {
                $("#errCaptcha").text("");
                submit = true;
            } else {
                $("#errCaptcha").text("Captcha Error. Please fix this.");
                $('#send-btn').prop('disabled', false);
                $('#send-btn').val("Submit");
            }
        });
        request.fail(function (response) {
            console.log(response)
        });
        request.always(function () {

        });

        return submit;
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