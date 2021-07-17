//checking wheather duration valid or not
let duration = $('input[name="contestDuration"]').val().trim();
let len = duration.length;
let flag = 1;

if (len >= 5 && len <= 6) {   //like 05:00 or 120:00
    if (duration[len - 1] >= '0' && duration[len - 1] <= '9' && duration[len - 2] >= '0' && duration[len - 2] <= '5' && duration[len - 3] == ':') {
        for (i = len - 4; i >= 0; i--) {
            if (duration[i] <= '0' && duration[i] >= '9') {
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