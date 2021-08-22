console.log("Script linked properly")

let OJ = "All", pNum = "", pName = "";
let problemList = [];

//getting 20 random problem list
getProblemList(OJ, pNum, pName);

function getProblemList(OJ, pNum, pName) {
    $.ajax({
        url: `/problemList?OJ=` + OJ + `&pNum=` + pNum + `&pName=` + pName,
        type: "GET",
        async: false,
        success: function (data) {
            problemList = data;  //assigning to a global variable
            showProblemList();
        },
        error: function () {
            console.log('Internal Server Error. Please try again after sometime or send us a feedback.');
        }
    });
};

function showProblemList() {
    $('#loadingGif').css("display", "none");        //hide loading gif image
    $('#loadingGifOJ').css("display", "none");
    $('#loadingGifPNum').css("display", "none");
    $('#loadingGifPName').css("display", "none");

    //removing current existing rows
    let rowSize = $('#problemTable tr').length;
    for (let i = 0; i < rowSize - 1; i++) {
        $('.problemRow' + i).remove();
    }

    if (problemList == null || problemList.length == 0) {
        $('#notFound').text("No Problem Found"); //if no problem found
    } else {
        $('#notFound').text("");                    //otherwise hide this message

        for (let i = 0; i < Math.min(20, problemList.length); i++) {
            let link = "/problemView/" + problemList[i].OJ + "-" + problemList[i].PNum;

            dataCreate = `<tr class="problemRow` + i + `">
                            <td id="OJ">` + problemList[i].OJ + `</td>
                            <td id="pNum">`+ problemList[i].PNum + `</td>
                            <td id="pName"><a href="`+ link + `" target="_blank">` + problemList[i].PName + `</a></td>
                            <td><button onclick="doSelect(` + i + `)">Select</button></td>
                        </tr>`;

            $('#problemTable').append(dataCreate);
        }
    }
}

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
$("input[name=pName]").change(function () {
    //console.log("Hello1")
    $('#loadingGifPName').css('display', 'inline-block');
    //console.log("Hello2")
    newReq();
})
function newReq() {
    OJ = $("select").val();
    pNum = $("input[name=pNum]").val().trim();
    pName = $("input[name=pName]").val().trim();

    getProblemList(OJ, pNum, pName);
}

let classIndex = 0, serialIndex = 65;

function doSelect(index) {
    //removing warning section
    $('div[role="alert"]').removeClass('d-block');
    $('div[role="alert"]').addClass('d-none');

    let selectedOJ = $('.problemRow' + index + ` #OJ`).text();
    let selectedPNum = $('.problemRow' + index + ` #pNum`).text();
    let selectedPName = $('.problemRow' + index + ` #pName`).text();
    let flag = 0;

    let rowSize = $('#contestProbSelectedTable tr').length;

    if (rowSize <= 26) {
        for (let i = 1; i < rowSize; i++) {
            let tempOJ = $("#contestProbSelectedTable tr:eq(" + i + ") #OJ").text();
            let tempPNum = $("#contestProbSelectedTable tr:eq(" + i + ") #pNum").text();

            if (tempOJ == selectedOJ && tempPNum == selectedPNum) {
                alert("This problem already selected!");
                flag = 1;
                break;
            }
        }

        if (flag == 0) {
            let dataCreate = `<tr class="selectedRow` + classIndex + `">
                            <td scope="col" id="serialNum">` + String.fromCharCode(serialIndex) + `</td>
                            <td scope="col" id="OJ"><textarea type="text" name="OJ`+ serialIndex + `" value="" readonly class="cTableInput">` + selectedOJ + `</textarea></td>
                            <td scope="col" id="pNum"><textarea type="text" name="pNum`+ serialIndex + `" value="" readonly class="cTableInput">` + selectedPNum + `</textarea></td>
                            <td scope="col" id="pName"><textarea type="text" name="pName`+ serialIndex + `" value="" readonly class="cTableInput">` + selectedPName + `</textarea></td>
                            <td scope="col" id="customName"><textarea type="text" name="customName`+ serialIndex + `" placeholder="Give a custom name" class="contestCustomInput" style="resize: none;overflow:hidden;"></textarea></td>
                            <td scope="col" ><button onclick="removeSelected(` + classIndex + `)">Remove</button></td>
                        </tr>`;

            $('#contestProbSelectedTable').append(dataCreate);

            classIndex++;
            serialIndex++;
        }
    } else {
        alert("Maximum 26 problem can be added!");
    }
}
function removeSelected(index) {
    // let selectedOJ = $('.selectedRow' + index + ` #OJ`).text();
    // let selectedPNum = $('.selectedRow' + index + ` #pNum`).text();
    // let selectedPName = $('.selectedRow' + index + ` #pName`).text();

    let nextSerial = $('.selectedRow' + index + ' #serialNum').text();  //current serialIndex will be the serialIndex cause current one is being removed
    let nextSerialInt = nextSerial.charCodeAt(0);

    let rowSize = $('#contestProbSelectedTable tr').length;
    let currRowIndex = $('.selectedRow' + index).index();

    //removing current row
    $('.selectedRow' + index).remove();;
    serialIndex--;

    //renaming serial Number
    for (let i = currRowIndex; i < rowSize; i++) {
        $("#contestProbSelectedTable tr:eq(" + i + ") #serialNum").text(nextSerial);
        $("#contestProbSelectedTable tr:eq(" + i + ") #OJ textarea").attr('name', 'OJ' + nextSerialInt);
        $("#contestProbSelectedTable tr:eq(" + i + ") #pNum textarea").attr('name', 'pNum' + nextSerialInt);
        $("#contestProbSelectedTable tr:eq(" + i + ") #pName textarea").attr('name', 'pName' + nextSerialInt);
        $("#contestProbSelectedTable tr:eq(" + i + ") #customName textarea").attr('name', 'customName' + nextSerialInt);

        nextSerialInt++;
        nextSerial = String.fromCharCode(nextSerialInt);
    }
}

