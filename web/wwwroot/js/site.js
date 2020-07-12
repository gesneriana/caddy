$(document).ready(function () {
    var windowHeight = $(window).height();
    if (windowHeight > 800) {
        $("#menu").height(windowHeight)
    }
    $("#menu").width($("#menu").parent().width())
})

window.onresize = function () {
    var windowHeight = $(window).height();
    if (windowHeight > 800) {
        $("#menu").css("height", windowHeight)
    }
    $("#menu").width($("#menu").parent().width())
}