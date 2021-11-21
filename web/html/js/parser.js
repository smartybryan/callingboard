
function parseRawData(endpoint) {
	const rawData = document.getElementById("rawdata");
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4) {
			let msg = this.responseText;
			if (msg === "null\n") {
				notify(nSUCCESS, MESSAGE_IMPORT_SUCCESSFUL)
				document.getElementById("rawdata").value = "";
				resetModel();
				refreshFromModel();
			} else {
				notify(nERROR, MESSAGE_IMPORT_FAILURE)
			}
		}
	};
	xhttp.open("POST", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send(rawData.value);
}
