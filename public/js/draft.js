$(document).ready(function(){
    $(function(){
        history_auto();
    })
    $(function(){
        $.get("/draft/status").done(function(data){
            updateCurrent(data);
            updateCaptainPoints(data);
            updateUpcoming(data);
        })
    });
    $(function(){
        $('#bid_input').onEnter(function(e){
            e.preventDefault();
            submission();
        });
    });
    $(function(){
        $('#bid').submit(function(e){
            e.preventDefault();
            submission();
        });
    });
});
var current_history = ""
function history_auto(){
    $.get("/draft/history").done(function(data){
        historyUpdater(data);
        if (strStartsWith(current_history, "WINNER:") || strStartsWith(current_history, "STARTING:") || strStartsWith(current_history, "NEXT:")) {
            everythingelse()
        }

    }, "json")
    setTimeout(history_auto, 400)
}

function everythingelse() {
    $.get("/draft/status").done(function(data){
            updateCurrent(data);
            updateCaptainPoints(data);
            updateUpcoming(data);
        })
}

function strStartsWith(str, prefix) {
    return str.indexOf(prefix) === 0;
}

function submission(){
    var form = $("#bid");
    var get_url = form.attr('action');
    var get_data = form.serialize();
    $("#bid_input").val("")
    $.ajax({
        type: 'GET',
        url: get_url, 
        data: get_data,
    });
}

function historyUpdater(history) {

    var str_html = "";
    $.each(history, function(num, line){
        if (num == 0) {
            current_history = line;
            str_html += "<h5 class='text-success'>" + line + "</h5>"
        } else {
            str_html += "<h5 class='text-muted'>" + line + "</h5>"
        }
    });

    $("#history").html(str_html)
}

function updateUpcoming(data) {
    var str_html = "";
    $.each(data.Unassigned, function(num, person){
        str_html += "<li class='list-group-item'>" + person.Ign + "  <span class='text-muted'>"+ person.Player.Tier +"</span></li>"
    });
    $("#upcoming").html(str_html);
}

function updateCaptainPoints(data) {
    captains = data.Auctioners;
    var arr = [];
    var str_html = "";
    $.each(captains, function(key, auctioner){
        if (auctioner.Team == "") {
            return
        }
        arr.push({key:"<li class='list-group-item'>" + auctioner.Team + " <span class='text-info'>"+ auctioner.Points +"</span></li>", value: auctioner.Points});
    });

    var sorted = arr.slice(0).sort(function(a, b) {
        return a.value - b.value;
    });

    var keys = [];
    for (var i = 0, len = sorted.length; i < len; ++i) {
        keys[i] = sorted[i].key;
    }

    keys.reverse();
   $("#auctionerpoints").html(keys.join(""));
}

function updateCurrent(data){
    var team = $("#team").text()
    $("#points").text(data.Auctioners[team].Points)
    $("#current_ign").text(data.Current.Ign)
    $("#current_name").text(data.Current.Player.Name)
    $("#current_tier").text(data.Current.Player.Tier)
    $("#current_role").text(data.Current.Player.RoleDescription)
    $("#current_lolking").text(data.Current.Player.Lolking)
    $("#current_lolking").attr("href", "http://www.lolking.net/summoner/na/"+data.Current.Id)
}

(function($) {
    $.fn.onEnter = function(func) {
        this.bind('keypress', function(e) {
            if (e.keyCode == 13) func.apply(this, [e]);    
        });               
        return this; 
     };
})(jQuery);
