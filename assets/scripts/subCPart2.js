console.log("Script linked properly");

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
    $('#submit-btn').val("Please wait...");

    let formData = $('#subForm').serializeArray();
    //console.log(formData);

    //sending ajax post request
    let request = $.ajax({
        async: true,
        type: "POST",
        url: "/submitC",
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

            //resetting previous value for now. These values will be added after getting verdict
            verdict.text("");
            time.text("");
            memory.text("");
            verdict.css('color', '#000');

            $('#subID').text(response.SubID);
            $('#pNumtd').text(formData[0].value);
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
        //console.log("Always")
        $('#submit-btn').prop('disabled', false);
        $('#submit-btn').val("Submit");
    });

    return false;
});

// for second/verdict modal
$(document).on('hidden.bs.modal', '.modal', function () {
    $('.modal:visible').length && $(document.body).addClass('modal-open');
});

$('#submit-solution').focus(function () {
    $('#emptyWarning').text("")
});

function addingSourceCode(srcCode) {
    //console.log(srcCode)
    $('#viewCodeC').text("");
    $('#viewCodeC').append(srcCode);

    //adding line number to the left of code segment
    let s = srcCode;
    $('#lineNumberC').text("");
    $('#lineNumberC').append(1);
    let lineNumberC = 1;
    for (i = 0; i < s.length; i++) {
        if (s[i] == '\n')
            $('#lineNumberC').append("<br>" + ++lineNumberC);
    }

    //for scrolling the line number with code segment
    (function () {
        let target = $("#lineNumberC");
        $("#viewCodeC").scroll(function () {
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