console.log("Script linked properly");

addLanguage($('#OJ').val(), $('#pNum').val());

function addLanguage(OJ, pNum) {
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
                $('#submit-language').append($("<option></option>")
                    .attr("value", value)
                    .text(key));
            }
        },
    });
};

$('#submit-solution').focus(function () {
    $('#emptyWarning').text("")
});

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
                window.location.href = "/login";
            } else if (response == "notVerified") {
                window.location.href = "/";
            } else if (response == "true") {
                $('#submissionModal').modal('show');
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

function addingSourceCode(srcCode) {
    $('#viewCode').text("");
    $('#viewCode').append(srcCode);

    //adding line number to the left of code segment
    let s = srcCode;
    $('#lineNumber').text("");
    $('#lineNumber').append(1);
    let lineNumber = 1;
    for (i = 0; i < s.length; i++) {
        if (s[i] == '\n')
            $('#lineNumber').append("<br>" + ++lineNumber);
    }

    //for scrolling the line number with code segment
    (function () {
        let target = $("#lineNumber");
        $("#viewCode").scroll(function () {
            target.prop("scrollTop", this.scrollTop)
                .prop("scrollLeft", this.scrollLeft);
        });
    })();

    // first, find all the div.code blocks
    document.querySelectorAll('pre code').forEach(el => {
        // then highlight each
        hljs.highlightElement(el);
    });
}

$('form').on('submit', function () {
    let sourceCode = $('#submit-solution');
    if (sourceCode.val().trim().length == 0) {
        sourceCode.val("");
        $('#emptyWarning').text("Source code should no be empty!");
        return false;   //cancel submission
    }
    let OJ = $('#OJform').val();
    if (OJ != "DimikOJ" && OJ != "Toph" && OJ != "URI" && sourceCode.val().trim().length < 50) {
        $('#emptyWarning').text("Source code is too short. Minimum 50 characters required to submit your solution!");
        return false;
    }

    $('#submit-btn').prop('disabled', true);
    $('#submit-btn').val("Submitting...");

    let formData = $('#subForm').serialize();
    //console.log(formData);
    //sending ajax post request
    let request = $.ajax({
        async: true,
        type: "POST",
        url: "/submit",
        data: formData,
    });
    request.done(function (response) {
        //console.log(response)

        if (response.error != "") {// for vj, got some error, submission not done
            $('#emptyWarning').text(response.error);
        } else {
            //now hide submission modal
            $('#submissionModal').modal('hide');
            $('#verdictModal').modal('show');

            $('#subID').text(response.SubID);
            $('#OJtd').text(response.OJ);
            $('#pNumLink').text(response.PNum);
            $('#pNumLink').attr("href", "/problemView/" + response.OJ + "-" + response.PNum);
            // time, memory & verdict will be added from result.js
            $('#language').text(response.Language);
            $('#submittedAt').text(response.SubmittedAt);
            addingSourceCode(response.SourceCode);

            getVerdict();
        }
    });
    request.fail(function (response) {
        console.log(response)
    });
    request.always(function () {
        $('#submit-btn').prop('disabled', false);
        $('#submit-btn').val("Submit");
    });
    return false;
});

// for second/verdict modal
$(document).on('hidden.bs.modal', '.modal', function () {
    $('.modal:visible').length && $(document.body).addClass('modal-open');
});


let verdict = $('#verdict')
let time = $('#time')
let memory = $('#memory')
let subID = $('#subID');

$('#rejudge').click(function () {
    $('#rejudge').css("display", "none");           //hide rejudge image
    $('#judging').css("display", "inline-flex");    //displaying judging gif image

    let url = "/rejudge/subID=" + subID.text();

    let counter = 0;
    let doCheck = setInterval(function () {
        $.getJSON(url, function (result) {
            counter++;
            verdict.text(result.Status);

            if (result.TerminalVerdict == true) {           //got final verdict - don't recall the verdict
                $('#judging').css("display", "none");       //hide judging gif image
                time.text(result.Runtime);
                memory.text(result.Memory);

                //setting a color to the verdict result
                if (result.Status == "Accepted") {
                    verdict.css('color', '#1d9563');
                } else {
                    verdict.css('color', '#de3b3b');
                }

                clearInterval(doCheck);
            } else if (counter >= 20) {
                $('#judging').css("display", "none");           //displaying judging gif image
                $('#rejudge').css("display", "inline-flex");    //displaying rejudge image
                //location.reload();
                clearInterval(doCheck);
            }
        });
    }, 5000);  //Delay here = 5 seconds
});

function getVerdict() {
    if (subID.text() != "") {
        $('#loadingVerdictGif').css("display", "inline-flex");     //displaying loading gif image

        let url = "/verdict/subID=" + subID.text();
        //console.log(url)

        let counter = 0;
        let doCheck = setInterval(function () {
            $.getJSON(url, function (result) {
                counter++;
                $('#loadingVerdictGif').css("display", "none");    //hide loading gif image
                $('#judging').css("display", "inline-flex");       //displaying judging gif image
                verdict.text(result.Status);

                if (result.TerminalVerdict == true) {                       //got final verdict - don't recall the verdict
                    $('#judging').css("display", "none");           //hide judging gif image
                    time.text(result.Runtime);
                    memory.text(result.Memory);

                    //setting a color to the verdict result
                    if (result.Status == "Accepted") {
                        verdict.css('color', '#1d9563');
                    } else {
                        verdict.css('color', '#de3b3b');
                    }

                    clearInterval(doCheck);
                } else if (counter >= 20) {
                    $('#judging').css("display", "none");           //displaying judging gif image
                    $('#rejudge').css("display", "inline-flex");    //displaying rejudge image
                    //location.reload();
                    clearInterval(doCheck);
                }
            });
        }, 5000);  //Delay here = 5 seconds
    }
}
