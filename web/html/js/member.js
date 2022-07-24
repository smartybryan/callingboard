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
	if (jsonObject == null &&
		(endpoint === "newly-available" || endpoint === "focus-members")) {
		return;
	}

	let wardId = getAuthValueFromCookie().wardid;
	//TODO: remove this after testing
	wardId = "256137";

	jsonObject.forEach(function (member) {
		let memberElement = document.createElement('li');
		let memberParts = member.split(";")
		let memberName = encodeURI(memberParts[0]);
		let memberImage = wardId + "/" + memberName + ".jpg";
		// in order the drag the image to a calling as well as the member li element,
		// the member element and the thumbnail have the same id
		let memberHTMLDisplay = `
<div class='thumbnail-container'>
	<img class="thumbnail" draggable="true" ondragstart="drag(event)" 
	style="display: none" id="` + memberParts[0] + `" onload="this.style.display=''" src="` + memberImage + `">
</div>
<div>` + memberParts[0] + `</div>
`
		memberElement.innerHTML = memberHTMLDisplay;
		memberElement.classList.add(memberTypeClass(memberParts[1]))
		memberElement.classList.add("member-row");
		memberElement.setAttribute("id", memberParts[0]);
		memberElement.setAttribute("draggable", "true");
		memberElement.addEventListener("click", function () {
			memberSelected(memberElement);
		});
		membersElement.appendChild(memberElement);
	});

	filterMembers();
}

function displayMembersImageUploader_do(response) {
	if (response === 401) {
		makeTabDefault("authentication");
		focusDefaultTab();
		return;
	}
	clearContainer(document.getElementById("member-photos-list"));

	const membersElement = document.getElementById("member-photos-list");
	let jsonObject = JSON.parse(response);

	let wardId = getAuthValueFromCookie().wardid;
	//TODO: remove this after testing
	wardId = "256137";

	jsonObject.forEach(function (member) {
		let memberElement = document.createElement('li');
		let memberParts = member.split(";")
		let memberName = encodeURI(memberParts[0]);
		let memberImage = wardId + "/" + memberName + ".jpg";
		let memberHTMLWithUpload = `
<form method="POST" action="/v1/image-upload?member=` + memberName + `" enctype="multipart/form-data" novalidate class="box">
	<div class="box__input">
		<div class='thumbnail-container'>
			<img class="thumbnail" style="display: none" onload="this.style.display=''" src="` + memberImage + `">
		</div>
		<div>` + memberParts[0] + `</div>
		<input type="file" name="imageFile" id="file" class="box__file"/>
		<button type="submit" class="box__button">Upload</button>
	</div>
</form>		
`
		memberElement.innerHTML = memberHTMLWithUpload;
		memberElement.classList.add(memberTypeClass(memberParts[1]))
		memberElement.classList.add("member-row");
		memberElement.setAttribute("id", memberParts[0]);
		memberElement.addEventListener("click", function () {
			memberSelected(memberElement);
		});
		membersElement.appendChild(memberElement);
	});

	initImageForms(document)
}

function memberTypeClass(memberType) {
	switch (memberType) {
		case "1":
			return "male";
		case "2":
			return "female";
	}
	return "none";
}

function clearMembersPanel() {
	clearContainer(document.getElementById("members"));
	clearFilter();
	clearCallingsHeldByMember();
	addCallingsHeldByMemberMessage();
}

function memberSelected(element) {
	displayMemberCallings(element.id);
}

function filterMembers() {
	let filter = document.getElementById("member-filter").value.toLowerCase();
	const memberElements = document.getElementById("members").getElementsByTagName("li");

	for (let i = 0; i < memberElements.length; i++) {
		memberElements[i].classList.remove("filtered");
		if (!memberElements[i].id.toLowerCase().includes(filter)) {
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

function addCallingsHeldByMemberMessage() {
	let memberCallings = document.getElementById("member-callings");
	let memberCallingsMessage = document.createElement("li");
	memberCallingsMessage.innerHTML = MEMBER_CALLING_MESSAGE;
	memberCallings.appendChild(memberCallingsMessage);
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
		.then(data => {
		})
		.catch(error => {
			console.log(error);
		})
}

function populateMemberPhotoList() {
	apiCall("members")
		.then(data => {
			populateMemberPhotoList_do(data);
		})
		.catch(error => {
			console.log(error);
		})
}

function populateMemberPhotoList_do(response) {
	displayMembersImageUploader_do(response)
}
