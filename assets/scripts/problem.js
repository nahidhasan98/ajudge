console.log("Script linked properly");

//declaring variables
let shiftX = 32, previousPage = 5, currPage = 1;
var OJ, pNum, pName;
let submissionList = [];

//checking cookies for getting recent OJ, pNum, pName searched
if (document.cookie != "") {
    var cookieArray = document.cookie.split(';');
    var cookieIndex = -1;

    //checking for desired cookie segment
    for (var k = 0; k < cookieArray.length; k++) {
        var temp = cookieArray[k].trim().substring(0, 3);

        if (temp == "OJ=") {
            cookieIndex = k;
            break;
        }
    }
    if (cookieIndex != -1) {
        var cookieText = cookieArray[cookieIndex].split('&');
        var OJText = cookieText[0].split('='), pNumText = cookieText[1].split('='), pNameText = cookieText[2].split('=');

        OJ = OJText[1];
        pNum = pNumText[1];
        pName = pNameText[1];
    }
}

if (OJ == undefined) OJ = "All";    //if no cookies found
if (pNum == undefined) pNum = "";
if (pName == undefined) pName = "";

//taking username for this user
var username = $('#name').text().trim();
getUserSubmission(username);

//getting the submission list for (logged in/null if not logged) user
function getUserSubmission(username) {
    $.ajax({
        url: "/userSubmission/" + username,
        method: "GET",
        success: function (data) {
            submissionList = data.SubList;
        },
    });
};

let ajaxCall;
getProblemList(OJ, pNum, pName, ""); //on page load
function getProblemList(OJ, pNum, pName, tag) {
    ajaxCall = $.ajax({
        url: `/problemList?OJ=` + OJ + `&pNum=` + pNum + `&pName=` + pName,
        method: 'GET',
        success: function (response) {
            //console.log(response);

            //hiding mini loading gif
            $('#loadingGifOJ').css("display", "none");
            $('#loadingGifPNum').css("display", "none");
            $('#loadingGifPName').css("display", "none");

            $('#loadingGifEasy').css("display", "none");
            $('#loadingGifMath').css("display", "none");
            $('#loadingGifGeometry').css("display", "none");
            $('#loadingGifDataStructure').css("display", "none");
            $('#loadingGifImplementation').css("display", "none");
            $('#loadingGifHard').css("display", "none");
            $('#loadingGifString').css("display", "none");
            $('#loadingGifGreedy').css("display", "none");
            $('#loadingGifSorting').css("display", "none");
            $('#loadingGifBFS').css("display", "none");
            $('#loadingGifDP').css("display", "none");
            $('#loadingGifProbability').css("display", "none");
            $('#loadingGifNumberTheory').css("display", "none");
            $('#loadingGifGraph').css("display", "none");

            if (tag == "") {
                //setting input value (like: placeholder)
                $("select[name=OJ]").val(OJ);
                $("input[name=pNum]").val(pNum);
                $("input[name=pName]").val(pName);
            } else {
                $("input[name=pName]").val("");
            }

            //calling the process function to display & other stuffs
            if (response == null) { //if no problem found
                //console.log("Null value returned")
                processPList([]);
            } else {
                processPList(response);
            }
        },
    });
};

function processPList(getPList) {
    //console.log(getPList)
    pList = getPList    //taking to a variable for using in other function

    //calculating some stuffs for pagination
    totalPages = Math.ceil(pList.length / 20);  //20 problem per page
    displayingPageNum = Math.min(9, totalPages) //by default 9 pageNum displayed on pagination
    rightBound = (totalPages - displayingPageNum) * -(shiftX);

    //previous & next button display
    $('#previous').css("display", "inline-block");
    $('#next').css("display", "block");

    //page number creation
    for (i = 1; i <= totalPages; i++) {
        $('#pageUL').append(`<li><a id="pageBox` + i + `" href="#" onclick="displayProblem(` + i + `)">` + i + `</a></li>`);
    }
    //width fixing of pagination
    pagerWidth = displayingPageNum * 32;
    $("#myPagination").css("width", pagerWidth + "px");
    $("#myPager").css("width", (pagerWidth + 140) + "px"); //extra 140px for 'pre/next' button

    //displaying problem of page 1
    displayProblem(1);
}

function displayProblem(active) {
    $('#loadingGif').css("display", "none");    //hide loading gif image

    if (pList.length == 0) {
        $('#notFound').text("No Problem Found");    //if no problem found

        $('#previous').css("display", "none");  //hiding buttons
        $('#next').css("display", "none");
    } else {
        $('#notFound').text("");    //otherwise hide this message
    }

    currPage = active;  //taking to a variable for using in other function

    //removing current existing rows
    var rowSize = $('#problemTable tr').length;
    for (i = 0; i < rowSize - 2; i++) {
        $('.problemRow').remove();
    }

    //adding new rows
    var start = 20 * (active - 1);
    var finish = start + 20;

    for (i = start; i < Math.min(finish, pList.length); i++) {
        var link = "/problemView/" + pList[i].OJ + "-" + pList[i].PNum;

        dataCreate = `<tr class="problemRow">`;
        if (username != null && username != "") {
            //setting a color to the verdict result
            let solvedStatus = getSolvedStatus(pList[i].OJ, pList[i].PNum);
            if (solvedStatus == "Solved") {
                dataCreate += `<td style="color:#1d9563">Solved</td>`;
            } else if (solvedStatus == "Tried") {
                dataCreate += `<td style="color:#e68a00">Tried</td>`;
            } else {
                dataCreate += `<td></td>`;
            }
        }
        dataCreate += `<td>` + pList[i].OJ + `</td>
                    <td>`+ pList[i].PNum + `</td>
                    <td><a href="`+ link + `">` + pList[i].PName + `</a></td>
                </tr>`;

        $('#problemTable').append(dataCreate);
    }
    updateActiveClass(active);
    shiftingPageNum(active);
}
function getSolvedStatus(OJ, pNum) {
    var result = "";

    if (submissionList != null) {
        for (var k = 0; k < submissionList.length; k++) {
            if (submissionList[k].OJ == OJ && submissionList[k].PNum == pNum) {
                result = "Tried";
                if (submissionList[k].Verdict == "Accepted") {
                    result = "Solved";
                    break;
                }
            }
        }
    }
    return result;
}

