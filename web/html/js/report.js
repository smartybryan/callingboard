
function generateReport() {
	let canvas = document.getElementById("report-canvas");
	let releases = getCallingDataFromContainer("releases");
	let sustainings = getCallingDataFromContainer("sustainings");
	let content = buildReportTable(
		"We have released the following.",
		"Those who wish to express appreciation for their service, please raise your hand.",
		"Release", releases);
	content += "<br><br><br>";
	content += buildReportTable(
		"We have called the following and ask that you stand when your name is read, and remain standing until the vote is taken.",
		"We propose that they be sustained. All in favor, please raise your hand. [Wait] Those opposed by the same sign.",
		"Sustain", sustainings);
	canvas.innerHTML = content;
}

function getCallingDataFromContainer(id) {
	let callingList = [];
	let elem = document.getElementById(id);
	let callings = elem.getElementsByTagName("li");
	for (let call of callings) {
		if (call.id.length === 0) {
			continue;
		}
		let idInfo = callingIdComponents(call.id);
		let calling = idInfo.callingName;
		let holder = idInfo.holderName;
		let org = call.getAttribute("data-org");
		callingList.push({org, calling, holder});
	}
	return callingList;
}

function buildReportTable(pretable, posttable, title, callings) {
	let content = `
	<style>
	.report-table {
	  border: 1px solid black;
	  border-collapse: collapse;
	  padding: 5px;
	}
	.report-table th {
		background-color: aliceblue;
		width: 250px;
	}
	</style>
	<h2>` + title + `</h2>
	<div class="report-verbiage">` + pretable + `</div>
	<br>
	<table class="report-table">
	<thead>
	<tr><th class="report-table"><strong>Member</strong></th>
	<th class="report-table"><strong>Calling</strong></th>
	<th class="report-table"><strong>Organization</strong></th></tr>
	</thead>
	<tbody>`;

	for (let call of callings) {
		content += "<tr>";
		content += '<td class="report-table">' + call.holder + "</td>";
		content += '<td class="report-table">' + call.calling + "</td>";
		content += '<td class="report-table">' + call.org + "</td>";
		content += "</tr>";
	}
	content += `</tbody></table>
	<br>
	<div class="report-verbiage">` + posttable + `</div>`;
	return content
}
