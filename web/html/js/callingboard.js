const VACANT = "Calling Vacant";
const ALL_ORGS = "All Organizations";
const RELEASE_DROP_ENABLER = "[Drop a calling here] &#x2935;";
const SUSTAIN_DROP_MESSAGE = "[Drop member on a vacant calling] &#x2192;";
const MEMBER_CALLING_MESSAGE = "[Click on member] &#x2191;";
let currentMemberListEndpoint = ""; // the last member query endpoint
// We use the image version in the member list img tag to assure
// the cache is not misrepresenting the current state of the images on the server.
let imageVersion = 0; // the current image version

window.onload = function () {
	registerDefaultEvents();
	initialize();
	focusDefaultTab();
};

function initialize() {
	imageVersion = Date.now();
	checkLoginStatus();
	if (!authCookieExists()) {
		return;
	}
	setupTreeStructure();
	displayMembers("members-with-callings");
	listModels();
	populateFocusList();
}

function registerDefaultEvents() {
	// Pressing enter on wardid will login
	let input = document.getElementById("wardid");
	input.addEventListener("keypress", function (event) {
		if (event.key === "Enter") {
			event.preventDefault();
			login();
		}
	});
}

function openTab(evt, tabName) {
	tabPreEvent(tabName);
	let i, tabcontent, tablinks, leavingTab;
	tabcontent = document.getElementsByClassName("tabcontent");
	for (i = 0; i < tabcontent.length; i++) {
		tabcontent[i].style.display = "none";
	}
	tablinks = document.getElementsByClassName("tablinks");
	for (i = 0; i < tablinks.length; i++) {
		if (tablinks[i].classList.contains("active")) {
			leavingTab = tablinks[i].name;
		}
		tablinks[i].classList.remove("active");
	}
	document.getElementById(tabName).style.display = "flex";
	evt.currentTarget.classList.add("active");

	tabPostEvent(leavingTab);
}

function tabPreEvent(enteringTab) {
	switch (enteringTab) {
		case "report":
			generateReport();
			break;
		case "authentication":
			let element = document.getElementById("username");
			element.focus();
			break;
	}
}

function tabPostEvent(leavingTab) {
	switch (leavingTab) {
		case "manage-focus-members":
			saveFocusList();
			break;
	}
}

function makeTabDefault(tabName) {
	let tabs = document.getElementsByClassName("tablinks")
	for (let tab of tabs) {
		tab.classList.remove("default");
		if (tab.name === tabName) {
			tab.classList.add("default");
		}
	}
}

function focusDefaultTab() {
	let tabs = document.getElementsByClassName("default");
	if (tabs.length === 0) {
		return
	}
	tabs[0].click();
	tabPreEvent(tabs[0].name);
}

function clearModeling() {
	clearMembersPanel();
	clearCallingsReleases();
	clearCallingTree();
}

//// login functions ////

function login() {
	let username = document.getElementById("username").value;
	let wardid = document.getElementById("wardid").value.toLowerCase();
	let params = "username=" + username + "&wardid=" + wardid;

	apiCall("login", params)
		.then(data => {
			initialize();
			focusDefaultTab();
		})
		.catch(error => {
			console.log(error);
		});
}

function logout() {
	let authData = getAuthValueFromCookie();
	let params = "username=" + authData.username + "&wardid=" + authData.wardid;
	apiCall("logout", params)
		.then(data => {
			clearLoggedInUsername();
			document.cookie = "id=; Max-Age=-9999999"
			clearModeling();
			makeTabDefault("authentication");
			focusDefaultTab();
		})
		.catch(error => {
			console.log(error);
		});
}

function checkLoginStatus() {
	if (!authCookieExists()) {
		makeTabDefault("authentication");
		focusDefaultTab();
		return
	}
	let auth = getAuthValueFromCookie();
	document.getElementById("username").value = auth.username;
	document.getElementById("wardid").value = auth.wardid;
	setLoggedInUsername();
	makeTabDefault("modeling");
	focusDefaultTab();
}

