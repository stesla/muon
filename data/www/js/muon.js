(function() {
    var input_, output_, ws_;

    function connect(fn) {
        //TODO: wss when we're https
        var ws = new WebSocket("ws://" + window.location.host + "/connect");
        ws.onopen = function() {
            fn("Connected!");
        }
        ws.onclose = function() {
            fn("Disconnected!");
        }
        ws.onmessage = function(msg) {
            fn(msg.data);
        }
        return ws;
    }

    function processInput(evt) {
        if (this.value && evt.keyCode == 13) {
            evt.preventDefault();
            ws_.send(this.value);
            console.log("sent:", this.value);
           this.value = "";
        }
    }

    window.onload = function() {
        input_ = document.querySelector('#container .input textarea');
        output_ = document.querySelector('#container .output');
        input_.addEventListener('keydown', processInput);
        ws_ = connect(function(text) {
            var line = document.createElement('p');
            line.innerHTML = text;
            output_.appendChild(line);
            input_.scrollIntoView();
        });
    };
})();
