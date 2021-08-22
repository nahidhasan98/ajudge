console.log("Script linked properly");

//adding a class to different OJ's problem for designing problem description setion
let OJ = $('#invisibleOJ').text();
let pNum = $('#invisiblePNum').text();
let username = $('#name a').text();

let subHistoryList = [];

//getting the submission list for (logged in/null if not logged) user
$(document).ready(function () {
    $.ajax({
        url: `/subHistory?OJ=` + OJ + `&pNum=` + pNum + `&user=` + username,
        type: "GET",
        async: false,
        success: function (data) {
            //console.log(data)

            subHistoryList = data;
            process();
        },
        error: function () {
            console.log('Internal Server Error. Please try again after sometime or send us a feedback.');
        }
    });
});

function process() {
    let dataCreate = "";

    if (username == "") {
        dataCreate = `<center id="notFound" style="padding:10px;">Please login to see your previous submission history.</center>`;
    } else if (subHistoryList == null) {
        dataCreate = `<center id="notFound" style="padding:10px;">No submission yet!</center>`;
    } else {
        $('#notFound').text("");

        dataCreate = `<table style="width:100%;">
                            <tr style="border-bottom:1px solid #2691d9;backgroound:#f2f2f2;height:32px;">
                                <th style="width:50%; text-align:center;">Verdict</th>
                                <th style="width:50%; text-align:center;">Submitted At</th>
                            </tr>`

        for (let i = 0; i < Math.min(8, subHistoryList.length); i++) {  //maximum 8 line will show
            //Converting Unix timestamp (like: 1549312452) to time (like: HH/MM/SS format)
            let formattedTime = timeConverter(subHistoryList[i].SubmittedAt)

            dataCreate += `<tr style="border-bottom:1px solid #bedef4;">`;

            if (subHistoryList[i].Verdict == "Accepted") { //color for accepted verdict
                dataCreate += `<td style="width:50%; text-align:center;padding: 3px 0px;color:#22aa71">` + subHistoryList[i].Verdict + `</td>`;
            } else {
                dataCreate += `<td style="width:50%; text-align:center;padding: 3px 0px;color:#de3b3b;">` + subHistoryList[i].Verdict + `</td>`;
            }

            dataCreate += `<td style="width:50%; text-align:center;">` + formattedTime + `</td></tr>`
        }
        if (subHistoryList.length > 8) {
            dataCreate += `<tr style = "border-bottom:1px solid #bedef4;">
                <td colspan="2" style="width:50%; text-align:center;padding: 3px 0px;height:40px;font-weight:bold;"><a href="/profile/`+ username + `?OJ=` + OJ + `&pNum=` + pNum + `">See More</a></td>
                        </tr>`
        }
        dataCreate += `</table>`
    }
    $('#subHistoryTable').append(dataCreate);
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

//for css work
if (OJ == "DimikOJ") {
    $('.leftContent').addClass("DimikOJ");
} else if (OJ == "Toph") {
    $('.leftContent').addClass("Toph");

    // var TophInternalLink = $("a[rel='nofollow']").attr('href');
    // firstPart = TophInternalLink.substr(0, 2);
    // if (firstPart == "/p") {
    //     probName = TophInternalLink.substr(3);
    //     $("a[rel='nofollow']").attr('href', "/problemView/Toph-" + probName);
    // }
} else if (OJ == "URI") {
    $('.leftContent').addClass("URI");
} else {
    $('.leftContent').addClass("VJudge");
}

if (OJ == "EOlymp") {
    $('i').css("display", "none");
} else if (OJ == "LightOJ") {
    $('.nav-link').css("text-decoration", "none");
    $('.dropdown-item').css("text-decoration", "none");

    $('.nav-link').css("color", "#595959");
    $('.dropdown-item').css("color", "#595959");

    $('.nav-link').hover(function () {
        $(this).css("color", "#2283c3");
    }, function () {
        $(this).css("color", "#595959");
    });

    $('.dropdown-item').hover(function () {
        $(this).css("color", "#2283c3");
    }, function () {
        $(this).css("color", "#595959");
    });
}