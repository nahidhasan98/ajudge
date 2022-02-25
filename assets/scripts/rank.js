console.log("Script linked properly")

const shiftX = 32;
let currPage = 1, totalPages;
let rankType = $('#rankType').text();
let rankList = [];

//getting the user rank list
$(document).ready(function () {
    let link = "";
    if (rankType == "OJ") {
        link = "/rank/list/oj";
    } else if (rankType == "User") {
        link = "/rank/list/user";
    }
    $.ajax({
        url: link,
        type: "GET",
        async: false,
        success: function (data) {
            rankList = data;  //assigning to a global variable

            process();
        },
        error: function () {
            console.log('Internal Server Error. Please try again after sometime or send us a feedback.');
        }
    });
});
function process() {
    //first calculating and crateing page number
    if (rankType == "User") {
        pageNumberCreate();
    }

    //onload-showing submission of page 1
    showList(1);
}
function pageNumberCreate() {
    if (rankList != null && rankList.length > 0) {
        totalPages = Math.ceil(rankList.length / 20);   //20 problem per page

        displayingPageNum = Math.min(9, totalPages)     //by default 9 pageNum displayed on pagination
        rightBound = (totalPages - displayingPageNum) * -(shiftX);

        //previous & next button display
        $('#previous').css("display", "inline-block");
        $('#next').css("display", "block");

        //page number creation
        for (let i = 1; i <= totalPages; i++) {
            $('#pageUL').append(`<li><a id="pageBox` + i + `" href="#page" onclick="showList(` + i + `)">` + i + `</a></li>`);
        }

        //width fixing of pagination
        pagerWidth = displayingPageNum * 32;
        $("#myPagination").css("width", pagerWidth + "px");
        $("#myPager").css("width", (pagerWidth + 140) + "px"); //extra 140px for 'pre/next' button
    }
}
function showList(activePage) {
    $('#loadingGif').css("display", "none");        //hide loading gif image

    if (rankList == null || rankList.length == 0) {
        //removing current existing rows
        let rowSize = $('#problemTable tr').length
        for (let i = 0; i < rowSize - 1; i++) {
            $('.problemRow').remove();
        }

        $('#notFound').text("Something went wrong. Please try again after sometime."); //if no problem found

        $('#previous').css("display", "none");      //hiding buttons
        $('#next').css("display", "none");
    } else {
        $('#notFound').text("");                    //otherwise hide this message

        currPage = activePage;                          //assigning to a global variable
        //removing current existing rows
        let rowSize = $('#problemTable tr').length
        for (let i = 0; i < rowSize - 1; i++) {
            $('.problemRow').remove();
        }

        //adding new rows
        let start = 20 * (activePage - 1);
        let finish;

        if (rankType == "OJ") {
            finish = 100;
        } else if (rankType == "User") {
            finish = start + 20;
        }

        for (let i = start; i < Math.min(finish, rankList.length); i++) {
            let serial = i + 1;

            if (rankType == "OJ") {
                tdata = rankList[i].OJ;
                tLink = "";
            } else if (rankType == "User") {
                tdata = rankList[i].FullName;
                tLink = ` [<a href="/profile/` + rankList[i].Username + `">` + rankList[i].Username + `</a>]`;
            }

            dataCreate = `<tr class="problemRow">
                        <td>`+ serial + `</td>
                        <td align="left" style="padding-left:5%;">` + tdata + tLink + `</td>
                        <td>`+ rankList[i].TotalSolved + `</td>
                        </tr>`

            $('#problemTable').append(dataCreate);
        }
        updateActiveClass(activePage);
        shifting(activePage);
    }
}
function updateActiveClass(currPage) {
    for (let i = 1; i <= totalPages; i++) {
        let id = "#pageBox" + i;

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
    showList(currPage - 1);
});
$('#nxt').click(function () {
    showList(currPage + 1);
});