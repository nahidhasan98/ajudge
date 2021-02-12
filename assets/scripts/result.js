var verdict = $('#verdict')
var time = $('#time')
var memory = $('#memory')
var subID = $('#subID');

$(document).ready(function () {
    if (subID.text() != "") {
        $('#loadingVerdictGif').css("display", "inline-flex");     //displaying loading gif image

        var url = "/verdict/subID=" + subID.text();
        //console.log(url)

        var counter = 0;
        var doCheck = setInterval(function () {
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
        }, 3000);  //Delay here = 3 seconds
    }
});

$('#rejudge').click(function () {
    $('#rejudge').css("display", "none");           //hide rejudge image
    $('#judging').css("display", "inline-flex");    //displaying judging gif image

    var url = "/rejudge/subID=" + subID.text();

    var counter = 0;
    var doCheck = setInterval(function () {
        $.getJSON(url, function (result) {
            counter++;
            verdict.text(result.Status);

            if (result.TerminalVerdict == true) {           //got final verdict - don't recall the verdict
                $('#judging').css("display", "none");       //hide judging gif image
                time.text(result.Runtime);
                memory.text(result.Memory);

                clearInterval(doCheck);
            } else if (counter >= 20) {
                $('#judging').css("display", "none");           //displaying judging gif image
                $('#rejudge').css("display", "inline-flex");    //displaying rejudge image
                //location.reload();
                clearInterval(doCheck);
            }
        });
    }, 3000);  //Delay here = 3 seconds
});