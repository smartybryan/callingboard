
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
		notify(nINFO,"Please select a model name to load.");
		return
	}
	modelOperation("load-trans", name)
	makeClean()
}

function deleteModel() {
	let name = getSelectedModelFile();
	if (!name) {
		notify(nINFO,"Please select a model name to delete.");
		return
	}
	let currentModel = document.getElementById("model-name").value;
	if (name === currentModel) {
		document.getElementById("model-name").value = "";
		resetModel();
	}
	//TODO: this is a race condition as resetmodel could still be executing
	modelOperation("delete-trans", name)
}

function saveModel() {
	let name = document.getElementById("model-name").value;
	if (!name) {
		notify(nINFO,"Please provide a model name to save.");
		return
	}
	modelOperation("save-trans", name)
	makeClean()
}

function resetModel() {
	modelOperation("reset-model");
	document.getElementById("model-name").value = "";
	notify(nSUCCESS, MESSAGE_MODEL_RESET);
	makeClean()
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
					notify(nSUCCESS,MESSAGE_MODEL_LOADED);
				}
				if (endpoint === "save-trans") {
					notify(nSUCCESS,MESSAGE_MODEL_SAVED);
				}
			} else {
				notify(nERROR, this.responseText)
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
