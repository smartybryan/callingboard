const MESSAGE_RELEASE_ONLY = "You can only drag this calling into Releases.";
const MESSAGE_RELEASE_VACANT = "You cannot drag a vacant calling anywhere.";
const MESSAGE_RELEASE_SUSTAINED_DROP = "Drop this in the Trash to Undo.";
const MESSAGE_MEMBER_ONLY_TO_TREE = "You can only drag this member into a Vacant Calling.";
const MESSAGE_MODEL_LOADED = "Model loaded";
const MESSAGE_MODEL_RESET = "Original data has been reloaded.";
const MESSAGE_MODEL_SAVED = "Model has been saved.";
const MESSAGE_IMPORT_SUCCESSFUL = "The data has been imported successfully.";
const MESSAGE_IMPORT_FAILURE = "Import was not successful.";
const NOT_AUTHENTICATED = "Please login";

const nSUCCESS = "Success";
const nERROR = "Error";
const nALERT = "Alert";
const nINFO = "Info";

function notify(type, description) {
	const Notify = new XNotify();
	let duration = 4000;

	switch (type) {
		case nSUCCESS:
			Notify.success({
				title: type,
				description: description,
				duration: duration
			});
			break;
		case nERROR:
			Notify.error({
				title: type,
				description: description,
				duration: duration
			});
			break;
		case nALERT:
			Notify.alert({
				title: type,
				description: description,
				duration: duration
			});
			break;
		default:
			Notify.info({
				title: type,
				description: description,
				duration: duration
			});
			break;
	}
}
