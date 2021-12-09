const VACANT = "Calling Vacant";
const ALL_ORGS = "All Organizations";
const RELEASE_DROP_ENABLER = "[Drop calling here] &#10549;";

window.onload = function () {
	initialize()
	focusDefaultTab();
};

function initialize() {
	checkLoginStatus();
	setupTreeStructure();
	displayMembers("members-with-callings");
	listModels();
	populateFocusList();
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
}

//// login functions ////

function login() {
	let username = document.getElementById("username").value;
	let wardid = document.getElementById("wardid").value;
	let params = "username=" + username + "&wardid=" + wardid;
	apiCall("login", params, login_callback);
}

function login_callback(response) {
	initialize();
	focusDefaultTab();
}

function logout() {
	let authData = getAuthValueFromCookie();
	let params = "username=" + authData.username + "&wardid=" + authData.wardid;
	apiCall("logout", params, logout_callback)
	document.cookie = "id=; Max-Age=-9999999"
}

function logout_callback(response) {
	initialize();
	clearLoggedInUsername();
	//TODO: clear the modeling data on the page
	//TODO: prepare the html for login

}

function checkLoginStatus() {
	if (document.cookie.length === 0) {
		makeTabDefault("authentication");
		return
	}
	let auth = getAuthValueFromCookie();
	document.getElementById("username").value = auth.username;
	document.getElementById("wardid").value = auth.wardid;

	setLoggedInUsername();
	makeTabDefault("modeling");
	focusDefaultTab();
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
	return callingName + "<br><span class=\"member-name\">" + holderName + "</span><br><span class=\"time-in-calling\">(" + timeInCalling + ")</span>";
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

	// find the UL and LI elements based on the dropped element
	let liElement = ev.target;
	if (liElement.tagName === "SPAN") {
		liElement = liElement.parentElement;
	}
	let dropTarget = liElement.parentElement;

	// dragging from the member row to a vacant calling in the tree
	if (movedElement.classList.contains("member-row")) {
		if (!liElement.classList.contains("vacant")) {
			notify(nALERT, MESSAGE_MEMBER_ONLY_TO_TREE)
			return
		}
		apiCall("add-member-calling", createTransactionParmsForMemberElement(movedElement, liElement), refreshFromModel);
		makeDirty()
		return
	}

	// dragging from releases or sustainings to trash
	if (movedElement.parentElement.id === "releases" || movedElement.parentElement.id === "sustainings") {
		if (ev.target.id !== "trash") {
			notify(nALERT, MESSAGE_RELEASE_SUSTAINED_DROP);
			return
		}
		let idComponents = callingIdComponents(movedElement.id);
		let params = "name=" + movedElement.parentElement.id + "&params=" + idComponents.holderName + ":" + movedElement.getAttribute("data-org") + ":" + idComponents.callingName;
		apiCall("backout-transaction", params, refreshFromModel);
		clearCallingsHeldByMember();
		makeDirty()
		return
	}

	// dragging from ward tree to ward tree, cancel the operation
	if (dropTarget.classList.contains("nested") && movedElement.parentElement.classList.contains("nested")) {
		notify(nALERT, MESSAGE_RELEASE_ONLY)
		return
	}

	// dragging from member-callings list
	if (movedElement.classList.contains("member-calling")) {
		// only allow it to be moved to releases
		if (dropTarget.id !== "releases") {
			notify(nALERT, MESSAGE_RELEASE_ONLY);
			return
		}
		movedElement.remove();
	}

	// dragging a vacant calling from the tree
	if (movedElement.classList.contains("vacant")) {
		notify(nALERT, MESSAGE_RELEASE_VACANT);
		return
	}

	apiCall("remove-member-calling", createTransactionParmsFromTreeElememt(movedElement), refreshFromModel);
	makeDirty()
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
	return "org=" + element.getAttribute("data-org") + "&calling=" + callingIdParts[0] + "&member=" + callingIdParts[1];
}

function createTransactionParmsForMemberElement(memberElement, callingElement) {
	let callingIdParts = callingElement.id.split("@");
	return "org=" + callingElement.getAttribute("data-org") + "&calling=" + callingIdParts[0] + "&member=" + memberElement.id;
}

//// api call ////

function apiCall(endpoint, params, callback) {
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function () {
		if (this.readyState === 4) {
			if (this.status === 200) {
				callback(this.responseText);
			} else {
				callback(this.status);
			}
		}
	};

	xhttp.open("GET", "/v1/" + endpoint + "?" + encodeURI(params));
	xhttp.setRequestHeader("Content-type", "text/plain");
	xhttp.send();
}
