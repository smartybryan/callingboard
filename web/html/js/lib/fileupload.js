// https://css-tricks.com/examples/DragAndDropFileUploading/
'use strict';

let maxImageFileSizeMB = 5;
let maxImageFileSize = maxImageFileSizeMB*1024*1024;

function initImageForms(document) {
	// applying the effect for every form
	let forms = document.querySelectorAll('.box');
	Array.prototype.forEach.call(forms, function (form) {
		let input = form.querySelector('input[type="file"]'),
			label = form.querySelector('label'),
			errorMsg = form.querySelector('.box__error span'),
			restart = form.querySelectorAll('.box__restart'),
			droppedFiles = false,
			triggerFormSubmit = function () {
				let event = document.createEvent('HTMLEvents');
				event.initEvent('submit', true, false);
				form.dispatchEvent(event);
			};

		// letting the server side to know we are going to make an Ajax request
		let ajaxFlag = document.createElement('input');
		ajaxFlag.setAttribute('type', 'hidden');
		ajaxFlag.setAttribute('name', 'ajax');
		ajaxFlag.setAttribute('value', 1);
		form.appendChild(ajaxFlag);

		// automatically submit the form on file select
		input.addEventListener('change', function (e) {
			triggerFormSubmit();
		});

		['dragover', 'dragenter', 'dragleave', 'drop'].forEach(function (event) {
			form.addEventListener(event, function (e) {
				// preventing the unwanted behaviours
				e.preventDefault();
				e.stopPropagation();
			});
		});
		['dragover', 'dragenter'].forEach(function (event) {
			form.addEventListener(event, function () {
				form.classList.add('is-dragover');
			});
		});
		['dragleave', 'dragend', 'drop'].forEach(function (event) {
			form.addEventListener(event, function () {
				form.classList.remove('is-dragover');
			});
		});
		form.addEventListener('drop', function (e) {
			droppedFiles = e.dataTransfer.files; // the files that were dropped
			triggerFormSubmit();
		});

		// if the form was submitted
		form.addEventListener('submit', function (e) {
			// preventing the duplicate submissions if the current one is in progress
			if (form.classList.contains('is-uploading')) return false;

			form.classList.add('is-uploading');
			form.classList.remove('is-error');

			e.preventDefault();

			// gathering the form data
			let fileName = "";
			let fileSize = 0;
			let ajaxData = new FormData(form);
			if (droppedFiles) {
				Array.prototype.forEach.call(droppedFiles, function (file) {
					ajaxData.append(input.getAttribute('name'), file);
					fileName = file.name.toLowerCase();
					fileSize = file.size;
				});
			}

			if (fileSize > maxImageFileSize ||
				(!fileName.endsWith("jpg") && !fileName.endsWith("jpeg") && !fileName.endsWith("png"))) {
				notify(nERROR, "You can only upload image files (jpg, png) with a max size of "+maxImageFileSizeMB+"MB.");
				form.classList.remove('is-uploading');
				return false;
			}

			// ajax request
			let ajax = new XMLHttpRequest();
			ajax.open(form.getAttribute('method'), form.getAttribute('action') + "&file=" + encodeURI(fileName), true);

			ajax.onload = function () {
				form.classList.remove('is-uploading');
				if (ajax.status === 200) {
					let data = ajax.responseText;
					let success = data === "success";
					form.classList.add(success ? 'is-success' : 'is-error');
					if (success) {
						notify(nSUCCESS, "File uploaded");
						populateMemberPhotoList();
						imageVersion++;
						displayMembers("");
					} else {
						errorMsg.textContent = data;
					}
				} else {
					notify(nERROR, ajax.responseText);
				}
			};

			ajax.onerror = function () {
				form.classList.remove('is-uploading');
				notify(nERROR, "Something went wrong: " + ajax.responseText);
			};

			ajax.send(ajaxData);
		});

		// restart the form if has a state of error/success
		Array.prototype.forEach.call(restart, function (entry) {
			entry.addEventListener('click', function (e) {
				e.preventDefault();
				form.classList.remove('is-error', 'is-success');
				input.click();
			});
		});

		// Firefox focus bug fix for file input
		input.addEventListener('focus', function () {
			input.classList.add('has-focus');
		});
		input.addEventListener('blur', function () {
			input.classList.remove('has-focus');
		});

	});
}
