
function generateReport() {
	let dateOptions = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric', hour: "numeric", minute: "numeric" };
	let today  = new Date();
	let dateTime = today.toLocaleDateString("en-US", dateOptions);
	let canvas = document.getElementById("report-canvas");
	let reportType = document.getElementById("reports").value;

	let content = "<div>Printed " + dateTime + " from Callingboard.org</div>"

	if (reportType === "relsus") {
		let releases = getCallingDataFromContainer("releases");
		let sustainings = getCallingDataFromContainer("sustainings");
		content += buildReleaseReportTable(
			"We have released the following.",
			"Those who wish to express appreciation for their service, please raise your hand.",
			"Release", releases);
		content += "<br><br>";
		content += buildSustainReportTable(
			"We have called the following and ask that you stand when your name is read, and remain standing until the vote is taken.",
			"We propose that they be sustained. All in favor, please raise your hand. [Wait] Those opposed by the same sign.",
			"Sustain", sustainings);
	}

	if (reportType === "memrep") {
		content += "<br><br><h2>Report under construction.</h2>"
	}

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

function buildReleaseReportTable(pretable, posttable, title, callings) {
	let content = `
	<h2>` + title + `</h2>
	<div class="report-verbiage">` + pretable + `</div>
	<br>
	<table class="report-table">
	<thead>
	<tr><th class="report-table"><strong>Member</strong></th>
	<th class="report-table"><strong>Calling</strong></th>
	<th class="report-table"  style="width: 150px;"><strong>Organization</strong></th>
	<th class="report-table" style="width: 100px;"><strong>Released</strong></th>
	<th class="report-table"><strong>Date</strong></th></tr>
	</thead>
	<tbody>`;

	for (let call of callings) {
		content += "<tr>";
		content += '<td class="report-table">' + call.holder + "</td>";
		content += '<td class="report-table">' + call.calling + "</td>";
		content += '<td class="report-table">' + call.org + "</td>";
		content += '<td class="report-table" style="text-align: center;">☐</td>';
		content += '<td class="report-table"></td>';
		content += "</tr>";
	}
	content += `</tbody></table>
	<br>
	<div class="report-verbiage">` + posttable + `</div>`;
	return content
}

function buildSustainReportTable(pretable, posttable, title, callings) {
	let content = `
	<h2>` + title + `</h2>
	<div class="report-verbiage">` + pretable + `</div>
	<br>
	<table class="report-table">
	<thead>
	<tr><th class="report-table"><strong>Member</strong></th>
	<th class="report-table"><strong>Calling</strong></th>
	<th class="report-table"  style="width: 150px;"><strong>Organization</strong></th>
	<th class="report-table" style="width: 100px;"><strong>Sustained</strong></th>
	<th class="report-table"><strong>Set Apart By/Date</strong></th></tr>
	</thead>
	<tbody>`;

	for (let call of callings) {
		content += "<tr>";
		content += '<td class="report-table">' + call.holder + "</td>";
		content += '<td class="report-table">' + call.calling + "</td>";
		content += '<td class="report-table">' + call.org + "</td>";
		content += '<td class="report-table" style="text-align: center;">☐</td>';
		content += '<td class="report-table"></td>';
		content += "</tr>";
	}
	content += `</tbody></table>
	<br>
	<div class="report-verbiage">` + posttable + `</div>`;
	return content
}
