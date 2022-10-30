function displayMembers(endpoint) {
	if (endpoint && endpoint.length > 0) {
		currentMemberListEndpoint = endpoint;
	} else {
		endpoint = currentMemberListEndpoint;
	}

	setButtonHighlight(endpoint);

	apiCall(endpoint)
		.then(data => {
			displayMembersImageUploader_do(data, endpoint);
		})
		.catch(error => {
			notify(nERROR, error);
			if (error === 401 || error === "Not logged in") {
				logout();
			} else {
				makeTabDefault("import");
				focusDefaultTab();
			}
		})
}

function displayMembersImageUploader_do(response, endpoint) {
	if (response === 401 || response === "Not logged in") {
		logout();
		return;
	}
	clearMembersPanel();

	const membersElement = document.getElementById("members");
	let jsonObject = JSON.parse(response);
	if (jsonObject == null &&
		(endpoint === "newly-available" || endpoint === "focus-members")) {
		return;
	}

	let wardId = getAuthValueFromCookie().wardid;

	if (document.getElementById("member-sort-firstname").checked) {
		jsonObject.sort(function (a, b) {
			return compareFirstNames(a.Name, b.Name);
		});
	}

	jsonObject.forEach(function (member) {
		let memberElement = document.createElement('li');
		let focusChecked = member.Focus ? "checked" : "";
		let memberParts = member.Name.split(";")
		let memberName = encodeURI(memberParts[0]);
		let memberImage = wardId + "/" + memberName + ".jpg";
		let memberNameParts = memberParts[0].split(",");
		if (memberNameParts.length < 2) {
			return false;
		}
		let memberLast = memberNameParts[0]
		let memberFirst = memberNameParts[1]
		// in order the drag the image to a calling as well as the member li element,
		// the member element and the thumbnail have the same id
		let memberHTMLWithUpload = `
<form method="POST" action="/v1/image-upload?member=` + memberName + `" enctype="multipart/form-data" novalidate class="box">
	<div class="box__input member-container">
		<span class="member-name-container">
			<div class="name-margins">
				<div class="member-name">` + memberLast + `,</div>
				<div class="indent-small">` + memberFirst + `</div>
			</div>
		</span>
		<span class='thumbnail-container'>
			<img class="thumbnail" draggable="true" ondragstart="drag(event)" 
			onload="this.style.display=''" style="display: none" id="` + memberParts[0] + `@img-member"
			src="` + memberImage + `?v=` + imageVersion + `">
			<img class="box__uploading" alt="" src="loading.gif">
		</span>
		<input type="file" name="imageFile" id="file" class="box__file"/>
		<button type="submit" class="box__button">Upload</button>
		<input id="` + memberParts[0] + `-focus" type="checkbox" onclick="setMemberFocus('` + memberParts[0] + `-focus')" title="Focus" ` + focusChecked + `/>
	</div>
</form>		
`
		memberElement.innerHTML = memberHTMLWithUpload;
		memberElement.classList.add(memberTypeClass(memberParts[1]))
		memberElement.classList.add("member-row");
		memberElement.setAttribute("id", memberParts[0]);
		memberElement.setAttribute("draggable", "true");
		memberElement.addEventListener("click", function () {
			memberSelected(memberElement);
		});
		membersElement.appendChild(memberElement);
	});

	initImageForms(document)
	filterMembers();
}

function compareFirstNames(a, b) {
	return buildFirstName(a).localeCompare(buildFirstName(b));
}

function buildFirstName(name) {
	let parts = name.split(";");
	if (parts.length < 2) {
		return name;
	}
	let name_parts = parts[0].split(",");
	if (name_parts.length < 2) {
		return name
	}
	return name_parts[1] + " " + name_parts[0];
}

function setButtonHighlight(endpoint) {
	let buttonIDs = ["ml-mem", "ml-new", "ml-foc", "ml-cal", "ml-awc"];
	let highlightClass = "member-button-active";

	for (let idx in buttonIDs) {
		document.getElementById(buttonIDs[idx]).classList.remove(highlightClass);
	}

	let buttonID = "";
	switch (endpoint) {
		case "members":
			buttonID = buttonIDs[0];
			break;
		case "newly-available":
			buttonID = buttonIDs[1];
			break;
		case "focus-members":
			buttonID = buttonIDs[2];
			break;
		case "members-with-callings":
			buttonID = buttonIDs[3];
			break;
		case "adults-without-calling":
			buttonID = buttonIDs[4];
			break;
	}
	document.getElementById(buttonID).classList.add(highlightClass);
}

function setMemberFocus(id) {
	let nameParts = id.split("-")
	let focusState = document.getElementById(id).checked
	apiCall("set-member-focus", "member=" + nameParts[0] + "&custom=" + focusState)
		.then(data => {})
		.catch(error => {
			console.log(error);
		})
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
	let male = document.getElementById("member-male").checked;
	let female = document.getElementById("member-female").checked;
	const memberElements = document.getElementById("members").getElementsByTagName("li");

	let count = memberElements.length;
	for (let i = 0; i < memberElements.length; i++) {
		memberElements[i].classList.remove("filtered");
		if ((!memberElements[i].id.toLowerCase().includes(filter)) ||
			(!male && memberElements[i].classList.contains("male")) ||
			(!female && memberElements[i].classList.contains("female"))) {
			memberElements[i].classList.add("filtered");
			count--
		}
	}

	document.getElementById("member-count").innerText = "" + count;
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
