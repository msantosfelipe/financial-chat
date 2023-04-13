$(function () {
    let websocket;
    let room;

    $("#input-form").on("submit", function (event) {
        event.preventDefault();

        if (!websocket) {
            // Create new connection
            room = $("#input-room")[0].value
            let url = "ws://" + window.location.host + "/websocket?room=" + room;
            websocket = new WebSocket(url);
            chatText = $("#chat-text");

            websocket.addEventListener("message", function (e) {
                let data = JSON.parse(e.data);
                let chatContent = `<p><strong>${data.timestamp} - ${data.username}</strong>: ${data.text}</p>`;
                chatText.append(chatContent);
                chatText.scrollTop = chatText.scrollHeight;
            });

            // disable inputs
            $("#input-username").prop('disabled', true);
            $("#input-room").prop('disabled', true);
            $("#connect-btn").prop('disabled', true);

            // show message textbox
            $("#div-text").prop('hidden', false);
            $("#message-btn").prop('hidden', false);
        } else {
            // Send message
            let username = $("#input-username")[0].value;
            let room = $("#input-room")[0].value;
            let text = $("#input-text")[0].value;

            websocket.send(
                JSON.stringify({
                    username: username,
                    room: room,
                    text: text,
                })
            );
            $("#input-text")[0].value = "";
        }
    });
});