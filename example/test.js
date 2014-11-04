function print(text) {
	var date = (new Date()).toString();

	var row = document.createElement("tr");
	row.innerHTML = "<td>" + date + "</td><td>" + text + "</td>";

	var output = document.getElementById("output")
	output.appendChild(row);
}

function test() {
	var session = NinchatClient.newSession();

	function hello(user_id) {
		var header = {
			action:       "send_message",
			user_id:      user_id,
			message_type: "ninchat.com/text"
		};

		var payload = [
			JSON.stringify({"text": "hello me"})
		];

		session.send(header, payload);
	}

	session.onSessionEvent(function(header) {
		print("SESSION: " + JSON.stringify(header));
		print("BINARY: " + session.binarySupported());

		if (header.event === "session_created") {
			hello(header.user_id);
		}
	});

	session.onEvent(function(header, payload) {
		print("EVENT: " + JSON.stringify(header));

		for (var i = 0; i < payload.length; i++) {
			print("PAYLOAD: " + NinchatClient.stringifyFrame(payload[i]));
		}
	});

	session.onConnState(function(state) {
		print("STATE: " + state);
	});

	session.onConnActive(function() {
		print("ACTIVITY");
	});

	session.onLog(function(message) {
		print("LOG: " + message);
	});

	if (true) {
		// don't even try websocket
		session.setTransport("longpoll");
	}

	var params = {
		message_types: [
			"ninchat.com/info/*",
			"ninchat.com/link",
			"ninchat.com/notice",
			"ninchat.com/text"
		]
	};

	if (false) {
		// log in with an identity
		params.identity_type = "email";
		params.identity_name = "bob@example.invalid";
		params.identity_auth = "password goes here";
	} else if (false) {
		// log in with a user agent authentication token
		params.user_id = "user id goes here";
		params.user_auth = "secret token goes here";
	} else if (false) {
		// create a new user and associate a new identity to it
		params.identity_type_new = "email";
		params.identity_name_new = "bob@example.invalid";
		params.identity_auth_new = "password goes here";
	} else {
		// create a new user without an identity
	}

	session.setParams(params);
	session.open();

	if (false) {
		setTimeout(session.close, 15000);
	}
}
