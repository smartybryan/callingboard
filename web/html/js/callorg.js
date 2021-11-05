const VACANT = "Calling Vacant";
const ALL_ORGS = "All Organizations";
const NONE = "None";
const MOVED = "-moved";
const COPY = "-copy";
const PROPOSED = "-proposed";
const RELEASE_DROP_ENABLER = "(Drag calling here to release)";
const MESSAGE_RELEASE_ONLY = "You can only drag this calling into Releases.";
const MESSAGE_RELEASE_DROP = "You may not drop a Released calling here. Drop it in the trash to undo.";
const MESSAGE_MEMBER_ONLY_TO_TREE = "You can only drag this member into a Vacant Calling.";
const MESSAGE_ALREADY_EXISTS = "The calling already exists at the drop zone.";
const MESSAGE_SUSTAIN_NOT_ALLOWED = "Sustain proposals cannot be dropped here.";

window.onload = function () {
	setupTree();
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

function replaceHolderInCallingInnards(element, holderName, timeInCalling) {
	let idComponents = element.id.split("@");
	return idComponents[0] + "<br><span class=\"member\">" + holderName + "</span> (" + timeInCalling + ")";
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

			jsonObject.forEach(function (calling) {
				counter++;
				let container = null;
				if (calling.SubOrg !== "") {
					container = document.getElementById(calling.Org + "@" + calling.SubOrg);
				} else {
					container = document.getElementById(calling.Org);
				}

				let callingInfo = document.createElement("li");
				callingInfo.setAttribute("id", callingId(calling.Name, calling.Holder, counter))
				callingInfo.setAttribute("draggable", "true");

				callingInfo.classList.add("calling-row", convertHolderToClass(calling.Holder));
				if (calling.Holder === VACANT) {
					callingInfo.classList.add("vacant");
				}
				callingInfo.innerHTML = callingInnards(calling.Name, calling.Holder, calling.PrintableTimeInCalling)
				container.appendChild(callingInfo);
			});
		}
	};

	let url = "callings?org=" + ALL_ORGS;
	xhttp.open("GET", "/v1/" + url);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

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
				callingInfo.setAttribute("draggable", "true");

				callingInfo.classList.add("calling-row", convertHolderToClass(calling.Holder));
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
				callingInfo.setAttribute("id", callingId(calling.Name, calling.Holder, 0))
				callingInfo.setAttribute("draggable", "true");

				callingInfo.classList.add("calling-row", convertHolderToClass(calling.Holder));
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
		if (!liElement.classList.contains("vacant") && !liElement.classList.contains("model-vacant")) {
			alert(MESSAGE_MEMBER_ONLY_TO_TREE)
			return
		}

		transaction("add-member-calling", createTransactionParmsForMemberElement(movedElement, liElement));

		// liElement.innerHTML = replaceHolderInCallingInnards(liElement, movedElement.innerHTML, "Proposed")
		// liElement.classList.remove("vacant", "model-vacant");
		// liElement.classList.add("proposed");
		//
		// let sustainingElement = liElement.cloneNode(true);
		// document.getElementById("sustainings").appendChild(sustainingElement);
		// liElement.id += PROPOSED;
		return
	}

	// dragging from releases
	if (movedElement.parentElement.id === "releases") {
		alert(MESSAGE_RELEASE_DROP);
		return
	}

	// dragging from sustainings
	if (movedElement.parentElement.id === "sustainings") {
		if (!dropTarget.classList.contains("nested")) {
			alert(MESSAGE_SUSTAIN_NOT_ALLOWED);
			return
		}

		let treeElement = document.getElementById(callingIdRemoveSuffix(movedElement.id, PROPOSED));
		movedElement.remove();
		//TODO: restore the element back to its vacant state (id, classes)
		//      This is a problem because we need to store the previous member in case they back out the sustain
		treeElement.classList.remove(PROPOSED);
		treeElement.classList.add(VACANT);
		treeElement.innerHTML = replaceHolderInCallingInnards(treeElement, VACANT, NONE)
	}

	// if the ward tree is both the source and destination, cancel the operation
	if (dropTarget.classList.contains("nested") && movedElement.parentElement.classList.contains("nested")) {
		alert(MESSAGE_RELEASE_ONLY)
		return
	}

	// if element came from member-callings list
	if (movedElement.classList.contains("member-calling")) {
		// only allow it to be moved to releases
		if (dropTarget.id !== "releases") {
			alert(MESSAGE_RELEASE_ONLY);
			return
		}

		// if element already exists in releases, ignore the drag
		currentId = callingIdRemoveSuffix(currentId, COPY);
		let checkExist = document.getElementById(currentId);
		if (checkExist && checkExist.parentElement.id === "releases") {
			alert(MESSAGE_ALREADY_EXISTS);
			return;
		}

		// otherwise it came FROM releases, so remove it and make the tree element the root
		movedElement.remove();
		movedElement = document.getElementById(currentId);
		newElement = document.getElementById(currentId).cloneNode(true);
	}

	// if dropTarget already has a moved element with same id, use it and delete source element.
	let movedId = callingIdWithSuffix(currentId, MOVED);
	let existingPreviouslyMovedElement = document.getElementById(movedId);
	if (existingPreviouslyMovedElement) {
		existingPreviouslyMovedElement.setAttribute("id", currentId)
		existingPreviouslyMovedElement.innerHTML = movedElement.innerHTML
		existingPreviouslyMovedElement.classList.remove("model-vacant")
		movedElement.remove()
		return
	}

	transaction("remove-member-calling", createTransactionParmsFromTreeElememt(movedElement));

	// fix up the old element
	// movedElement.setAttribute("id", callingIdWithSuffix(currentId, MOVED));
	// movedElement.innerHTML = callingInnards(callingIdComponents(currentId).callingName, VACANT, NONE);
	// movedElement.classList.add("model-vacant")

	// add a copy of the moved element to the new destination
	// dropTarget.appendChild(newElement);
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
	const parentElement = document.getElementById("member-callings");
	clearContainer(parentElement);

	let memberCallings = document.getElementsByClassName(convertHolderToClass(name));
	for (let i = 0; i < memberCallings.length; i++) {
		// skip copying member row element
		if (memberCallings[i].classList.contains("member-row")) {
			continue;
		}
		// skip copying element already released
		if (memberCallings[i].classList.contains("model-vacant")) {
			continue;
		}

		let currentId = memberCallings[i].id;
		let newElement = document.getElementById(currentId).cloneNode(true);
		newElement.classList.remove(convertHolderToClass(name)); // avoid endless loop as memberCallings will continue to grow
		newElement.classList.add("member-calling");
		newElement.id = callingIdWithSuffix(currentId, COPY);
		parentElement.appendChild(newElement);
	}
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

function createTransactionParmsFromTreeElememt(element) {
	let org = element.parentElement.id.split("@")[0];
	let callingIdParts = element.id.split("@");
	return "org=" + org + "&calling=" + callingIdParts[0] + "&member=" + callingIdParts[1];
}

function createTransactionParmsForMemberElement(memberElement, callingElement) {
	let org = callingElement.parentElement.id.split("@")[0];
	let callingIdParts = callingElement.id.split("@");
	return "org=" + org + "&calling=" + callingIdParts[0] + "&member=" + memberElement.id;
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
