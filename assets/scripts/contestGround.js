console.log("Script linked correctly.")

let problemSetLength = parseInt($('#problemSetLength').text().trim());
let timerCurrentTime
let timerBeforeStart

$(document).ready(function () {
    //v-pills ul
    $('#A-tab').addClass('active');
    $('#A-tab').attr('aria-selected', 'true');

    //v-pills div
    $('#problemsA').addClass('show active');

    if ($('#iframePA').text().trim() == "") {
        let dataCreate = `<iframe id="myIframeA" scrolling="no" frameborder="0" src="/problemSet/` + $('#contestID').text() + `/A/` + $('#problemsA #OJ').text() + `-` + $('#problemsA #pNum').text() + `" style="width: 100%;"></iframe>
                        <script src="../assets/scripts/iframeResizer.min.js"></script>
                        <script>iFrameResize({ log: false }, '#myIframeA')</script>`
        $('#iframePA').append(dataCreate)
    }

    //setting up dashboard tab
    for (let i = 65; i < (65 + problemSetLength); i++) {
        //console.log(i, String.fromCharCode(i));
        let serialIndex = String.fromCharCode(i);

        //for link design
        $('#goto' + serialIndex).hover(function () {
            $(this).css({ 'text-decoration': 'underline', 'cursor': 'pointer' });
        }, function () {
            $(this).css({ 'text-decoration': 'none', 'cursor': 'default' });
        });

        //for tab action (if link is clicked in dashboard tab)
        $('#goto' + serialIndex).click(function () {
            let OJ = $('#problems' + serialIndex + ' #OJ').text();
            let pNum = $('#problems' + serialIndex + ' #pNum').text();

            if ($('#iframeP' + serialIndex).text().trim() == "") {
                let dataCreate = `<iframe id="myIframe` + serialIndex + `" scrolling="no" frameborder="0" src="/problemSet/` + $('#contestID').text() + `/` + serialIndex + `/` + OJ + `-` + pNum + `" style="width: 100%;"></iframe>
                        <script src="../assets/scripts/iframeResizer.min.js"></script>
                        <script>iFrameResize({ log: false }, '#myIframe` + serialIndex + `')</script>`
                $('#iframeP' + serialIndex).append(dataCreate)
            }

            gotoTab(serialIndex);
        });

        //for retrieving problem description (if problem number is clicked)
        $('#' + serialIndex + '-tab').click(function () {
            let OJ = $('#problems' + serialIndex + ' #OJ').text();
            let pNum = $('#problems' + serialIndex + ' #pNum').text();

            if ($('#iframeP' + serialIndex).text().trim() == "") {
                let dataCreate = `<iframe id="myIframe` + serialIndex + `" scrolling="no" frameborder="0" src="/problemSet/` + $('#contestID').text() + `/` + serialIndex + `/` + OJ + `-` + pNum + `" style="width: 100%;"></iframe>
                        <script src="../assets/scripts/iframeResizer.min.js"></script>
                        <script>iFrameResize({ log: false }, '#myIframe` + serialIndex + `')</script>`
                $('#iframeP' + serialIndex).append(dataCreate)
            }
        });
    }

    let callData;
    getContestData();   //calling one time on page load

    //fixing the timer
    let contestStartAtStr = $('#contestStartAt').text().trim();
    let contestStartAt = parseInt(contestStartAtStr) * 1000;
    let now = new Date().getTime()

    let contestDurationStr = $('#contestDuration').text().trim();
    let cdHH = contestDurationStr.substr(0, 2);
    let cdMM = contestDurationStr.substr(3, 2);
    let cdHour = parseInt(cdHH)
    let cdMin = parseInt(cdMM)

    let contestDuration = (cdHour * 60 * 60 * 1000) + (cdMin * 60 * 1000)

    if (now < contestStartAt) {
        showCurrentTime();   //for avoiding first 1 sec delay
        showBeforeStart();   //for avoiding first 1 sec delay
        timerCurrentTime = setInterval(showCurrentTime, 1000);
        timerBeforeStart = setInterval(showBeforeStart, 1000);
    } else if (now >= contestStartAt && now < (contestStartAt + contestDuration)) {
        showRunningTimer();
        (function updateContestData() {
            getContestData();
            callData = setTimeout(updateContestData, 15000); //every 15 seconds this function will be calling
        })();
    } else {
        clearTimeout(callData);
    }
});

