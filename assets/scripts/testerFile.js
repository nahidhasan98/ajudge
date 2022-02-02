console.log("Script linked correctly.")

$(document).ready(function () {
    restart();
});

function restart() {
    $('#p3Freeze').css("display", "none");
    $('#p3Loading').css("display", "block");

    let request = $.ajax({
        async: true,
        type: "POST",
        url: "/test2",
    });
    // at this point, server is restarting, so response won't come.
    // request will be failed
    request.done(function (response) {
        console.log(response)
    });
    request.fail(function (response, status) {
        // request failed means, 
        console.log(response.status, status)
    });
    request.always(function () {
        console.log("always")
    });
}