const VACANT = "Calling Vacant";
const NONE = "None";
const MOVED = "-moved";
const COPY = "-copy";
const MESSAGE_RELEASE_ONLY = "You can only drag this calling into Releases."
const MESSAGE_ALREADY_EXISTS = "The calling already exists at the drop zone."

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

				callingInfo.classList.add("calling-row", convertHolderToClass(calling.Holder));
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
	let movedElement = document.getElementById(currentId);
	let newElement = document.getElementById(currentId).cloneNode(true);

	// find the UL container of the dropped element
	let dropTarget = ev.target;
	if (dropTarget.tagName === "LI") {
		dropTarget = ev.target.parentElement
	} else if (dropTarget.tagName === "SPAN") {
		dropTarget = ev.target.parentElement.parentElement
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


	// fix up the old element
	movedElement.setAttribute("id", callingIdWithSuffix(currentId, MOVED));
	movedElement.innerHTML = callingInnards(callingIdComponents(currentId).callingName, VACANT, NONE);
	movedElement.classList.add("model-vacant")

	// add a copy of the moved element to the new destination
	dropTarget.appendChild(newElement);
}

//// member functions ////

function displayMembers(endpoint) {
	const membersElement = document.getElementById("members");
	//TODO: clear each element like member-callings
	membersElement.innerHTML = "";

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
					memberSelected(this);
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
	//displayCallings("callings-for-member", "member", memberOptions[index].text);
}

function displayMemberCallings(name) {
	const parentElement = document.getElementById("member-callings");
	while (parentElement.firstChild) {
		parentElement.lastChild.remove();
	}

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