function authCookieExists() {
	return document.cookie.length > 0;
}

function clearLoggedInUsername() {
	let usernameElements = document.getElementsByClassName("logged-in-username");
	for (let usernameElement of usernameElements) {
		usernameElement.innerHTML = "";
	}

	document.getElementById("login-panel").classList.remove("filtered")
	document.getElementById("logout-panel").classList.add("filtered")
	document.getElementById("header-logout-button").classList.add("filtered")
}

function setLoggedInUsername() {
	let usernameElements = document.getElementsByClassName("logged-in-username");
	let cookieData = getAuthValueFromCookie();
	for (let usernameElement of usernameElements) {
		usernameElement.innerHTML = "Logged in as " + cookieData.username;
	}

	document.getElementById("login-panel").classList.add("filtered")
	document.getElementById("logout-panel").classList.remove("filtered")
	document.getElementById("header-logout-button").classList.remove("filtered")
}

function getAuthValueFromCookie() {
	let username = "";
	let wardid = "";
	let cookieSplit = document.cookie.split("=");
	if (cookieSplit.length < 2) {
		return {username, wardid};
	}

	let cookieValueSplit = cookieSplit[1].split(":");
	if (cookieValueSplit.length < 2) {
		return {username, wardid};
	}
	username = cookieValueSplit[0];
	wardid = cookieValueSplit[1];
	return {username, wardid};
}

//// calling id functions ////

function callingId(callingName, callingHolder, counter) {
	return callingName + "@" + callingHolder + "@" + counter;
}

function callingInnards(callingName, holderName, timeInCalling) {
	let wardId = getAuthValueFromCookie().wardid;
	let memberName = encodeURI(holderName);
	let memberImage = wardId + "/" + memberName + ".jpg";
	let imageElement = "<span class=\"thumbnail-container\"><img id=\"" + holderName + "@img-calling\" class=\"thumbnail\" onload=\"this.style.display=''\" style=\"display: none\" src=\"" + memberImage + "?v=" + imageVersion + "\" alt></span>";

	return "<div class=\"calling-container\"><span class=\"calling-name-container\">" + callingName + "<br><span class=\"member-name indent\">" + holderName + "</span><br><span class=\"indent\">(" + timeInCalling + ")</span></span>" + imageElement + "</div>";
}

