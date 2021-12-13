
function listModels() {
	modelOperation("list-trans-files")
		.then(data => {})
		.catch(error => {
			console.log(error);
		})
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
		.then(data => {
			makeClean();
		})
		.catch(error => {
			console.log(error);
		})
}

function mergeModel() {
	let name = getSelectedModelFile();
	if (!name) {
		notify(nINFO,"Please select a model name to merge with current model.");
		return
	}
	modelOperation("merge-trans", name)
		.then(data => {
			makeClean();
		})
		.catch(error => {
			console.log(error);
		})
}

function deleteModel() {
	let name = getSelectedModelFile();
	if (!name) {
		notify(nINFO,"Please select a model name to delete.");
		return
	}
	let currentModel = document.getElementById("model-name").value;

	modelOperation("delete-trans", name)
		.then(data => {
			if (name === currentModel) {
				document.getElementById("model-name").value = "";
				resetModel();
			}
		})
		.catch(error => {
			console.log(error);
		})
}

function saveModel() {
	let name = document.getElementById("model-name").value;
	if (!name) {
		notify(nINFO,"Please provide a model name to save.");
		return
	}
	modelOperation("save-trans", name)
		.then(data => {
			makeClean();
		})
		.catch(error => {
			console.log(error);
		})
}

function resetModel() {
	modelOperation("reset-model")
		.then(data => {
			document.getElementById("model-name").value = "";
			notify(nSUCCESS, MESSAGE_MODEL_RESET);
			makeClean();
		})
		.catch(error => {
			console.log(error);
		})
}

let modelOperation = (endpoint, name) => {
	let params = ""
	if (name) {
		params = "?name=" + name;
	}

	return new Promise((resolve, reject) => {
		const xhttp = new XMLHttpRequest();
		xhttp.open("GET", "/v1/" + endpoint + encodeURI(params));
		xhttp.setRequestHeader("Content-type", "text/plain");
		xhttp.onload = () => {
			if (xhttp.status >= 200 && xhttp.status < 300) {
				if (endpoint === "list-trans-files") {
					updateFileList(xhttp.responseText);
				} else if (xhttp.responseText === "null\n") {
					refreshFromModel();
					listModels();
					if (endpoint === "load-trans") {
						notify(nSUCCESS,MESSAGE_MODEL_LOADED);
					}
					if (endpoint === "save-trans") {
						notify(nSUCCESS,MESSAGE_MODEL_SAVED);
					}
				} else {
					notify(nERROR, xhttp.responseText)
				}
				resolve(xhttp.responseText);
			} else {
				reject(xhttp.statusText);
			}
		}
		xhttp.onerror = () => {
			reject(xhttp.statusText);
		}
		xhttp.send();
	})
}
