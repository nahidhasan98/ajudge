console.log("Script linked properly")

const shiftX = 32;
let currPage = 1, totalPages;
let self = $('#self').text();
let submissionList = [], tempList = [];

//taking username for this user
let username = window.location.pathname //url-path like /profile/aaa
username = username.substring(9)

//getting the submission list for (logged in/null if not logged) user
$(document).ready(function () {
    $.ajax({
        url: "/userSubmission/" + username,
        type: "GET",
        async: false,
        success: function (data) {
            $('#problemSolved').text(": " + data.ProblemSolved)
            $('#problemAttempted').text(": " + data.ProblemAttempted)

            submissionList = data.SubList;  //assigning to a global variable
            tempList = data.SubList;        //assigning to a global variable

            process();
        },
        error: function () {
            alert('Internal Server Error. Please try again after sometime or send us a feedback.');
        }
    });

    //if lastpage is provlemView & clicked on see more, then auto select OJ & pNum for see
    let extraQuery = window.location.search.substring(1);
    if (extraQuery.length > 0) {
        let newOJ = getURLParameter('OJ');
        let newPNum = getURLParameter('pNum');

        if (newOJ == "") newOJ = "All";

        $("select").val(newOJ);
        $("input[name=pNum]").val(newPNum);

        newReq();
    }
});

function process() {
    //first calculating and crateing page number
    pageNumberCreate();

    //onload-showing submission of page 1
    showSubmission(1);
}
function pageNumberCreate() {
    if (tempList != null && tempList.length > 0) {
        totalPages = Math.ceil(tempList.length / 20);   //20 problem per page

        displayingPageNum = Math.min(9, totalPages)     //by default 9 pageNum displayed on pagination
        rightBound = (totalPages - displayingPageNum) * -(shiftX);

        //previous & next button display
        $('#previous').css("display", "inline-block");
        $('#next').css("display", "block");

        //page number creation
        for (let i = 1; i <= totalPages; i++) {
            $('#pageUL').append(`<li><a id="pageBox` + i + `" href="#page" onclick="showSubmission(` + i + `)">` + i + `</a></li>`);
        }

        //width fixing of pagination
        pagerWidth = displayingPageNum * 32;
        $("#myPagination").css("width", pagerWidth + "px");
        $("#myPager").css("width", (pagerWidth + 140) + "px"); //extra 140px for 'pre/next' button
    }
}
function showSubmission(activePage) {
    $('#loadingGif').css("display", "none");        //hide loading gif image
    $('#loadingGifOJ').css("display", "none");
    $('#loadingGifPNum').css("display", "none");

    if (tempList == null || tempList.length == 0) {
        //removing current existing rows
        let rowSize = document.getElementById("submissionTable").rows.length;
        for (let i = 0; i < rowSize - 1; i++) {
            $('.submissionRow').remove();
        }

        $('#notFound').text("No Submission Found"); //if no problem found

        $('#previous').css("display", "none");      //hiding buttons
        $('#next').css("display", "none");
    } else {
        $('#notFound').text("");                    //otherwise hide this message

        currPage = activePage;                          //assigning to a global variable
        //removing current existing rows
        let rowSize = document.getElementById("submissionTable").rows.length;
        for (let i = 0; i < rowSize - 1; i++) {
            $('.submissionRow').remove();
        }

        //adding new rows
        let start = 20 * (activePage - 1);
        let finish = start + 20;

        for (let i = start; i < Math.min(finish, tempList.length); i++) {
            //Converting Unix timestamp (like: 1549312452) to time (like: HH/MM/SS format)
            let formattedTime = timeConverter(tempList[i].SubmittedAt)

            dataCreate = `<tr class="submissionRow">
                        <td>`+ tempList[i].SubID + `</td>
                        <td>`+ tempList[i].OJ + `</td>
                        <td><a href="/problemView/`+ tempList[i].OJ + "-" + tempList[i].PNum + `">` + tempList[i].PNum + `</a></td>
                        <td>`+ checkRejudge(tempList[i].TerminalVerdict, tempList[i].SubID, i);

            if (tempList[i].Verdict == "Accepted") {    //color for accepted verdict
                dataCreate += `<p id="verdict` + i + `" style="color:#1d9563">` + tempList[i].Verdict + `</p>`;
            } else {
                dataCreate += `<p id="verdict` + i + `" style="color:#de3b3b">` + tempList[i].Verdict + `</p>`;
            }

            dataCreate += `</td>
                        <td id="time`+ i + `">` + tempList[i].TimeExec + `</td>
                        <td id="memory`+ i + `">` + tempList[i].MemoryExec + `</td>
                        <td>`+ tempList[i].Language + `</td>
                        <td>`+ formattedTime + `</td>`

            if (self == "true") {
                dataCreate += `<td><button onclick="displayCode(` + i + `)">View Code</button></td>`
            }
            dataCreate += `</tr>`

            $('#submissionTable').append(dataCreate);
        }
        updateActiveClass(activePage);
        shifting(activePage);
    }
}
function updateActiveClass(currPage) {
    for (let i = 1; i <= totalPages; i++) {
        id = "#pageBox" + i;

        if ($(id).attr("class")) {
            $(id).removeClass("activePage disableClick")
        }

        if ($(id).text() == currPage) {
            $(id).addClass("activePage disableClick")
        }
    }

    //previous button
    if (currPage == 1) {
        $('#pre').addClass("disableClick")
    } else {
        $('#pre').removeClass("disableClick")
    }
    //next button
    if (currPage == totalPages) {
        $('#nxt').addClass("disableClick")
    } else {
        $('#nxt').removeClass("disableClick")
    }
}
function shifting(currPage) {
    if (currPage >= 1 && currPage <= 5) {
        $("#pageUL").css({ "transform": "translateX(" + -0 + "px)", "transition": "transform 1s ease" });
    } else if (currPage <= totalPages && currPage >= (totalPages - 4)) {
        $("#pageUL").css({ "transform": "translateX(" + rightBound + "px)", "transition": "transform 1s ease" });
    } else {
        diff = currPage - 5;
        $("#pageUL").css({ "transform": "translateX(" + -(diff * 32) + "px)", "transition": "transform 1s ease" });
    }
}
//for previous and next button
$('#pre').click(function () {
    showSubmission(currPage - 1);
});
$('#nxt').click(function () {
    showSubmission(currPage + 1);
});

