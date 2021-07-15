console.log("Script linked properly");

$(document).ready(function () {
    $('#submitAllow').click(function () {
        $('#submitAllow').prop('disabled', true);
        $('#submitAllow').text("Please wait...");

        let request = $.ajax({
            url: "/checkLogin",
            type: "GET",
            async: true,
        });
        request.done(function (response) {
            //console.log(response);
            if (response == "false") {
                window.parent.location.href = "/login";
            } else if (response == "notVerified") {
                window.parent.location.href = "/";
            } else if (response == "true") {
                window.parent.$('#submissionModal').modal('show');

                let OJ = $('#invisibleOJ').text();
                let pNum = $('#invisiblePNum').text();
                let contestID = $('#contestID').text();
                let serialID = $('#serialID').text();

                window.parent.$('#conIDSerial').val(contestID + " - " + serialID);
                window.parent.$('#conIDSerial').attr("value", contestID + " - " + serialID);
                window.parent.$('#OJform').attr("value", OJ);
                window.parent.$('#pNumform').attr("value", pNum);
                addLanguage(OJ, pNum);
                window.parent.$('#submit-solution').val("");

                //console.log(OJ, pNum, contestID, serialID);
            }
        });
        request.fail(function (response) {
            console.log(response);
        });
        request.always(function () {
            $('#submitAllow').prop('disabled', false);
            $('#submitAllow').text("Submit");
        });
    });
});

function addLanguage(OJ, pNum) {
    window.parent.$('#submit-language').empty();

    $.ajax({
        url: `/lang?OJ=` + OJ + `&pNum=` + pNum,
        method: "GET",
        success: function (data) {
            //console.log(data)

            //taking data into a map
            let map = new Map();
            $.each(data, function (key, value) {
                map.set(value, key);
            });

            //sorting map
            let mapSorted = new Map([...map.entries()].sort());

            for (let [key, value] of mapSorted) {
                // console.log(key + ' = ' + value)
                window.parent.$('#submit-language').append($("<option></option>")
                    .attr("value", value)
                    .text(key));
            }
        },
    });
};