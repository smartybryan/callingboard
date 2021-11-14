const VACANT = "Calling Vacant";
const ALL_ORGS = "All Organizations";
const RELEASE_DROP_ENABLER = "[Drop calling here] &#10549;";
const MESSAGE_RELEASE_ONLY = "You can only drag this calling into Releases.";
const MESSAGE_RELEASE_VACANT = "You cannot drag a vacant calling into Releases.";
const MESSAGE_RELEASE_SUSTAINED_DROP = "You may not drop a Released or Sustained calling here. Drop it in the Trash to undo.";
const MESSAGE_MEMBER_ONLY_TO_TREE = "You can only drag this member into a Vacant Calling.";
const MESSAGE_MODEL_RESET = "Data has been reloaded and the model has been cleared.";

window.onload = function () {
	setupTreeStructure();
	listModels();
	document.getElementById("default-tab").click();
};

function openTab(evt, tabName) {
	var i, tabcontent, tablinks;
	tabcontent = document.getElementsByClassName("tabcontent");
	for (i = 0; i < tabcontent.length; i++) {
		tabcontent[i].style.display = "none";
	}
	tablinks = document.getElementsByClassName("tablinks");
	for (i = 0; i < tablinks.length; i++) {
		tablinks[i].classList.remove("active");
	}
	document.getElementById(tabName).style.display = "block";
	evt.currentTarget.classList.add("active");
}

//// tree functions ////

function setupTreeStructure() {
	const wardOrgs = document.getElementById("ward-organizations");
	const xhttp = new XMLHttpRequest();

	function setupNestedContainers(containerName, parentContainer, orgContainer, containerId) {
		let caret = document.createElement("span");
		caret.classList.add("caret");
		caret.innerText = containerName;
		orgContainer.appendChild(caret);
		parentContainer.appendChild(orgContainer);

		let nested = document.createElement("ul");
		nested.classList.add("nested");
		nested.setAttribute("id", containerId);
		orgContainer.appendChild(nested);
		return nested;
	}

	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let counter = 0;
			let jsonObject = JSON.parse(this.responseText);
			jsonObject.forEach(function (calling) {
				counter++;
				let orgContainer = document.getElementById(calling.Org);
				if (!orgContainer) {
					let containerId = calling.Org;
					orgContainer = document.createElement("li");
					orgContainer = setupNestedContainers(calling.Org, wardOrgs, orgContainer, containerId);
				}

				if (calling.SubOrg !== "") {
					let containerId = calling.Org + "@" + calling.SubOrg;
					let subOrgContainer = document.getElementById(containerId);
					if (!subOrgContainer) {
						subOrgContainer = document.createElement("li");
						subOrgContainer.setAttribute("id", calling.SubOrg)
						orgContainer = setupNestedContainers(calling.SubOrg, orgContainer, subOrgContainer, containerId);
					} else {
						orgContainer = subOrgContainer;
					}
				}
				orgContainer.classList.add("leaf-container")
				orgContainer.setAttribute("ondrop", "drop(event)")
				orgContainer.setAttribute("ondragover", "allowDrop(event)")
				orgContainer.setAttribute("ondragstart", "drag(event)")
			});

			startTreeListeners();
			refreshFromModel();
			expandCollapseTree('c');
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

function callingInnards(callingName, holderName, timeInCalling) {
	return callingName + "<br><span class=\"member\">" + holderName + "</span> (" + timeInCalling + ")";
}

function callingIdComponents(id) {
	let components = id.split("@");
	let callingName = components[0], holderName = components[1];
	return {callingName, holderName}
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
	// expand all organizations by default
	if (op === 'c') {
		toggleElement(document.getElementById("ward-organizations"));
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

function createCallingElement(calling, counter) {
	let callingInfo = document.createElement("li");
	callingInfo.setAttribute("id", callingId(calling.Name, calling.Holder, counter));
	callingInfo.setAttribute("data-org", calling.Org);
	callingInfo.setAttribute("draggable", "true");
	callingInfo.classList.add("calling-row");
	if (calling.Holder === VACANT) {
		callingInfo.classList.add("vacant");
	}
	callingInfo.innerHTML = callingInnards(calling.Name, calling.Holder, calling.PrintableTimeInCalling);
	return callingInfo;
}

function refreshTree() {
	let showNewVacancies = document.getElementById("new-vacancies").checked
	let showAllVacancies = document.getElementById("all-vacancies").checked

	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			let counter = 0;
			let jsonObject = JSON.parse(this.responseText);
			let workingObjects = jsonObject;
			if (showNewVacancies) {
				workingObjects = jsonObject.NewVacancies;
			}
			// clear all organizations
			let leafContainers = document.getElementsByClassName("leaf-container");
			for (let leafContainer of leafContainers) {
				clearContainer(leafContainer);
			}

			// repopulate organizations
			workingObjects.forEach(function (calling) {
				if (showAllVacancies && calling.Holder != VACANT) {
					return;
				}
				counter++;
				let container = findContainerFromCalling(calling);
				container.appendChild(createCallingElement(calling, counter));
			});
		}
	};

	let url = "callings?org=" + ALL_ORGS;
	if (showNewVacancies) {
		url = "diff"
	}
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
				container.appendChild(createCallingElement(calling, 0));
			});

			container = document.getElementById("sustainings");
			clearContainer(container);
			jsonObject.Sustainings.forEach(function (calling) {
				container.appendChild(createCallingElement(calling, 0));
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

function findContainerFromCalling(calling) {
	if (calling.SubOrg !== "") {
		return document.getElementById(calling.Org + "@" + calling.SubOrg);
	} else {
		return document.getElementById(calling.Org);
	}
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

	// dragging from releases or sustainings to trash
	if (movedElement.parentElement.id === "releases" || movedElement.parentElement.id === "sustainings") {
		if (ev.target.id !== "trash") {
			alert(MESSAGE_RELEASE_SUSTAINED_DROP);
			return
		}
		let idComponents = callingIdComponents(movedElement.id);
		let params = "name=" + movedElement.parentElement.id + "&params=" + idComponents.holderName + ":" + movedElement.getAttribute("data-org") + ":" + idComponents.callingName;
		transaction("backout-transaction", params);
		clearCallingsHeldByMember();
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
				memberElement.classList.add("member-row");
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
				let callingInfo = createCallingElement(calling, 0);
				callingInfo.classList.add("member-calling");
				container.appendChild(callingInfo);
			});
		}
	};

	let endpoint = "callings-for-member?member=" + name;
	xhttp.open("GET", "/v1/" + encodeURI(endpoint));
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}