function gotoTab(tempSerial) {
    //pills ul
    $('#dashboard-tab').removeClass('active');
    $('#submissions-tab').removeClass('active');
    $('#standings-tab').removeClass('active');

    $('#dashboard-tab').attr('aria-selected', 'false');
    $('#submissions-tab').attr('aria-selected', 'false');
    $('#standings-tab').attr('aria-selected', 'false');

    $('#problems-tab').addClass('active');
    $('#problems-tab').attr('aria-selected', 'true');

    //pills div
    $('#dashboard').removeClass('show active');
    $('#submissions').removeClass('show active');
    $('#standings').removeClass('show active');

    $('#problems').addClass('show active');

    for (let i = 65; i < (65 + problemSetLength); i++) {
        //console.log(i, String.fromCharCode(i));
        let serialIndex = String.fromCharCode(i);

        if (serialIndex == tempSerial) {
            //v-pills ul
            $('#' + tempSerial + '-tab').addClass('active');
            $('#' + tempSerial + '-tab').attr('aria-selected', 'true');

            //v-pills div
            $('#problems' + tempSerial).addClass('show active');
        } else {
            //v-pills ul
            $('#' + serialIndex + '-tab').removeClass('active');
            $('#' + serialIndex + '-tab').attr('aria-selected', 'false');

            //v-pills div
            $('#problems' + serialIndex).removeClass('show active');
        }
    }
}

