let offset = new Date().getTimezoneOffset();
console.log(offset);

let timeZoneOff = document.getElementById('timeZoneOff');
timeZoneOff.innerHTML = offset * 60;