$(document).ready(function(){

    var ws = new WebSocket(url("/" + fbId));
    ws.onmessage = function(event) {
      msg = JSON.parse(event.data);
      handle(msg);
    };

    ws.onopen = function() {
      var msg = {
        type: "login",
        text: ""+fbId,
        type: ""+fbId
      };

//      ws.send(JSON.stringify(msg));
    };

    var handle = function(message){
        handlers[message.type](message);
    }

    var team = function(message) {
        $("#team").html(message.text);
    }

    var points_handler = function(message) {
        $("#points").html(message.text);
    }

    var event_handler = function(message){
        var str = $("#history").html();
        var new_line = "<h5 class='text-success'>" + message.text + "</h5>";
        str = new_line + str;
        $("#history").html(str);

    }

    var text_updater = function(query) {
        return function(message){
            $(query).html(message.text);
        }
    }

    var handlers = {"team": team, "points": points_handler, "event": event_handler, "captains": text_updater("#auctioners"), "upcoming": text_updater("#upcoming"),
    "current-player": text_updater("#current"), "current-header": text_updater("#current_ign")};

    $(function(){
        $('#bid_input').onEnter(function(e){
            e.preventDefault();
            submission(ws);
        });
    });
    $(function(){
        $('#bid').submit(function(e){
            e.preventDefault();
            submission(ws);
        });
        $("#bid_5").click(function(e){
            var msg = {
                type: "bid-more",
                from: ""+fbId,
                text: "5"
            };

            ws.send(JSON.stringify(msg));
        });

        $("#bid_1").click(function(){
            var msg = {
                type: "bid-more",
                from: ""+fbId,
                text: "1"
            };

            ws.send(JSON.stringify(msg));
        })
    });
});


function strStartsWith(str, prefix) {
    return str.indexOf(prefix) === 0;
}

function submission(ws){
    var form = $("#bid");
    var get_url = form.attr('action');
    var get_data = $("#bid_input").val();
    $("#bid_input").val("");

    var msg = {
        type: "bid",
        from: ""+fbId,
        text: get_data
    };

    ws.send(JSON.stringify(msg));
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

function url(s) {
    var l = window.location;
    return ((l.protocol === "https:") ? "wss://" : "ws://") + l.hostname + (((l.port != 80) && (l.port != 443)) ? ":" + l.port : "") + l.pathname + s;
}

(function($) {
    $.fn.onEnter = function(func) {
        this.bind('keypress', function(e) {
            if (e.keyCode == 13) func.apply(this, [e]);
        });
        return this;
     };
})(jQuery);
