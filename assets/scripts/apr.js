console.log("Script linked properly");

$(document).ready(function () {
    $('#p1Freeze').css("display", "none");
    $('#p1Loading').css("display", "block");

    let request = $.ajax({
        async: true,
        type: "POST",
        url: "/apr/pull",
    });
    request.done(function (response) {
        let s = response.message;
        if (s[s.length - 1] == "\n") {
            s = s.substr(0, s.length - 1);
        }
        let resMsg = s.replace(/\n/g, "<br>- ")
        resMsg = "- " + resMsg;

        console.log(response);

        if (response.status == "error") {
            $('#p1Err').html(resMsg);
            $('#p1Loading').css("display", "none");
            $('#p1Cross').css("display", "block");

            $('#p4What').css("display", "none");
            $('#p4Cross').css("display", "block");
        } else if (resMsg.indexOf("Already") > -1) {
            $('#p1Err').html(resMsg);
            $('#p1Loading').css("display", "none");
            $('#p1Tick').css("display", "block");

            $('#p4What').css("display", "none");
            $('#p4Tick').css("display", "block");
        } else {
            $('#p1Err').html(resMsg);
            $('#p1Loading').css("display", "none");
            $('#p1Tick').css("display", "block");

            $('#p2Freeze').css("display", "none");
            $('#p2Loading').css("display", "block");

            // calling restart
            restart();

            // checking if server restarted or not
            ping();
        }
    });
    request.fail(function () {
        $('#p1Err').text("something went wrong while pulling");
        $('#p1Loading').css("display", "none");
        $('#p1Cross').css("display", "block");

        $('#p4What').css("display", "none");
        $('#p4Cross').css("display", "block");
    });
    request.always(function () {
        console.log("always")
    });
});

function restart() {
    $('#p2Loading').css("display", "none");
    $('#p2Tick').css("display", "block");

    $('#p3Freeze').css("display", "none");
    $('#p3Loading').css("display", "block");

    let request = $.ajax({
        async: true,
        type: "POST",
        url: "/apr/restart",
    });
    request.done(function (response) {
        console.log(response)
    });
    request.fail(function (response) {
        console.log(response)
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
            type: "POST",
            url: "/apr/pull",
        });
        request.done(function (response) {
            console.log(response)

            $('#p3Loading').css("display", "none");
            $('#p3Tick').css("display", "block");

            $('#p4What').css("display", "none");
            $('#p4Tick').css("display", "block");

            clearInterval(doCheck);
        });
        request.fail(function (response) {
            console.log(response)
        });
        request.always(function () {
            console.log("always")
        });

        counter++;
        if (counter == 10) {
            clearInterval(doCheck);
        }
    }, 3000);  //Delay here = 3 seconds
}