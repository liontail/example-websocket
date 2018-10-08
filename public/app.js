new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        username: null, // Our username
        joined: false // True if username have been filled in
    },
    created: function () {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function (e) {
            var message = JSON.parse(e.data);
            let username = self.$cookies.get('username');
            if (message.length) {
                message.forEach(function (msg) {
                    self.createMessage(username, msg);
                });
            } else {
                self.createMessage(username, message);
            }
        });

        if (self.$cookies.get('username')) {
            this.username = self.$cookies.get('username')
            this.joined = true;
        }
    },
    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text(), // Strip out html
                    }
                    ));
                this.newMsg = ''; // Reset newMsg
            }
        },
        join: function () {
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            let username = $('<p>').html(this.username).text();
            var self = this;
            var xhr = new XMLHttpRequest();
            var url = 'http://' + window.location.host + '/join';
            xhr.open('POST', url, true);
            xhr.setRequestHeader('Content-Type', 'application/json');
            xhr.onreadystatechange = function () {
                if (xhr.readyState === 4 && xhr.status === 200) {
                    if (xhr.responseText == 'ok') {
                        self.username = username
                        self.$cookies.set('username', username);
                        self.joined = true;
                    } else if (xhr.responseText == 'Duplicate') {
                        Materialize.toast('Username already have been chosen', 4000);
                        self.joined = false;
                        return
                    } else {
                        Materialize.toast(xhr.responseText, 4000);
                        self.joined = false;
                        return
                    }
                }
            };
            var data = JSON.stringify({ username: username });
            xhr.send(data);
        },
        createMessage: function (username, msg) {
            let div = document.createElement('div');
            let content = '';
            if (msg.username == username) {
                content += '<div class="row to-right margin-bottom-zero"> <div class="col s6 offset-s6">'
            } else {
                content += '<div class="row to-left margin-bottom-zero"> <div class="col s6">'
            }
            content += '<div class="col s12">' + msg.username + '</div>'
            if (msg.username == username) {
                content += '<div class="col s10 chip text-message my-message">'
            } else {
                content += '<div class="col s10 chip text-message">'
            }

            content += emojione.toImage(msg.message) + '</div>'; // Parse emojis

            if (msg.username == username) {
                content += '<div class="col s2 text-message my-time-message">' + msg.updated_at + '</div>'
            } else {
                content += '<div class="col s2 text-message time-message">' + msg.updated_at + '</div>'
            }
            content += '</div> </div>'
            div.innerHTML = content
            var element = document.getElementById('chat-messages');
            element.appendChild(div)
            element.scrollTop = element.scrollHeight // Auto scroll to the bottom
        }
    }
});