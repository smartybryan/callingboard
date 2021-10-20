function parseRawData(endpoint) {
	const rawData = document.getElementById("rawdata");
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
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
	membersElement.innerHTML = "";

	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let jsonObject = JSON.parse(this.responseText)
			jsonObject.forEach(function (member) {
				let opt = document.createElement('option');
				opt.value = member;
				opt.innerText = member;
				membersElement.appendChild(opt);
			});
		}
	};
	xhttp.open("GET", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

function memberSelected(index) {
	let memberOptions = document.getElementById('members').options;
	displayCallings("callings-for-member", "member", memberOptions[index].text);
}

function displayOrganizations(endpoint) {
	const orgElement = document.getElementById("organizations");
	orgElement.innerHTML = "";
	let opt = document.createElement('option');
	opt.value = "All Organizations";
	opt.innerText = "All Organizations";
	orgElement.appendChild(opt);

	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let jsonObject = JSON.parse(this.responseText)
			jsonObject.forEach(function (member) {
				let opt = document.createElement('option');
				opt.value = member;
				opt.innerText = member;
				orgElement.appendChild(opt);
			});
		}
	};
	xhttp.open("GET", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

function orgSelected(index) {
	let orgOptions = document.getElementById('organizations').options;
	displayCallings("callings", "org", orgOptions[index].text)
}

function displayCallings(endpoint, argName, arg) {
	const callingElement = document.getElementById("callings");
	callingElement.innerHTML = "";

	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let jsonObject = JSON.parse(this.responseText)
			jsonObject.forEach(function (calling) {
				let opt = document.createElement('option');
				opt.value = calling.Name;
				opt.innerText = calling.Org + "/" + calling.SubOrg + " ; " + calling.Name + " ; " + calling.Holder + " ; " + calling.PrintableSustained + " (" + calling.PrintableTimeInCalling + ")";
				callingElement.appendChild(opt);
			});
		}
	};

	let url = endpoint + "?" + argName + "=" + encodeURI(arg);
	xhttp.open("GET", "/v1/" + url);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}
