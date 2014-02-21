(function(window) {
    var ws_;

    function connect(fn) {
        //TODO: wss when we're https
        var ws = new WebSocket("ws://" + window.location.host + "/connect" + window.location.search);
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

    function setTextAreaWidthInCharacters(width) {
        /* Presently, the only
         * reliable way to do that is to actually get the width of one
         * character as rendered on the screen because none of the CSS
         * units work right. */
        var oneChar = document.getElementById('oneChar');
        var sheets = document.styleSheets;
        var sheet = sheets[sheets.length - 1];
        sheet.insertRule('#container { width: ' + width * oneChar.clientWidth + 'px }', 0);
    }

    window.addEventListener('load', function() {
        var input = document.querySelector('#container .input textarea');
        var output = document.querySelector('#container .output');

        setTextAreaWidthInCharacters(80);

        input.addEventListener('keydown', processInput);

        ws_ = connect(function(text) {
            var line = document.createElement('p');
            line.innerHTML = text;
            output.appendChild(line);
            input.scrollIntoView();
        });
    });
})(window);
