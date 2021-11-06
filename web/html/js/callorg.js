const VACANT = "Calling Vacant";
const ALL_ORGS = "All Organizations";
const NONE = "None";
const MOVED = "-moved";
const COPY = "-copy";
const PROPOSED = "-proposed";
const RELEASE_DROP_ENABLER = "(Drag calling here to release)";
const MESSAGE_RELEASE_ONLY = "You can only drag this calling into Releases.";
const MESSAGE_RELEASE_VACANT = "You cannot drag a vacant calling into Releases.";
const MESSAGE_RELEASE_DROP = "You may not drop a Released calling here. Drop it in the trash to undo.";
const MESSAGE_SUSTAIN_DROP = "You may not drop a Sustained calling here. Drop it in the trash to undo.";
const MESSAGE_MEMBER_ONLY_TO_TREE = "You can only drag this member into a Vacant Calling.";
const MESSAGE_ALREADY_EXISTS = "The member has already been released from this calling.";

window.onload = function () {
	setupTreeStructure();
};

//// tree functions ////

function setupTreeStructure() {
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

				refreshFromModel();
			});

			startTreeListeners()
		}
	};

	let url = "callings?org=" + encodeURI(ALL_ORGS);
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
	return {callingName, holderName}
}

function callingInnards(callingName, holderName, timeInCalling) {
	return callingName + "<br><span class=\"member\">" + holderName + "</span> (" + timeInCalling + ")";
}

function callingIdWithSuffix(id, suffix) {
	return id + suffix;
}

function callingIdRemoveSuffix(id, suffix) {
	return id.replace(suffix, '');
}

function convertHolderToClass(holder) {
	return holder.replace(",", "").replaceAll(" ", "")
}

//// tree functions ////

