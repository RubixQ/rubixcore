(function(){
    var ws = new WebSocket("ws://localhost:5000/ws/status")
    ws.addEventListener('message', function(e){
        var list = document.querySelector("#messages");
        var li = document.createElement("li");
        var message = document.createTextNode(e.data);
        li.appendChild(message)

        list.appendChild(li)
    });
})();