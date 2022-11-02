function generateReport() {
	let dateOptions = {
		weekday: 'long',
		year: 'numeric',
		month: 'long',
		day: 'numeric',
		hour: "numeric",
		minute: "numeric"
	};
	let today = new Date();
	let dateTime = today.toLocaleDateString("en-US", dateOptions);
	let canvas = document.getElementById("report-canvas");
	let reportType = document.getElementById("reports").value;

	let content = "<div>Printed " + dateTime + " from Callingboard.org</div>"

	switch (reportType) {
		case "relsus":
			generateCallingReleasesReport(canvas, content);
			break;
		case "memrep":
			generateMemberCallingReport(canvas, content);
			break;
	}
}

//////////////////////////////////////////////////////

function generateCallingReleasesReport(canvas, content) {
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

	canvas.innerHTML = content
}

function getCallingDataFromContainer(id) {
	let callingList = [];
	let elem = document.getElementById(id);
	let callings = elem.getElementsByTagName("li");
	for (let call of callings) {
		if (call.id.length === 0) {
			continue;
		}
		let calling = call.getAttribute("data-callname");
		let holder = call.getAttribute("data-holder")
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

//////////////////////////////////////////////////////////////////////////////

function generateMemberCallingReport(canvas, content) {
	apiCall("members")
		.then(data => {
			getCallingData(data, content, canvas, "Member Callings");
		})
		.catch(error => {
			console.log(error);
			if (error === 401 || error === NOT_LOGGED_IN) {
				logout();
			}
		});

	return content
}

function getCallingData(memberData, content, canvas, title) {
	let endpoint = "callings";
	let params = "org=" + ALL_ORGS;

	apiCall(endpoint, params)
		.then(data => {
			buildMemberReportTable(data, memberData, content, canvas, title);
		})
		.catch(error => {
			console.log(error);
		});
}

function parseCallingData(callingJson) {
	const callingMap = new Map();

	callingJson.forEach(function (calling) {
		let callings = callingMap.get(calling.Holder);
		if (!callings) {
			callings = [];
		}
		callings.push({
			"name": calling.Name,
			"org": calling.Org,
			"suborg": calling.SubOrg,
			"time": calling.PrintableTimeInCalling
		});
		callingMap.set(calling.Holder, callings);
	});
	return callingMap
}

function buildMemberReportTable(callingResponse, memberResponse, content, canvas, title) {
	let wardId = getAuthValueFromCookie().wardid;

	content += `
	<h2>` + title + `</h2>
	<span>
		<span class="filter">Show:</span>
		<input onclick="redisplayMemberReport()" type="checkbox" id="member-report-male" checked>M</span>
		<input onclick="redisplayMemberReport()" type="checkbox" id="member-report-female" checked>F</span>
	</span>
	<br><br>
	
	<table class="report-table">
	<thead>
	<tr>
	<th class="report-table"><strong>Name</strong></th>
	<th class="report-table" style="width: 200px;"><strong>Photo</strong></th>
	<th class="report-table" style="width: 600px;"><strong>Calling</strong></th>
	</tr>
	</thead>
	<tbody>`;

	let callingMap = parseCallingData(JSON.parse(callingResponse));
	let memberJson = JSON.parse(memberResponse);
	memberJson.forEach(function (member) {
		let memberParts = member.Name.split(";")
		let memberName = memberParts[0];
		let memberImage = wardId + "/" + encodeURI(memberName) + ".jpg";

		content += "<tr class='member-report-row " + findGender(memberParts[1]) + "'>";
		content += '<td class="report-table">' + memberName + "</td>";
		content += '<td class="report-table" style="text-align:center;">' + '<img class="thumbnail-report" onload="this.style.display=\'\'" style="display: none;" src="' + memberImage + '?v=' + imageVersion + '" alt="">' + "</td>";
		content += '<td class="report-table">' + getMemberCallings(callingMap, memberName) + "</td>";
		content += "</tr>";
	});

	content += `</tbody></table>`
	canvas.innerHTML = content;
}

function getMemberCallings(callingMap, memberName) {
	let callings = callingMap.get(memberName);
	if (!callings) {
		return "";
	}

	let callingInfo = "";
	callings.forEach(function (calling) {
		if (callingInfo.length > 0) {
			callingInfo += "<br>"
		}
		callingInfo += "<div>" + calling.name + "</div>";
		callingInfo += "<div>(" + calling.org;
		if (calling.suborg.length > 0) {
			callingInfo += ", " + calling.suborg;
		}
		callingInfo += ")</div>"
		callingInfo += "<div>" + calling.time + "</div>"
	});

	return callingInfo;
}

function findGender(genderIdentifier) {
	switch (genderIdentifier) {
		case "1": return "male-report";
		case "2": return "female-report";
	}
	return "";
}

function redisplayMemberReport() {
	let includeFemale = document.getElementById("member-report-female").checked;
	let includeMale = document.getElementById("member-report-male").checked;

	let reportRows = document.getElementsByClassName("member-report-row");
	for (let row of reportRows) {
		if ((row.classList.contains("male-report") && includeMale) || (row.classList.contains("female-report") && includeFemale)) {
			row.classList.remove("filtered");
		} else {
			row.classList.add("filtered");
		}
	}
}
