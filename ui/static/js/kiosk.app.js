(function(){
    var ws = new WebSocket("ws://localhost:5000/ws/kiosks")
    ws.addEventListener('message', function(e){
        console.log(e.data);
    });
})();