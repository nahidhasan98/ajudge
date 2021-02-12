console.log("Script linked properly");

var OJ = $('#OJ').val();
var pNum = $('#pNum').val();
var languageBox = $('#submit-language');

addLanguage(OJ, pNum);

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
            var mapSorted = new Map([...map.entries()].sort());

            for (let [key, value] of mapSorted) {
                // console.log(key + ' = ' + value)
                languageBox.append($("<option></option>")
                    .attr("value", value)
                    .text(key));
            }
        },
    });
};

$(document).ready(function () {
    $('form').on('submit', function () {
        var sourceCode = $('#submit-solution');

        if (sourceCode.val().trim().length > 0) {
            return true;    //submit form
        }
        sourceCode.val("");
        return false;   //cancel submission
    });
});