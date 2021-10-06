function parseRawData(endpoint) {
	const rawData = document.getElementById("rawdata");
	const xhttp = new XMLHttpRequest();
	xhttp.open("POST", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send(rawData.value);
}

function displayMembers(endpoint) {
	const membersElement = document.getElementById("members");
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			membersElement.innerText = this.responseText
		}
	};
	xhttp.open("GET", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