let cSubmissionList = [];
let cSolvedStatus = [];
function getContestData() {
    //console.log("Hello");
    $.ajax({
        url: "/dataContest/" + $('#contestID').text(),
        type: "GET",
        async: false,
        success: function (data) {
            //console.log(data);
            cSubmissionList = data.CSubmissionList;
            displaySubmissionList();
            displaySolvedStatus(data.CSolvedStatus, data.CAttempedStatus, data.CTotalSolved, data.CTotalSubmission)
            displayStandings(data.CContestantData);
        },
        error: function () {
            alert('Internal Server Error. Please try again after sometime or send us a feedback.');
        }
    });
}
function displaySubmissionList() {
    $('#loadingGif').css("display", "none");        //hide loading gif image

    if (cSubmissionList == null || cSubmissionList.length == 0) {
        //removing current existing rows
        let rowSize = $('#submissionTable tr').length;
        for (let i = 0; i < rowSize - 1; i++) {
            $('.submissionRow').remove();
        }
        $('#notFound').text("No Submission Found"); //if no problem foun
    } else {
        $('#notFound').text("");                    //otherwise hide this message

        //removing current existing rows
        let rowSize = $('#submissionTable tr').length;
        for (let i = 0; i < rowSize - 1; i++) {
            $('.submissionRow').remove();
        }

        //adding new rows
        for (let i = 0; i < cSubmissionList.length; i++) {
            //Converting Unix timestamp (like: 1549312452) to time (like: HH/MM/SS format)
            let formattedTime = timeConverter(cSubmissionList[i].SubmittedAt)

            let dataCreate = `<tr class="submissionRow">
                <td>`+ cSubmissionList[i].SubID + `</td>
                <td><a href="/profile/`+ cSubmissionList[i].Username + `">` + cSubmissionList[i].Username + `</a></td>
                <td>`+ cSubmissionList[i].SerialIndex + `</td>
                <td>`+ checkRejudge(cSubmissionList[i].TerminalVerdict, cSubmissionList[i].SubID, i);

            if (cSubmissionList[i].Verdict == "Accepted") {    //color for accepted verdict
                dataCreate += `<p id="verdict` + i + `" style="color:#1d9563">` + cSubmissionList[i].Verdict + `</p>`;
            } else {
                dataCreate += `<p id="verdict` + i + `" style="color:#de3b3b">` + cSubmissionList[i].Verdict + `</p>`;
            }

            dataCreate += `</td>
            <td id="time`+ i + `">` + cSubmissionList[i].TimeExec + `</td>
            <td id="memory`+ i + `">` + cSubmissionList[i].MemoryExec + `</td>
            <td>`+ cSubmissionList[i].Language + `</td>
            <td>`+ formattedTime + `</td>`

            if ($('#contestUser').text().trim() == cSubmissionList[i].Username || $('#contestUser').text().trim() == $('#contestAuthor').text().trim()) {
                dataCreate += `<td><button onclick="displayCode(` + i + `)" data-toggle="modal" data-target="#modal">View Code</button></td>`
            } else {
                dataCreate += `<td><button disabled style="opacity: 50%;color: black;cursor: default;">View Code</button></td>`
            }
            dataCreate += `</tr>`

            $('#submissionTable').append(dataCreate);
        }
    }
}
function displayCode(index) {
    $('#viewCode').text("");
    $('#viewCode').append(cSubmissionList[index].SourceCode);

    //adding line number to the left of code segment
    let s = cSubmissionList[index].SourceCode;
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
//Converting Unix timestamp (like: 1549312452) to time (like: DATE MONTH YEAR HH/MM/SS AM format)
function timeConverter(UNIX_timestamp) {
    let a = new Date(UNIX_timestamp * 1000);
    let months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    let year = a.getFullYear();
    let month = months[a.getMonth()];
    let date = a.getDate();
    let timeDate = date + '-' + month + '-' + year;

    let time12format = new Date(UNIX_timestamp * 1000).toLocaleTimeString("en-US");
    let time = timeDate + ' (' + time12format + ')';

    return time;
}
function checkRejudge(terminalVerdict, subID, i) {
    let res = "";

    if (!terminalVerdict) {
        res = `<img src="../assets/images/rejudge.png" id="rejudge` + i + `" style="width: 20px; display:inline-flex; cursor:pointer;" title="Click to rejudge" onclick="doRejudge(` + subID + `,` + i + `)">
        <img src="../assets/images/loadingMini.gif" id="judging`+ i + `" style="width: 20px; display:none; cursor:pointer;" title="Judging...">`;
    }

    return res;
}
function doRejudge(subID, i) {
    $('#rejudge' + i).css("display", "none");           //hide rejudge image
    $('#judging' + i).css("display", "inline-flex");    //displaying judging gif image

    let url = "/rejudge/subID=" + subID;

    let counter = 0;
    let doCheck = setInterval(function () {
        $.getJSON(url, function (result) {
            counter++;
            $('#verdict' + i).text(result.Status);

            if (result.TerminalVerdict == true) {           //got final verdict - don't recall the verdict
                $('#judging' + i).css("display", "none");       //hide judging gif image
                $('#time' + i).text(result.Runtime);
                $('#memory' + i).text(result.Memory);

                clearInterval(doCheck);
            } else if (counter >= 20) {
                $('#judging' + i).css("display", "none");           //displaying judging gif image
                $('#rejudge' + i).css("display", "inline-flex");    //displaying rejudge image

                clearInterval(doCheck);
            }
        });
    }, 3000);  //Delay here = 3 seconds
};
function displaySolvedStatus(solvedStatus, attempedStatus, totalSolved, totalSubmission) {
    //setting up dashboard tab
    for (let i = 65; i < (65 + problemSetLength); i++) {
        //for dashboard solved status
        let serialIndex = String.fromCharCode(i);

        if (solvedStatus[serialIndex] == true) {
            $('#myStatus' + serialIndex).text("Accepted");
            $('#myStatus' + serialIndex).css('color', '#1d9563');
        } else if (attempedStatus[serialIndex] == true) {
            $('#myStatus' + serialIndex).text("Tried");
            $('#myStatus' + serialIndex).css('color', '#e68a00');
        }

        //for overall solved/submission
        let solveCounter = 0, submissionCounter = 0;
        let res = "";

        if (totalSubmission[serialIndex] > 0) {
            if (totalSolved[serialIndex] > 0) {
                res = totalSolved[serialIndex] + "/" + totalSubmission[serialIndex];
            } else {
                res = "0/" + totalSubmission[serialIndex];
            }
        }
        $('#totalStatus' + serialIndex).text(res);
    }
}
function displayStandings(contestantData) {
    //console.log(contestantData);
    $('#loadingGifStanding').css("display", "none");        //hide loading gif image

    if (contestantData == null || contestantData.length == 0) {
        //removing current existing rows
        let rowSize = $('#standingTable tr').length;
        for (let i = 0; i < rowSize - 1; i++) {
            $('.standingRow').remove();
        }
        $('#notFound').text("No Submissions Yet"); //if no problem foun
    } else {
        $('#notFound').text("");                    //otherwise hide this message

        //removing current existing rows
        let rowSize = $('#standingTable tr').length;
        for (let i = 0; i < rowSize - 1; i++) {
            $('.standingRow').remove();
        }

        //calculating col width
        let totalCol = (1 + 4 + 2 + problemSetLength);  //rank=x1,user=x4,score=x2,perProb=x1
        let x1 = 100 / totalCol;
        let x2 = 2 * x1;
        let x4 = 100 - ((x1 * (1 + problemSetLength)) + (x2 * 1));

        //adding new rows
        for (let i = 0; i < contestantData.length; i++) {
            let rank = i + 1;

            let dataCreate = `<tr class="standingRow">
                <td style="width:`+ x1 + `%;">` + rank + `</td>
                <td style="width:`+ x4 + `%;"><a href="/profile/` + contestantData[i].Username + `">` + contestantData[i].Username + `</a></td>
                <td style="width:`+ x2 + `%;" title="Solved: ` + contestantData[i].TotalSolved + ` / Penalty Time: ` + makeMinute(contestantData[i].TotalTime) + `"><span style="color:#1d9563;font-weight: bold;font-size: 15px;">` + contestantData[i].TotalSolved + `</span><br>` + makeMinute(contestantData[i].TotalTime) + `</td>`;

            if (contestantData[i].SubDetails != null) { //at least 1 submission done(except compilation error) by this user
                for (let j = 65; j < (65 + problemSetLength); j++) { //for A to Z problem
                    let serialIndex = String.fromCharCode(j);
                    if (contestantData[i].SubDetails[serialIndex] != undefined) {
                        //console.log(contestantData[i].SubDetails[serialIndex]);
                        let tempVerdict = contestantData[i].SubDetails[serialIndex].Verdict;
                        let penaltyCount = contestantData[i].SubDetails[serialIndex].Penalty;
                        let comErrorCount = contestantData[i].SubDetails[serialIndex].CompilationError;
                        let tdTitle = " / Penalty: " + penaltyCount + " / Compilation Error: " + comErrorCount;

                        if (tempVerdict == "Accepted") {
                            dataCreate += `<td style="width:` + x1 + `%;" title="Accepted: 1` + tdTitle + `">` + checkFirstSolved(contestantData, i, serialIndex) + `<br><span style="color:#1d9563;">1</span> / <span style="color:#de3b3b;">` + penaltyCount + `</span> / <span style="color:#e68a00;">` + comErrorCount + `</span></td>`;
                        } else {
                            dataCreate += `<td style="width:` + x1 + `%; "title="Accepted: 0` + tdTitle + `"><img src="../assets/images/cross.png" style="width:15px;"><br><span style="color:#1d9563;">0</span> / <span style="color:#de3b3b;">` + penaltyCount + `</span> / <span style="color:#e68a00;">` + comErrorCount + `</span></td>`;
                        }
                    } else {
                        dataCreate += `<td style="width:` + x1 + `%;"></td>`;
                    }
                }
            }

            $('#standingTable').append(dataCreate);
        }
    }
}
function makeMinute(timeSeconds) {
    let timeMinute = Math.floor(timeSeconds / 60);
    return timeMinute;
}
//current timing section
function showCurrentTime() {
    let time = new Date();
    let hour = time.getHours();
    let min = time.getMinutes();
    let sec = time.getSeconds();
    am_pm = "AM";

    if (hour > 12) {
        hour -= 12;
        am_pm = "PM";
    }
    if (hour == 0) {
        hr = 12;
        am_pm = "AM";
    }

    if (hour < 10) hour = "0" + hour;
    if (min < 10) min = "0" + min;
    if (sec < 10) sec = "0" + sec;

    let currentTime = hour + ":" + min + ":" + sec + " " + am_pm;

    $('#currentTime').text(currentTime);

    let a = new Date();
    let months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    let year = a.getFullYear();
    let month = months[a.getMonth()];
    let date = a.getDate();

    if (date < 10) date = "0" + date;
    let today = date + '-' + month + '-' + year;

    $('#currentDate').text(today);

    let contestStartAtStr = $('#contestStartAt').text().trim(); //this is in seconds
    let contestStartAt = parseInt(contestStartAtStr) * 1000;
    let timeNow = new Date().getTime();

    if (timeNow >= contestStartAt) {
        $('#cDurationLabel').text("Contest started. Refresh the page.")
        $('#cDurationLabel').css('color', 'Green');
        clearInterval(timerCurrentTime);
    }
}

function showBeforeStart() {
    // Timer to be stopped at
    let contestStartAtStr = $('#contestStartAt').text().trim();
    let contestStartAt = parseInt(contestStartAtStr) * 1000;

    let now = new Date().getTime()
    let distance = contestStartAt - now;

    //let days = Math.floor(distance / (1000 * 60 * 60 * 24));
    //let hours = Math.floor((distance / (1000 * 60 * 60)) % 24);
    let hours = Math.floor(distance / (1000 * 60 * 60));
    let minutes = Math.floor((distance / (1000 * 60)) % 60);
    let seconds = Math.floor((distance / 1000) % 60);

    if (hours < 10) hours = "0" + hours;
    if (minutes < 10) minutes = "0" + minutes;
    if (seconds < 10) seconds = "0" + seconds;

    $('#beforeStart').text(hours + ":" + minutes + ":" + seconds);

    if (distance <= 1000) {
        clearInterval(timerBeforeStart);
    }
}

function showRunningTimer() {
    //console.log($('#currentTimeLabel').text())
    let timer = setInterval(function () {
        let contestStartAtStr = $('#contestStartAt').text().trim();
        let contestStartAt = parseInt(contestStartAtStr) * 1000;
        let timeNow = new Date().getTime();

        let elapsed = timeNow - contestStartAt;
        //console.log(elapsed);

        //let days = Math.floor(distance / (1000 * 60 * 60 * 24));
        //let hours = Math.floor((distance / (1000 * 60 * 60)) % 24);
        let hours = Math.floor(elapsed / (1000 * 60 * 60));
        let minutes = Math.floor((elapsed / (1000 * 60)) % 60);
        let seconds = Math.floor((elapsed / 1000) % 60);

        if (hours < 10) hours = "0" + hours;
        if (minutes < 10) minutes = "0" + minutes;
        if (seconds < 10) seconds = "0" + seconds;

        $('#currentTime').text(hours + ":" + minutes + ":" + seconds);

        let contestDurationStr = $('#contestDuration').text().trim();
        let cdHH = contestDurationStr.substr(0, 2);
        let cdMM = contestDurationStr.substr(3, 2);
        let cdHour = parseInt(cdHH)
        let cdMin = parseInt(cdMM)

        let contestDuration = (cdHour * 60 * 60 * 1000) + (cdMin * 60 * 1000)
        let totalTime = contestStartAt + contestDuration;

        let timeLeft = totalTime - timeNow

        let hoursLeft = Math.floor(timeLeft / (1000 * 60 * 60));
        let minutesLeft = Math.floor((timeLeft / (1000 * 60)) % 60);
        let secondsLeft = Math.floor((timeLeft / 1000) % 60);

        if (hoursLeft == -1) hoursLeft = "0";   //to avoiding 0-1:0-1:0-1 format
        if (minutesLeft == -1) minutesLeft = "0";
        if (secondsLeft == -1) secondsLeft = "0";

        if (hoursLeft < 10) hoursLeft = "0" + hoursLeft;
        if (minutesLeft < 10) minutesLeft = "0" + minutesLeft;
        if (secondsLeft < 10) secondsLeft = "0" + secondsLeft;

        $('#beforeStart').text(hoursLeft + ":" + minutesLeft + ":" + secondsLeft);

        if (totalTime <= timeNow) {
            $('#currentTimeLabel').text("Full Time");
            $('#beforeStartLabel').text("Contest Finished");

            clearInterval(timer);
        }
    }, 1000);
}

function checkFirstSolved(conData, contestantSerial, serial) {
    let res = "";
    let temp = Number.MAX_SAFE_INTEGER;

    for (let i = 0; i < conData.length; i++) {
        if (conData[i].SubDetails[serial] != undefined) {
            temp = Math.min(temp, conData[i].SubDetails[serial].AcceptedAt);
        }
    }

    //if current user submission(accepted) is first
    if (conData[contestantSerial].SubDetails[serial].AcceptedAt == temp) {
        res = `<img src="../assets/images/tickCircle.png" style="width:16px;">`; //first AC
    } else {
        res = `<img src="../assets/images/tick.png" style="width:15px;">`;  //normal AC
    }

    return res;
}

//for highlighter js to detect language each time
$('.modal-close-icon').click(function () {
    $('#viewCode').removeAttr('class');
});
