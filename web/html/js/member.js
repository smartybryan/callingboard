
function displayMembers(endpoint) {
	apiCall(endpoint)
		.then(data => {
			displayMembers_do(data, endpoint);
		})
		.catch(error => {
			console.log(error);
			if (error === 401) {
				logout();
			} else {
				makeTabDefault("import");
				focusDefaultTab();
			}
		})
}

function displayMembers_do(response, endpoint) {
	if (response === 401) {
		makeTabDefault("authentication");
		focusDefaultTab();
		return;
	}
	clearMembersPanel()

	const membersElement = document.getElementById("members");
	let jsonObject = JSON.parse(response);
	if (jsonObject == null && endpoint === "newly-available") {
		return;
	}

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

	filterMembers();
}

function clearMembersPanel() {
	clearContainer(document.getElementById("members"));
	clearFilter();
	clearCallingsHeldByMember();
}

function memberSelected(element) {
	displayMemberCallings(element.id);
}

function filterMembers() {
	let filter = document.getElementById("member-filter").value.toLowerCase();
	const memberElements = document.getElementById("members").getElementsByTagName("li");

	for (let i = 0; i < memberElements.length; i++) {
		memberElements[i].classList.remove("filtered");
		if (!memberElements[i].id.toLowerCase().startsWith(filter)) {
			memberElements[i].classList.add("filtered");
		}
	}
}

function clearFilter() {
	document.getElementById("member-filter").value = "";
	filterMembers();
}

function displayMemberCallings(name) {
	apiCall("callings-for-member", "member=" + name)
		.then(data => {
			displayMemberCallings_do(data);
		})
		.catch(error => {
			console.log(error);
		})
}

function displayMemberCallings_do(response) {
	const container = document.getElementById("member-callings");
	clearContainer(container);

	let jsonObject = JSON.parse(response)
	jsonObject.forEach(function (calling) {
		let callingInfo = createCallingElement(calling, 0);
		callingInfo.classList.add("member-calling");
		container.appendChild(callingInfo);
	});
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

function populateFocusList() {
	apiCall("members-with-focus")
		.then(data => {
			populateFocusList_do(data);
		})
		.catch(error => {
			console.log(error);
		})
}

function populateFocusList_do(response) {
	const tableContainer = document.getElementById("focus-member-list");
	clearContainer(tableContainer);

	let jsonObject = JSON.parse(response)
	jsonObject.forEach(function (memberRecord) {
		let rowElement = document.createElement("tr")
		let memberElement = document.createElement("td");
		memberElement.innerHTML = memberRecord.Name;
		let focusElement = document.createElement("td")
		focusElement.classList.add("focus-column");
		let cbElement = document.createElement("input");
		cbElement.setAttribute("type", "checkbox");
		if (memberRecord.Focus) {
			cbElement.checked = true;
		}
		focusElement.appendChild(cbElement);
		rowElement.appendChild(memberElement);
		rowElement.appendChild(focusElement);
		tableContainer.appendChild(rowElement);
	});
}

function saveFocusList() {
	const tableContainer = document.getElementById("focus-member-list");
	let focusMembers = "";

	let rows = tableContainer.getElementsByTagName("tr");
	for (let row of rows) {
		let columns = row.getElementsByTagName("td");
		if (columns[1].getElementsByTagName("input")[0].checked) {
			if (focusMembers.length > 0) {
				focusMembers += "|";
			}
			focusMembers += columns[0].innerText;
		}
	}

	apiCall("put-focus-members", "member=" + focusMembers)
		.then(data => {})
		.catch(error => {
			console.log(error);
		})
}