function startTreeListeners() {
	let caret = document.getElementsByClassName("caret");
	let i;

	for (i = 0; i < caret.length; i++) {
		caret[i].addEventListener("click", function () {
			if (this.innerText === ALL_ORGS) {
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

//// refresh functions ////

function refreshFromModel() {
	refreshTree();
	refreshCallingChanges();
}

function refreshTree() {
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let counter = 0;
			let jsonObject = JSON.parse(this.responseText);

			// clear all organizations
			jsonObject.forEach(function (calling) {
				let container = null;
				//TODO: refactor duplicate lines in here and setupTree()
				if (calling.SubOrg !== "") {
					container = document.getElementById(calling.Org + "@" + calling.SubOrg);
				} else {
					container = document.getElementById(calling.Org);
				}
				clearContainer(container);
			});

			// repopulate organizations
			jsonObject.forEach(function (calling) {
				counter++;
				let container = null;
				if (calling.SubOrg !== "") {
					container = document.getElementById(calling.Org + "@" + calling.SubOrg);
				} else {
					container = document.getElementById(calling.Org);
				}

				let callingInfo = document.createElement("li");
				callingInfo.setAttribute("id", callingId(calling.Name, calling.Holder, counter));
				callingInfo.setAttribute("data-org", calling.Org);
				callingInfo.setAttribute("draggable", "true");
				callingInfo.classList.add("calling-row");
				if (calling.Holder === VACANT) {
					callingInfo.classList.add("vacant");
				}
				callingInfo.innerHTML = callingInnards(calling.Name, calling.Holder, calling.PrintableTimeInCalling);
				container.appendChild(callingInfo);
			});
		}
	};

	let url = "callings?org=" + ALL_ORGS;
	xhttp.open("GET", "/v1/" + url);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

//TODO: refactor to remove duplicate code
function refreshCallingChanges() {
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let jsonObject = JSON.parse(this.responseText);

			let container = document.getElementById("releases");
			clearContainer(container);
			addReleaseDropEnabler(container);
			jsonObject.Releases.forEach(function (calling) {
				let callingInfo = document.createElement("li");
				callingInfo.setAttribute("id", callingId(calling.Name, calling.Holder, 0));
				callingInfo.setAttribute("data-org", calling.Org);
				callingInfo.setAttribute("draggable", "true");

				callingInfo.classList.add("calling-row");
				if (calling.Holder === VACANT) {
					callingInfo.classList.add("vacant");
				}
				callingInfo.innerHTML = callingInnards(calling.Name, calling.Holder, calling.PrintableTimeInCalling)
				container.appendChild(callingInfo);
			});

			container = document.getElementById("sustainings");
			clearContainer(container);
			jsonObject.Sustainings.forEach(function (calling) {
				let callingInfo = document.createElement("li");
				callingInfo.setAttribute("id", callingId(calling.Name, calling.Holder, 0));
				callingInfo.setAttribute("data-org", calling.Org);
				callingInfo.setAttribute("draggable", "true");

				callingInfo.classList.add("calling-row");
				if (calling.Holder === VACANT) {
					callingInfo.classList.add("vacant");
				}
				callingInfo.innerHTML = callingInnards(calling.Name, calling.Holder, calling.PrintableTimeInCalling)
				container.appendChild(callingInfo);
			});
		}
	};

	let url = "diff";
	xhttp.open("GET", "/v1/" + url);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

function addReleaseDropEnabler(container) {
	let dropEnabler = document.createElement("li");
	dropEnabler.innerHTML = RELEASE_DROP_ENABLER;
	container.appendChild(dropEnabler);
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
	let movedElement = document.getElementById(currentId);
	let newElement = document.getElementById(currentId).cloneNode(true);

	// find the UL and LI elements based on the dropped element
	let liElement = ev.target;
	if (liElement.tagName === "SPAN") {
		liElement = liElement.parentElement;
	}
	let dropTarget = liElement.parentElement;

	// dragging from the member row to a vacant calling in the tree
	if (movedElement.classList.contains("member-row")) {
		if (!liElement.classList.contains("vacant")) {
			alert(MESSAGE_MEMBER_ONLY_TO_TREE)
			return
		}
		transaction("add-member-calling", createTransactionParmsForMemberElement(movedElement, liElement));
		return
	}

	// dragging from releases
	if (movedElement.parentElement.id === "releases") {
		alert(MESSAGE_RELEASE_DROP);
		return
	}

	// dragging from sustainings
	if (movedElement.parentElement.id === "sustainings") {
		alert(MESSAGE_SUSTAIN_DROP);
		return
	}

	// dragging from ward tree to ward tree, cancel the operation
	if (dropTarget.classList.contains("nested") && movedElement.parentElement.classList.contains("nested")) {
		alert(MESSAGE_RELEASE_ONLY)
		return
	}

	// dragging from member-callings list
	if (movedElement.classList.contains("member-calling")) {
		// only allow it to be moved to releases
		if (dropTarget.id !== "releases") {
			alert(MESSAGE_RELEASE_ONLY);
			return
		}
		movedElement.remove();
	}

	// dragging a vacant calling from the tree
	if (movedElement.classList.contains("vacant")) {
		alert(MESSAGE_RELEASE_VACANT);
		return
	}

	transaction("remove-member-calling", createTransactionParmsFromTreeElememt(movedElement));
}


//// member functions ////

function displayMembers(endpoint) {
	const membersElement = document.getElementById("members");
	clearContainer(membersElement);

	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let jsonObject = JSON.parse(this.responseText)
			jsonObject.forEach(function (member) {
				let memberElement = document.createElement('li');
				memberElement.innerHTML = member;
				memberElement.setAttribute("id", member);
				memberElement.setAttribute("draggable", "true");
				memberElement.classList.add("member-row", convertHolderToClass(member));
				memberElement.addEventListener("click", function () {
					memberSelected(memberElement);
				});
				membersElement.appendChild(memberElement);
			});
		}
	};

	xhttp.open("GET", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

function memberSelected(element) {
	displayMemberCallings(element.id);
}

function displayMemberCallings(name) {
	const container = document.getElementById("member-callings");
	clearContainer(container);

	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let jsonObject = JSON.parse(this.responseText)
			jsonObject.forEach(function (calling) {
				let callingInfo = document.createElement("li");
				callingInfo.setAttribute("id", callingId(calling.Name, calling.Holder, 0));
				callingInfo.setAttribute("data-org", calling.Org);
				callingInfo.setAttribute("draggable", "true");

				callingInfo.classList.add("calling-row", "member-calling");
				callingInfo.innerHTML = callingInnards(calling.Name, calling.Holder, calling.PrintableTimeInCalling)
				container.appendChild(callingInfo);

			});
		}
	};

	let endpoint = "callings-for-member?member=" + name;
	xhttp.open("GET", "/v1/" + encodeURI(endpoint));
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

function clearContainer(element) {
	if (!element) {
		return;
	}
	while (element.firstChild) {
		element.lastChild.remove();
	}
}

//// transactions ////

function transaction(endpoint, parms) {
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4) {
			refreshFromModel();
		}
	};
	xhttp.open("GET", "/v1/" + endpoint + "?" + encodeURI(parms));
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

//TODO: remove dup code from the two functions below
function createTransactionParmsFromTreeElememt(element) {
	let callingIdParts = element.id.split("@");
	return "org=" + element.getAttribute("data-org") + "&calling=" + callingIdParts[0] + "&member=" + callingIdParts[1];
}

function createTransactionParmsForMemberElement(memberElement, callingElement) {
	let callingIdParts = callingElement.id.split("@");
	return "org=" + callingElement.getAttribute("data-org")  + "&calling=" + callingIdParts[0] + "&member=" + memberElement.id;
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