function clearCallingsHeldByMember() {
	clearContainer(document.getElementById("member-callings"));
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

function transaction(endpoint, params) {
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4) {
			refreshFromModel();
		}
	};
	xhttp.open("GET", "/v1/" + endpoint + "?" + encodeURI(params));
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

//// model file operations ////

function listModels() {
	modelOperation("list-trans-files")
}

function updateFileList(response) {
	let transFileList = document.getElementById("model-names")
	clearContainer(transFileList)

	let jsonObject = JSON.parse(response)
	jsonObject.forEach(function (name) {
		let tr = document.createElement('tr');
		let td = document.createElement('td');
		td.innerText = name;
		tr.appendChild(td)
		transFileList.appendChild(tr);
	});
}

function selectModelFile(ev) {
	let currentSelected = ev.target.parentElement.parentElement.getElementsByClassName("selected")
	if (currentSelected.length > 0) {
		currentSelected[0].classList.remove("selected")
	}
	ev.target.classList.add("selected")
}

function getSelectedModelFile() {
	let table = document.getElementById("model-names")
	let selected = table.getElementsByClassName("selected")
	if (selected.length > 0) {
		return selected[0].innerHTML;
	}
	return ""
}

function loadModel() {
	let name = getSelectedModelFile();
	if (!name) {
		alert("Please select a model name to load.");
		return
	}
	modelOperation("load-trans", name)
}

function deleteModel() {
	let name = getSelectedModelFile();
	if (!name) {
		alert("Please select a model name to delete.");
		return
	}
	modelOperation("delete-trans", name)
}

function saveModel() {
	let name = document.getElementById("model-name").value;
	if (!name) {
		alert("Please provide a model name to save.");
		return
	}
	modelOperation("save-trans", name)
}

function resetModel() {
	modelOperation("reset-model");
	alert(MESSAGE_MODEL_RESET);
}

function modelOperation(endpoint, name) {
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4 && this.status === 200) {
			if (endpoint === "list-trans-files") {
				updateFileList(this.responseText);
			} else if (this.responseText === "null\n") {
				refreshFromModel();
				listModels();
				if (endpoint === "load-trans") {
					alert("Model loaded");
				}
			} else {
				alert(this.responseText)
			}
		}
	};

	let params = ""
	if (name) {
		params = "?name=" + name;
	}
	xhttp.open("GET", "/v1/" + endpoint + encodeURI(params));
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
			document.getElementById("rawdata").value = "";
			resetModel();
			refreshFromModel();
		}
	};
	xhttp.open("POST", "/v1/" + endpoint);
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send(rawData.value);
}
