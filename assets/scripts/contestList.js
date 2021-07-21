console.log("Script linked properly")

const shiftX = 32;
let currPage = 1, totalPages;
let contestList = [];

//getting the user rank list
$(document).ready(function () {
    $.ajax({
        url: "/listContest",
        type: "GET",
        async: false,
        success: function (data) {
            //console.log(data)
            contestList = data;  //assigning to a global variable

            process();
            $('#combinedRankBtn').css("display","block");
        },
        error: function (resp) {
            console.log(resp);
        }
    });
});
function process() {
    //first calculating and crateing page number
    pageNumberCreate();

    //onload-showing submission of page 1
    showList(1);
}
function pageNumberCreate() {
    if (contestList != null && contestList.length > 0) {
        totalPages = Math.ceil(contestList.length / 20);   //20 problem per page

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

    //removing current existing rows
    let rowSize = $('#problemTable tr').length
    for (let i = 0; i < rowSize - 1; i++) {
        $('.problemRow').remove();
    }

    if (contestList == null || contestList.length == 0) {
        $('#notFound').text("Something went wrong. Currently contest list is empty. Please try again after sometime."); //if no problem found

        $('#previous').css("display", "none");          //hiding buttons
        $('#next').css("display", "none");
    } else {
        $('#notFound').text("");                        //otherwise hide this message

        currPage = activePage;                          //assigning to a global variable

        //adding new rows
        let start = 20 * (activePage - 1);
        let finish = start + 20;

        for (let i = start; i < Math.min(finish, contestList.length); i++) {
            let link = ``;

            dataCreate = `<tr class="problemRow">
                        <td>`+ contestList[i].ContestID + `</td>
                        <td align="left" style="padding-left:3%;"><a href="/contest/`+ contestList[i].ContestID + `">` + contestList[i].Title + `</a></td>
                        <td>`+ startTime(contestList[i].StartAt) + `</td>
                        <td>`+ startDate(contestList[i].StartAt) + `</td>
                        <td>`+ contestList[i].Duration + ` h</td>
                        <td><a href="/profile/`+ contestList[i].Author + `">` + contestList[i].Author + `</a></td>
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

function startTime(startAt) {
    let time12format = new Date(startAt * 1000).toLocaleTimeString("en-US");

    return time12format;
}
function startDate(startAt) {
    let a = new Date(startAt * 1000);
    let months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    let year = a.getFullYear();
    let month = months[a.getMonth()];
    let date = a.getDate();
    let timeDate = date + '-' + month + '-' + year;

    return timeDate;
}

$('#combinedRankBtn').click(function () {
    let btnText = $(this).text();

    if (btnText == "Combined Rank List"){
        $(this).text("Show Rank");
        $(this).prop('disabled', true);
        $(this).css('cursor', 'not-allowed');
    
        $('#problemTable tr').each(function(i, row){
            if ( i == 0 ) {
                let data = $(row).html();
                $(row).empty();
                data = `<th>Select</th>` + data;
                $(row).append(data);
            } else {
                let data = $(row).html();
                let cID = $(row)[0]["cells"][0]["childNodes"][0]["textContent"];
                //console.log(id);
                insertingData = `<td><input type="checkbox" class="chkBox" onclick="checkBoxFn(` + cID + `, this)" style="width: 20px;height: 20px;"></td>`
                $(row).empty();
                data = insertingData + data;
                $(row).append(data);
            }
            $('.chkBox').val($(this).is(':checked'));
        });
    } else if (btnText =="Show Rank"){
        //console.log("HHH");
        $('#cRankModal').modal('show');
    
        $('#loadingGifStanding').css("display", "block");        //hide loading gif image
        //removing current existing rows
        let rowSize = $('#standingTable tr').length;
        for (let i = 0; i < rowSize - 1; i++) {
            $('.standingRow').remove();
        }

        //removing cIDs column
        $('.cIDAdded').remove();

        getCombinedData();
    }
});

let cList = new Set();
function checkBoxFn(cID, item){
    //console.log(cID, $(item));
    let isChecked = $(item).is(':checked');

    if (isChecked == true) {
        cList.add(cID);
    } else {
        cList.delete(cID);
    }
    //console.log(cList.size);
    //show/hide send button
    if (parseInt(cList.size) > 0) {
        $('#combinedRankBtn').prop('disabled', false);
        $('#combinedRankBtn').css('cursor', 'pointer');
    } else {
        $('#combinedRankBtn').prop('disabled', true);
        $('#combinedRankBtn').css('cursor', 'not-allowed');
    }
    //console.log(cList)
}

function getCombinedData(){
    $('#combinedRankBtn').prop('disabled', true);
    $('#combinedRankBtn').text("Please wait...");
    //console.log(cList);
    //sending ajax post request
    let request = $.ajax({
        async: true,
        type: "POST",
        url: "/getCombinedStandings",
        data: {ids: JSON.stringify(Array.from(cList))},
    });

    request.done(function (response) {
        //console.log(response)
        displayStandings(response);
    });

    request.fail(function (response) {
        console.log(response)
    });

    request.always(function () {
        //console.log("Always")
        $('#combinedRankBtn').prop('disabled', false);
        $('#combinedRankBtn').text("Show Rank");
    });
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
        $('#notFoundStanding').text("No Submissions Found"); //if no problem found
    } else {
        $('#notFoundStanding').text("");                    //otherwise hide this message

        //removing current existing rows
        let rowSize = $('#standingTable tr').length;
        for (let i = 0; i < rowSize - 1; i++) {
            $('.standingRow').remove();
        }

        //removing cIDs column
        $('.cIDAdded').remove();
        //creating cIDs column
        cListSorted = [...(new Set(cList))].sort()
        for (let item of cListSorted) {
            //console.log(item);
            $('#standingTable').find('tr:first').append(`<th class="cIDAdded">Contest ID `+item+`</th>`);
        }

        //calculating col width
        let totalCol = (1 + 4 + 2 + cListSorted.length);  //rank=x1,user=x4,score=x2,perProb=x1
        let x1 = 100 / totalCol;
        let x2 = 2 * x1;
        let x3 = 3 * x1;
        let x4 = 100 - ((x1 * (1 + cListSorted.length)) + (x2 * 1));

        //adding new rows
        for (let i = 0; i < contestantData.length; i++) {
            let rank = i + 1;

            let dataCreate = `<tr class="standingRow">
                <td style="width:`+ x1 + `%;">` + rank + `</td>
                <td style="width:`+ x3 + `%;"><a href="/profile/` + contestantData[i].Username + `">` + contestantData[i].Username + `</a></td>
                <td style="width:`+ x2 + `%;background: aliceblue;" title="Solved: ` + contestantData[i].TotalSolved + ` / Penalty Time: ` + makeMinute(contestantData[i].TotalTime) + `"><span style="color:#1d9563;font-weight: bold;font-size: 15px;">` + contestantData[i].TotalSolved + `</span><br>` + makeMinute(contestantData[i].TotalTime) + `</td>`;

            let idx = 0;
            for (let item of cListSorted) {
                if (idx >= contestantData[i].PerContestStatus.length){
                    dataCreate += `<td></td>`;
                    continue;
                }
                if (contestantData[i].PerContestStatus[idx].ConID == item) {
                    dataCreate += `<td style="width:`+ x1 + `%;" title="Solved: ` + contestantData[i].PerContestStatus[idx].PerSolved + ` / Penalty Time: ` + makeMinute(contestantData[i].PerContestStatus[idx].PerTime) + `"><span style="color:#1d9563;font-weight: bold;font-size: 15px;">` + contestantData[i].PerContestStatus[idx].PerSolved + `</span><br>` + makeMinute(contestantData[i].PerContestStatus[idx].PerTime) + `</td>`;
                    idx++;
                } else {
                    dataCreate += `<td></td>`;
                }
            }
            dataCreate += `</tr>`;
            $('#standingTable').append(dataCreate);
        }
    }
}
function makeMinute(timeSeconds) {
    let timeMinute = Math.floor(timeSeconds / 60);
    return timeMinute;
}