function callingIdComponents(id) {
	let components = id.split("@");
	let callingName = components[0], holderName = components[1];
	return {callingName, holderName}
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

	// keep moving up until dropTarget is a droppable element
	let dropTarget = ev.target;
	for (let i = 0; i < 10; i++) { // limit the attempts to find a droppable element
		if (dropTarget.classList.contains("droppable")) {
			break
		}
		dropTarget = dropTarget.parentElement
	}

	if (!dropTarget.classList.contains("droppable")) {
		notify(nINFO, "Cannot drop anything here.");
		return
	}

	// drop image from member into trash
	if (movedElement.id.endsWith("@img-member") && dropTarget.id === "trash") {
		let imageNameParts = movedElement.id.split("@");
		apiCall("image-delete", "member=" + encodeURI(imageNameParts[0]))
			.then(data => {
				notify(nSUCCESS, "Image deleted");
				imageVersion++;
				displayMembers("");
				refreshFromModel();
				refreshTree()
			})
			.catch(error => {
				notify(nERROR, error);
				imageVersion++;
				displayMembers("");
			});
		return
	}

	// dragging image from a calling - ignore event
	if (movedElement.id.endsWith("@img-calling")) {
		notify(nALERT, MESSAGE_DELETE_IMAGE_DRAG_FROM_MEMBER)
		return
	}

	// dragging from the member row to a vacant calling in the tree
	if (movedElement.classList.contains("member-row")) {
		if (dropTarget.tagName === "UL") {
			notify(nALERT, MESSAGE_DRAG_ONTO_CALLING)
			return
		}
		if (!dropTarget.classList.contains("vacant")) {
			notify(nALERT, MESSAGE_MEMBER_ONLY_TO_TREE)
			return
		}
		apiCall("add-member-calling", createTransactionParmsForMemberElement(movedElement, dropTarget))
			.then(data => {
				refreshFromModel();
				makeDirty()
			})
			.catch(error => {
				console.log(error);
			});
		return
	}

	// dragging from releases or sustainings to trash
	if (movedElement.parentElement.id === "releases" ||
		movedElement.parentElement.id === "sustainings") {
		if (ev.target.id !== "trash") {
			notify(nALERT, MESSAGE_RELEASE_SUSTAINED_DROP);
			return
		}
		let idComponents = callingIdComponents(movedElement.id);
		let params = "name=" + movedElement.parentElement.id + "&params=" + idComponents.holderName + ":" + movedElement.getAttribute("data-org") + ":" + movedElement.getAttribute("data-suborg") + ":" + idComponents.callingName;
		apiCall("backout-transaction", params)
			.then(data => {
				refreshFromModel();
				clearCallingsHeldByMember();
				makeDirty()
			})
			.catch(error => {
				console.log(error);
			})
		return
	}

	// dragging a vacant calling from the tree
	if (movedElement.classList.contains("vacant")) {
		notify(nALERT, MESSAGE_RELEASE_VACANT);
		return
	}

	// last resort to verify a drop to releases
	if (dropTarget.id !== "releases") {
		for (let i = 0; i < 10; i++) { // limit the attempts to find a droppable element
			if (dropTarget.id === "releases") {
				break
			}
			dropTarget = dropTarget.parentElement
		}
	}
	if (dropTarget.id !== "releases") {
		notify(nALERT, MESSAGE_RELEASE_ONLY);
		return
	}

	// dragging from member-callings list
	if (movedElement.classList.contains("member-calling")) {
		movedElement.remove();
	}

	apiCall("remove-member-calling", createTransactionParmsFromTreeElememt(movedElement))
		.then(data => {
			refreshFromModel();
			makeDirty();
		})
		.catch(error => {
			console.log(error);
		})
}

function makeDirty() {
	_manageSaveButtons('d')
}

function makeClean() {
	_manageSaveButtons('c')
}

function _manageSaveButtons(op) {
	let saveButtons = document.getElementsByClassName("save-button");
	for (let saveButton of saveButtons) {
		if (op === 'd') {
			saveButton.classList.add("dirty")
		} else {
			saveButton.classList.remove("dirty")
		}
	}
}

//// transactions ////

function createTransactionParmsFromTreeElememt(element) {
	let callingIdParts = element.id.split("@");
	return "org=" + element.getAttribute("data-org") + "&suborg=" + element.getAttribute("data-suborg") + "&calling=" + callingIdParts[0] + "&member=" + callingIdParts[1];
}

function createTransactionParmsForMemberElement(memberElement, callingElement) {
	let callingIdParts = callingElement.id.split("@");
	return "org=" + callingElement.getAttribute("data-org") + "&suborg=" + callingElement.getAttribute("data-suborg") + "&calling=" + callingIdParts[0] + "&member=" + memberElement.id;
}

function undoLast() {
	apiCall("undo-trans")
		.then(data => {
			refreshFromModel();
		})
		.catch(error => {
			console.log(error)
		})
}

function redoLast() {
	apiCall("redo-trans")
		.then(data => {
			refreshFromModel();
		})
		.catch(error => {
			console.log(error)
		})
}

//// api call ////

let apiCall = (endpoint, params) => {
	return new Promise((resolve, reject) => {
		const xhttp = new XMLHttpRequest();
		xhttp.open("GET", "/v1/" + endpoint + "?" + encodeURI(params));
		xhttp.setRequestHeader("Content-type", "text/plain");
		xhttp.onload = () => {
			if (xhttp.status >= 200 && xhttp.status < 300) {
				resolve(xhttp.responseText);
			} else {
				reject(xhttp.responseText);
			}
		}
		xhttp.onerror = () => {
			reject(xhttp.responseText);
		}
		xhttp.send();
	})
}
