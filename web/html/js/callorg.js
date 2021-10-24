window.onload = function () {
	setupTree()
};

//// tree functions ////

function setupTree() {
	const wardOrgs = document.getElementById("ward-organizations");

	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let jsonObject = JSON.parse(this.responseText);
			jsonObject.forEach(function (calling) {
				let container = document.getElementById(calling.Org);
				if (container == null) {
					container = document.createElement("li");
					let caret = document.createElement("span");
					caret.setAttribute("class", "caret");
					caret.innerText = calling.Org;
					container.appendChild(caret);
					wardOrgs.appendChild(container);
					let nested = document.createElement("ul");
					nested.setAttribute("class", "nested");
					nested.setAttribute("id", calling.Org);
					container.appendChild(nested);
					container = nested;
				}

				if (calling.SubOrg !== "") {
					let subOrg = document.getElementById(calling.Org + calling.SubOrg);
					if (subOrg == null) {
						subOrg = document.createElement("li");
						subOrg.setAttribute("id", calling.SubOrg)
						let caret = document.createElement("span");
						caret.setAttribute("class", "caret");
						caret.innerText = calling.SubOrg;
						subOrg.appendChild(caret);
						container.appendChild(subOrg);
						let nested = document.createElement("ul");
						nested.setAttribute("class", "nested");
						nested.setAttribute("id", calling.Org + calling.SubOrg);
						subOrg.appendChild(nested);
						container = nested;
					} else {
						container = subOrg;
					}
				}

				let callingInfo = document.createElement("li");
				callingInfo.innerText = calling.Name + " ; " + calling.Holder + " ; " + calling.PrintableSustained + " (" + calling.PrintableTimeInCalling + ")";
				container.appendChild(callingInfo);
			});

			startTreeListeners()
		}
	};

	let url = "callings?org=All+Organizations";
	xhttp.open("GET", "/v1/" + url);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

function startTreeListeners() {
	let caret = document.getElementsByClassName("caret");
	let i;

	for (i = 0; i < caret.length; i++) {
		caret[i].addEventListener("click", function () {
			if (this.innerText === "All Organizations") {
				toggleElement(this);
			} else {
				expandChildren(this);
			}
		});
	}
}

function expandCollapseTree(op) {
	let caret = document.getElementsByClassName("caret");
	let i;

	for (i = 0; i < caret.length; i++) {
		switch(op) {
			case 'e':
				caret[i].parentElement.querySelector(".nested").classList.add("active");
				caret[i].classList.add("caret-down");
				break;
			case 'c':
				caret[i].parentElement.querySelector(".nested").classList.remove("active");
				caret[i].classList.remove("caret-down");
				break;
		}
	}
}

function expandChildren(element) {
	let toggler = element.parentElement.getElementsByClassName("caret");
	let i;
	for (i = 0; i < toggler.length; i++) {
		toggleElement(toggler[i]);
	}
}

function toggleElement(element) {
	element.parentElement.querySelector(".nested").classList.toggle("active");
	element.classList.toggle("caret-down");
}

//// misc ////

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

//// parsers ////

function parseRawData(endpoint) {
	const rawData = document.getElementById("rawdata");
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4) {
			let msg = this.responseText;
			if (msg === "null\n") {
				msg = "Import successful"
			}
			alert(msg);
		}
	};
	xhttp.open("POST", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send(rawData.value);
}
