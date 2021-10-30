const VACANT = "Calling Vacant";
const NONE = "None";
const MOVED = "-moved";

window.onload = function () {
	setupTree()
};

//// tree functions ////

function setupTree() {
	const wardOrgs = document.getElementById("ward-organizations");

	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let counter = 0;
			let jsonObject = JSON.parse(this.responseText);
			jsonObject.forEach(function (calling) {
				counter++;
				let container = document.getElementById(calling.Org);
				if (!container) {
					container = document.createElement("li");

					let caret = document.createElement("span");
					caret.classList.add("caret");
					caret.innerText = calling.Org;
					container.appendChild(caret);
					wardOrgs.appendChild(container);

					let nested = document.createElement("ul");
					nested.classList.add("nested");
					nested.setAttribute("id", calling.Org);
					container.appendChild(nested);
					container = nested;
				}

				if (calling.SubOrg !== "") {
					let subOrg = document.getElementById(calling.Org + "@" + calling.SubOrg);
					if (!subOrg) {
						subOrg = document.createElement("li");
						subOrg.setAttribute("id", calling.SubOrg)

						let caret = document.createElement("span");
						caret.classList.add("caret");
						caret.innerText = calling.SubOrg;
						subOrg.appendChild(caret);
						container.appendChild(subOrg);

						let nested = document.createElement("ul");
						nested.classList.add("nested");
						nested.setAttribute("id", calling.Org + "@" + calling.SubOrg);
						subOrg.appendChild(nested);
						container = nested;
					} else {
						container = subOrg;
					}
				}
				container.setAttribute("ondrop", "drop(event)")
				container.setAttribute("ondragover", "allowDrop(event)")
				container.setAttribute("ondragstart", "drag(event)")

				let callingInfo = document.createElement("li");
				callingInfo.setAttribute("id", callingId(calling.Name, calling.Holder, counter))
				callingInfo.setAttribute("draggable", "true");

				callingInfo.classList.add("calling-row");
				if (calling.Holder === VACANT) {
					callingInfo.classList.add("vacant");
				}
				callingInfo.innerHTML = callingInnards(calling.Name, calling.Holder, calling.PrintableTimeInCalling)
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

//// calling id functions ////

function callingId(callingName, callingHolder, counter) {
	return callingName + "@" + callingHolder + "@" + counter;
}

function callingIdComponents(id) {
	let components = id.split("@");
	let callingName = components[0], holderName = components[1];
	return { callingName, holderName }
}

function callingInnards(callingName, holderName, timeInCalling) {
	return callingName + "<br><span class=\"member\">" + holderName + "</span> (" + timeInCalling + ")";
}

function callingIdMoved(id) {
	return id + MOVED;
}

function callingIdRemoveMoved(id) {
	return id.replace(MOVED,'');
}

//// tree functions ////

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
		switch (op) {
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

//// drag and drop ////

function allowDrop(ev) {
	ev.preventDefault();
}

function drag(ev) {
	ev.dataTransfer.setData("calling", ev.target.id);
}

function drop(ev) {
	ev.preventDefault();
	let currentId = ev.dataTransfer.getData("calling");

	let dropTarget = ev.target;
	if (dropTarget.tagName === "LI") {
		dropTarget = ev.target.parentElement
	} else if (dropTarget.tagName === "SPAN") {
		dropTarget = ev.target.parentElement.parentElement
	}

	//TODO: if the tree is both the source and destination, cancel the operation

	let movedElement = document.getElementById(currentId);
	let elementCopy = document.getElementById(currentId).cloneNode(true);

	// if dropTarget already has a moved element, use it and delete source element.
	let movedId = callingIdMoved(currentId);
	let targetElement = document.getElementById(movedId);
	if (targetElement) {
		targetElement.setAttribute("id", currentId)
		targetElement.innerHTML = movedElement.innerHTML
		targetElement.classList.remove("vacant")
		movedElement.remove()
		return
	}

	// finish the element copy and make source calling vacant
	movedElement.setAttribute("id", callingIdMoved(currentId));
	movedElement.innerHTML = callingInnards(callingIdComponents(currentId).callingName, VACANT, NONE);
	movedElement.classList.add("vacant")
	dropTarget.appendChild(elementCopy);
}

//// member functions ////

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