function sortingFunction() {
    if (tempList != null) {
        tempList.reverse();
        showSubmission(1);
    }
}
function displayCode(index) {
    $('.popup-wrap').fadeIn(500);
    $('.popup-box').removeClass('transform-out').addClass('transform-in');

    $('#viewCode').text("").remove
    $('#viewCode').append(tempList[index].SourceCode)

    //adding line number to the left of code segment
    let s = tempList[index].SourceCode;
    $('#lineNumber').text("").remove
    $('#lineNumber').append(1)
    let lineNumber = 1
    for (i = 0; i < s.length; i++) {
        if (s[i] == '\n')
            $('#lineNumber').append("<br>" + ++lineNumber)
    }

    //for scrolling the line number with code segment
    (function () {
        let target = $("#lineNumber");
        $("#viewCode").scroll(function () {
            target.prop("scrollTop", this.scrollTop)
                .prop("scrollLeft", this.scrollLeft);
        });
    })();

    //for highlighting code syntax
    document.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightBlock(block);
    });
}
$('.popup-close').click(function () {
    $('.popup-wrap').fadeOut(500);
    $('.popup-box').removeClass('transform-in').addClass('transform-out');
});

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

//for new query
$("select[name=OJ]").change(function () {
    //displaying mini loading gif image
    $('#loadingGifOJ').css("display", "inline-block");

    //request for new search
    newReq();
})
$("input[name=pNum]").change(function () {
    $('#loadingGifPNum').css("display", "inline-block");
    newReq();
})
function newReq() {
    let OJ = $("select").val();
    let pNum = $("input[name=pNum]").val().trim();

    //selecting submission according to requirement
    tempList = [];

    if (OJ != "All" && pNum != "") {
        for (i = 0; i < submissionList.length; i++) {
            pNum = pNum.toLowerCase();
            listPNum = submissionList[i].PNum.toLowerCase();

            if (OJ == submissionList[i].OJ && listPNum.indexOf(pNum) != -1) {
                tempList.push(submissionList[i]);
            }
        }
    } else if (OJ != "All" && pNum == "") {
        for (i = 0; i < submissionList.length; i++) {
            if (OJ == submissionList[i].OJ) {
                tempList.push(submissionList[i]);
            }
        }
    } else if (OJ == "All" && pNum != "") {
        for (i = 0; i < submissionList.length; i++) {
            pNum = pNum.toLowerCase();
            listPNum = submissionList[i].PNum.toLowerCase();

            if (listPNum.indexOf(pNum) != -1) {
                tempList.push(submissionList[i]);
            }
        }
    } else if (OJ == "All" && pNum == "") {
        tempList = submissionList;
    }
    process();
}

function getURLParameter(sParam) {
    let sPageURL = window.location.search.substring(1);
    let sURLVariables = sPageURL.split('&');

    for (let i = 0; i < sURLVariables.length; i++) {
        let sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] == sParam) {
            return sParameterName[1];
        }
    }
}