function updateActiveClass(currPage) {
    for (i = 1; i <= totalPages; i++) {
        id = "#pageBox" + i;

        if ($(id).attr("class")) {
            $(id).removeClass("activePage disableClick");
        }

        if ($(id).text() == currPage) {
            $(id).addClass("activePage disableClick");
        }
    }

    //previous button
    if (currPage == 1) {
        $('#pre').addClass("disableClick");
    } else {
        $('#pre').removeClass("disableClick");
    }
    //next button
    if (currPage == totalPages) {
        $('#nxt').addClass("disableClick");
    } else {
        $('#nxt').removeClass("disableClick");
    }
}
function shiftingPageNum(currPage) {
    if (currPage == 1 || currPage == 2 || currPage == 3 || currPage == 4 || currPage == 5) {
        $("#pageUL").css({ "transform": "translateX(" + -0 + "px)", "transition": "transform 1s ease" });
    } else if (currPage == totalPages || currPage == (totalPages - 1) || currPage == (totalPages - 2) || currPage == (totalPages - 3) || currPage == (totalPages - 4)) {
        $("#pageUL").css({ "transform": "translateX(" + rightBound + "px)", "transition": "transform 1s ease" });
    } else {
        diff = currPage - 5;
        $("#pageUL").css({ "transform": "translateX(" + -(diff * 32) + "px)", "transition": "transform 1s ease" });
    }
    previousPage = currPage;
}

//for previous and next button
$('#pre').click(function () {
    displayProblem(currPage - 1);
});
$('#nxt').click(function () {
    displayProblem(currPage + 1);
});

//for new query
$("select[name=OJ]").change(function () {
    ajaxCall.abort();

    //displaying mini loading gif image
    $('#loadingGifOJ').css("display", "inline-block");

    //request for new search
    newReq();
})
$("input[name=pNum]").change(function () {
    ajaxCall.abort();

    $('#loadingGifPNum').css("display", "inline-block");
    newReq();
})
$("input[name=pName]").change(function () {
    ajaxCall.abort();

    $('#loadingGifPName').css("display", "inline-block");
    newReq();
})

function newReq(tag = "") {
    if (tag == "") {
        OJ = $("select").val();

        pNum = $("input[name=pNum]").val().trim();
        pName = $("input[name=pName]").val().trim();

        //setting cookies
        setCookie(OJ, pNum, pName, 1);
    } else {
        OJ = "All";
        pNum = "";
        pName = tag;

        //no cookies set for tag
    }

    //gettting new problem list by searching
    getProblemList(OJ, pNum, pName, tag);
}
// $(document).ready(function () {
// });

function setCookie(cOJ, cPNum, cPName, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
    var expires = "expires=" + d.toGMTString();
    document.cookie = "OJ=" + cOJ + "&pNum=" + cPNum + "&pName=" + cPName + ";" + expires + ";path=/";
}
//for tag button click
$('#tagEasy').click(function () {
    ajaxCall.abort();

    $('#loadingGifEasy').css("display", "inline-block");
    newReq("easy");
});
$('#tagMath').click(function () {
    ajaxCall.abort();

    $('#loadingGifMath').css("display", "inline-block");
    newReq("math");
});
$('#tagGeometry').click(function () {
    ajaxCall.abort();

    $('#loadingGifGeometry').css("display", "inline-block");
    newReq("geometry");
});
$('#tagDataStructure').click(function () {
    ajaxCall.abort();

    $('#loadingGifDataStructure').css("display", "inline-block");
    newReq("data+structure");
});
$('#tagImplementation').click(function () {
    ajaxCall.abort();

    $('#loadingGifImplementation').css("display", "inline-block");
    newReq("implementation");
});
$('#tagHard').click(function () {
    ajaxCall.abort();

    $('#loadingGifHard').css("display", "inline-block");
    newReq("hard");
});
$('#tagString').click(function () {
    ajaxCall.abort();

    $('#loadingGifString').css("display", "inline-block");
    newReq("string");
});
$('#tagGreedy').click(function () {
    ajaxCall.abort();

    $('#loadingGifGreedy').css("display", "inline-block");
    newReq("greedy");
});
$('#tagSorting').click(function () {
    ajaxCall.abort();

    $('#loadingGifSorting').css("display", "inline-block");
    newReq("sorting");
});
$('#tagBFS').click(function () {
    ajaxCall.abort();

    $('#loadingGifBFS').css("display", "inline-block");
    newReq("bfs");
});
$('#tagDP').click(function () {
    ajaxCall.abort();

    $('#loadingGifDP').css("display", "inline-block");
    newReq("dp");
});
$('#tagProbability').click(function () {
    ajaxCall.abort();

    $('#loadingGifProbability').css("display", "inline-block");
    newReq("probability");
});
$('#tagNumberTheory').click(function () {
    ajaxCall.abort();

    $('#loadingGifNumberTheory').css("display", "inline-block");
    newReq("number+theory");
});
$('#tagGraph').click(function () {
    ajaxCall.abort();

    $('#loadingGifGraph').css("display", "inline-block");
    newReq("graph");
});