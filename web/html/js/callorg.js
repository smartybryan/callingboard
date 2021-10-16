function parseRawData(endpoint) {
	const rawData = document.getElementById("rawdata");
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState === 4 && this.status === 422) {
			alert(this.responseText);
		}
	};
	xhttp.open("POST", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send(rawData.value);
}

function displayMembers(endpoint) {
	const membersElement = document.getElementById("members");
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState === 4 && this.status === 200) {
			membersElement.innerText = this.responseText;
		}
	};
	xhttp.open("GET", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}


function displayOrganizations(endpoint) {
	const orgElement = document.getElementById("organizations");
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState === 4 && this.status === 200) {
			orgElement.innerText = this.responseText;
		}
	};
	xhttp.open("GET", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

function displayCallings(endpoint) {
	const callingElement = document.getElementById("callings");
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState === 4 && this.status === 200) {
			callingElement.innerText = this.responseText;
		}
	};
	let url = endpoint;
	if (endpoint === 'callings') {
		url += "?org=" + encodeURI(document.getElementById("calling-org").value);
	}
	xhttp.open("GET", "/v1/" + url);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}
