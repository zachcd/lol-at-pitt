$(document).ready(function(){
    $(function(){
        $.get("/draft/history").done(function(data){
            historyUpdater(data);
        }, "json")
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
        str_html += "<li class='list-group-item'>" + person.Ign + "  <span class='text-info'>" + person.Player.Lolking + "</span></li>"
    });
    $("#upcoming").html(str_html);
}

function updateCaptainPoints(data) {
    captains = data.Auctioners;
    var str_html = "";
    $.each(captains, function(key, auctioner){
        if (auctioner.Team == "") {
            return
        }
        str_html += "<li class='list-group-item'>" + auctioner.Team + " <span class='text-info'>"+ auctioner.Points +"</span></li>";
    });

   $("#auctionerpoints").html(str_html);
}

function updateCurrent(data){
    $("#current_ign").text(data.Current.Ign)
    $("#current_name").text(data.Current.Player.Name)
    $("#current_tier").text(data.Current.Player.Tier)
    $("#current_role").text(data.Current.Player.NormalizedIgn)
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