//setting a default date (today's date)
let a = new Date();
let year = a.getFullYear();
let month = a.getMonth() + 1;   //months are from 0 to 11, so added 1
let day = a.getDate();

if (month < 10) {
    month = '0' + month;
}
if (day < 10) {
    day = '0' + day;
}

//if it is for update page, then date should be set
let pathname = window.location.pathname;
if (pathname.indexOf("/contestUpdate/") == 0) {
    let cDate = $('#cDate').text().trim();
    $('input[name="contestDate"]').attr('value', cDate);

    //initial set up index variable because already some problem is added if it is update page
    classIndex = 0 + parseInt($('#cProbSetLength').text().trim()), serialIndex = 65 + parseInt($('#cProbSetLength').text().trim());
} else {
    let today = year + '-' + month + '-' + day; //input[type=date] takes this format
    $('input[name="contestDate"]').attr('min', today);
    $('input[name="contestDate"]').attr('value', today);
}

$(document).ready(function () {
    //setting up client side time zone offset. js gives -360 as minute for GMT +06:00
    let offset = new Date().getTimezoneOffset();
    $('input[name="timeZoneOffset"]').val(offset * 60);

    //on form submit
    $('form').on('submit', function () {
        $('form').bind(); //prevent default submitting

        // console.log($('#serialNum').text());
        // console.log($('input[name="contestTitle"]').val());
        // console.log($('input[name="contestDate"]').val());
        // console.log($('input[name="contestTime"]').val());
        // console.log($('input[name="contestDuration"]').val());

        //checking wheather problem set empty or not
        if ($('#serialNum').text() == "") {
            $('div[role="alert"]').text("Problem set empty! Please select at least one problem.");
            $('div[role="alert"]').addClass('d-block');

            return false;
        }
        //checking wheather title empty or not
        if ($('input[name="contestTitle"]').val().trim() == "") {
            //setting up empty if there present spaces
            $('input[name="contestTitle"]').val("");
            $('input[name="contestTitle"]').addClass('alert alert-danger');

            $('div[role="alert"]').text("Contest title cannot be empty!");
            $('div[role="alert"]').addClass('d-block');

            return false;
        }

        //checking wheather start time valid or not
        if (pathname.indexOf("/contestUpdate/") != 0) { //skip time validation for update
            let flag = 1;
            //checking 24-hour input format is valid
            let tInput = $('input[name="contestTime"]').val().trim();

            let h1 = parseInt(tInput.substr(0, 1), 10);
            let h2 = parseInt(tInput.substr(1, 1), 10);
            let hh = parseInt(tInput.substr(0, 2), 10);

            let m1 = parseInt(tInput.substr(3, 1), 10);
            let m2 = parseInt(tInput.substr(4, 1), 10);
            let mm = parseInt(tInput.substr(3, 2), 10);

            let s1 = parseInt(tInput.substr(6, 1), 10);
            let s2 = parseInt(tInput.substr(7, 1), 10);
            let ss = parseInt(tInput.substr(6, 2), 10);
            //console.log(tInput, h1, h2, hh, m1, m2, mm, s1, s2, ss)

            if (isNaN(h1) || h1 < 0 || h1 > 2 || isNaN(h2) || h2 < 0 || h2 > 9 || isNaN(h2) || hh < 0 || hh > 23) flag = 0;
            if (isNaN(m1) || m1 < 0 || m1 > 5 || isNaN(m2) || m2 < 0 || m2 > 9 || isNaN(mm) || mm < 0 || mm > 59) flag = 0;
            if (isNaN(s1) || s1 < 0 || s1 > 5 || isNaN(s2) || s2 < 0 || s2 > 9 || isNaN(ss) || ss < 0 || ss > 59) flag = 0;

            //console.log(flag)
            if (flag == 0) {
                //setting up empty if there present spaces
                //$('input[name="contestTime"]').val("");
                $('input[name="contestTime"]').addClass('alert alert-danger');

                $('div[role="alert"]').text("Start time should be in 24-hour format!");
                $('div[role="alert"]').addClass('d-block');

                return false;
            } else { //24-hour format ok
                tInput = "";
                if (hh < 10) tInput = "0";
                tInput += hh.toString() + ":";
                if (mm < 10) tInput += "0";
                tInput += mm.toString() + ":";
                if (ss < 10) tInput += "0";
                tInput += ss.toString()

                $('input[name="contestTime"]').val(tInput);
            }

            let startTime = new Date($('input[name="contestDate"]').val() + " " + $('input[name="contestTime"]').val().trim());
            let currentTime = new Date();
            // console.log(startTime);
            // console.log(currentTime);
            if (startTime <= currentTime) {
                //setting up empty if there present spaces
                //$('input[name="contestTime"]').val("");
                $('input[name="contestTime"]').addClass('alert alert-danger');

                $('div[role="alert"]').text("Start time must be later from now!");
                $('div[role="alert"]').addClass('d-block');

                return false;
            }
        }

        //checking wheather duration valid or not
        let duration = $('input[name="contestDuration"]').val().trim();
        let len = duration.length;
        let flag = 1;

        if (len >= 5 && len <= 6) {   //like 05:00 or 120:00
            if (duration[len - 1] >= '0' && duration[len - 1] <= '9' && duration[len - 2] >= '0' && duration[len - 2] <= '5' && duration[len - 3] == ':') {
                for (i = len - 4; i >= 0; i--) {
                    if (duration[i] < '0' || duration[i] > '9') {
                        flag = 0;   //invalid
                        break;
                    }
                }

                if (flag) { //still valid? then do another check
                    //checking all value/character is zero or not like: 00:00
                    let flag2 = 0;  //suppose invalid
                    for (i = 0; i < len; i++) {
                        if (duration[i] != ':' && duration[i] != '0') {
                            flag2 = 1;  //valid
                            break;
                        }
                    }
                    if (flag2 == 0) {   //if invalid
                        flag = 0;   //invalid
                    }
                }
            } else {
                flag = 0;
            }
        } else {
            flag = 0;
        }

        if (flag) {
            return true;
        } else {
            //setting up empty if there present spaces
            $('input[name="contestDuration"]').val("");
            $('input[name="contestDuration"]').addClass('alert alert-danger');

            $('div[role="alert"]').text("Give the duration in hh:mm format!");
            $('div[role="alert"]').addClass('d-block');

            return false;
        }
    });
});

$('input[name="contestTitle"]').focus(function () {
    $('input[name="contestTitle"]').removeClass('alert alert-danger');

    $('div[role="alert"]').removeClass('d-block');
    $('div[role="alert"]').addClass('d-none');
});
$('input[name="contestDate"]').focus(function () {
    $('input[name="contestDate"]').removeClass('alert alert-danger');

    $('div[role="alert"]').removeClass('d-block');
    $('div[role="alert"]').addClass('d-none');
});
$('input[name="contestTime"]').focus(function () {
    $('input[name="contestTime"]').removeClass('alert alert-danger');

    $('div[role="alert"]').removeClass('d-block');
    $('div[role="alert"]').addClass('d-none');
});
$('input[name="contestDuration"]').focus(function () {
    $('input[name="contestDuration"]').removeClass('alert alert-danger');

    $('div[role="alert"]').removeClass('d-block');
    $('div[role="alert"]').addClass('d-none');
});