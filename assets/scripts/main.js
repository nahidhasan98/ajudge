//for highlighting current page on navbar //
let url = window.location.pathname;
if (url == "/") url = window.location.origin + '/'; //for homepage '/' wil be 'http://localhost:8080/'

//now grab every link from the navigation
$('.navbar a').each(function () {
    let currNavLink = $(this).attr('href');
    if (currNavLink == "/") currNavLink = window.location.origin;   //for homepage

    navLinkRegExp = new RegExp(currNavLink); //create regexp to match current url pathname

    //checking wheather this navLink(regExp) present or not in the URL pathname
    if (navLinkRegExp.test(url)) {  //this means, string contains something like regexp.
        if (currNavLink == "/rankOJ" || currNavLink == "/rankUser") {
            $(this).parent().prev().addClass('activeNav');
        } else {
            $(this).addClass('activeNav');
        }
    }
});