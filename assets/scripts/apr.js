console.log("Script linked properly");

$(document).ready(function () {
    pull();
});

function pull() {
    $('#p1Freeze').css("display", "none");
    $('#p1Loading').css("display", "block");

    let request = $.ajax({
        async: true,
        type: "POST",
        url: "/apr/pull",
    });
    request.done(function (response) {
        console.log(response);
        let resMsg = formatMessage(response.message);

        if (response.status == "error") {
            $('#p1Err').html(resMsg);
            $('#p1Loading').css("display", "none");
            $('#p1Cross').css("display", "block");

            $('#p4What').css("display", "none");
            $('#p4Cross').css("display", "block");
        } else if (resMsg.indexOf("Already") > -1) { // (if branch) Already up to date.
            $('#p1Err').html(resMsg);
            $('#p1Loading').css("display", "none");
            $('#p1Tick').css("display", "block");

            $('#p4What').css("display", "none");
            $('#p4Tick').css("display", "block");
        } else {    // successfully pulled from remote branch
            $('#p1Err').html(resMsg);
            $('#p1Loading').css("display", "none");
            $('#p1Tick').css("display", "block");

            // pull complete, calling build
            build();
        }
    });
    request.fail(function (jqXHR) {
        console.log(jqXHR);

        $('#p1Err').text("something went wrong while pulling");
        $('#p1Loading').css("display", "none");
        $('#p1Cross').css("display", "block");

        $('#p4What').css("display", "none");
        $('#p4Cross').css("display", "block");
    });
    request.always(function () {
        console.log("always")
    });
}

function formatMessage(s) {
    if (s[s.length - 1] == "\n") {  // removing last newline
        s = s.substr(0, s.length - 1);
    }
    let msg = s.replace(/\n/g, "<br>- ") // replacing internal newline
    msg = "- " + msg; // adding - at front of the message

    return msg;
}

function build() {
    $('#p2Freeze').css("display", "none");
    $('#p2Loading').css("display", "block");

    let request = $.ajax({
        async: true,
        type: "POST",
        url: "/apr/build",
    });
    request.done(function (response) {
        console.log(response)
        let resMsg = formatMessage(response.message);

        if (response.status == "error") {   // couldn't built the app
            $('#p2Err').html(resMsg);
            $('#p2Loading').css("display", "none");
            $('#p2Cross').css("display", "block");

            $('#p4What').css("display", "none");
            $('#p4Cross').css("display", "block");
        } else {    // successfully built the app
            $('#p2Err').html(resMsg);
            $('#p2Loading').css("display", "none");
            $('#p2Tick').css("display", "block");

            // build complete, calling restart
            restart();
        }
    });
    request.fail(function (jqXHR) {
        console.log(jqXHR)
    });
    request.always(function () {
        console.log("always")
    });
}

function restart() {
    $('#p3Freeze').css("display", "none");
    $('#p3Loading').css("display", "block");

    let request = $.ajax({
        async: true,
        type: "POST",
        url: "/apr/restart",
    });
    // at this point, server is restarting, so response won't come.
    // request will be failed
    request.done(function (response) {
        console.log(response)
    });
    request.fail(function (jqXHR, textStatus, errorThrown) {
        // request failed means, server is down due to restart
        // that means it is okay
        if (jqXHR.status == 502) {
            console.log("expected error occurred. nice...")
            console.log(textStatus, jqXHR.status, errorThrown)
        } else {
            console.log("unexpected error occurred")
        }

        // keep checking if server is up or not
        ping();
    });
    request.always(function () {
        console.log("always")
    });
}

function ping() {
    let counter = 0;
    let doCheck = setInterval(function () {
        let request = $.ajax({
            async: true,
            type: "GET",
            url: "/",
        });
        request.done(function (response) {
            // console.log(response)
            $('#p3Err').html("- server restarted successfully");

            $('#p3Loading').css("display", "none");
            $('#p3Tick').css("display", "block");

            $('#p4What').css("display", "none");
            $('#p4Tick').css("display", "block");

            clearInterval(doCheck);
        });
        request.fail(function (jqXHR) {
            console.log(jqXHR)
        });
        request.always(function () {
            console.log("always")
        });

        counter++;
        if (counter == 10) {
            clearInterval(doCheck);
        }
    }, 3000);  // Delay here = 3 seconds
}