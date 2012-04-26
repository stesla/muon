(function() {
    function processInput(evt) {
        if (this.value && evt.keyCode == 13) {
            evt.preventDefault();
            console.log("input:", this.value);
            this.value = "";
        }
    }

    window.onload = function() {
        var input = document.querySelector('#container .input textarea');
        input.addEventListener('keydown', processInput);
    };
})();
