// see https://towardsdatascience.com/quickly-extract-all-links-from-a-web-page-using-javascript-and-the-browser-console-49bb6f48127b
var x = document.querySelectorAll('a[download]');
var linkTextToFilter = 'DESCARGAR COMPROBANTE'
var myarray = []
for (var i = 0; i < x.length; i++) {
    var linkText = x[i].textContent;
    var cleanText = linkText.replace(/\s+/g, ' ').trim();
    if (cleanText.trim() === linkTextToFilter) {
        var cleanLink = x[i].href;
        myarray.push([cleanLink]);
    }
};
function print_links() {
    var linksList = '';

    for (var i = 0; i < myarray.length; i++) {
        linksList += myarray[i][0] + '<br/>';
    };

    var w = window.open("");
    w.document.write(linksList);
}
print_links()