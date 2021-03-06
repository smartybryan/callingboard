package engine

import (
	"bufio"
	"bytes"
	"strings"
	"time"
)

func (this *Callings) ParseCallingsFromRawData(data []byte) (callingCount int) {
	saveCallings := this.copy()
	*this = NewCallings(this.initialSize, this.filePath)

	var fileLines []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		fileLines = append(fileLines, strings.TrimSpace(scanner.Text()))
	}

	var currentOrganization, currentSubOrganization string
	withinOrganization := false

	for idx := 0; idx < len(fileLines); idx++ {
		if strings.HasPrefix(fileLines[idx], "Position") {
			currentSubOrganization = string(fileLines[idx-1])
			if org, found := OrganizationParseMap[currentSubOrganization]; found {
				currentOrganization = org
				if currentOrganization == "" {
					currentOrganization = currentSubOrganization
					currentSubOrganization = ""
				}
				this.OrganizationOrder = append(this.OrganizationOrder, currentOrganization)
			}
			idx++
			withinOrganization = true
		}

		if strings.HasPrefix(fileLines[idx], "Count") ||
			strings.HasPrefix(fileLines[idx], "Add Another") ||
			strings.HasPrefix(fileLines[idx], "No matching") ||
			strings.HasPrefix(fileLines[idx], "* customs") {
			withinOrganization = false
			continue
		}

		if !withinOrganization {
			continue
		}

		if fileLines[idx] == "" {
			continue
		}

		calling := Calling{Name: fileLines[idx]}
		if strings.HasPrefix(calling.Name, "*") {
			calling.Name = strings.TrimSpace(calling.Name[1:])
			calling.CustomCalling = true
		}

		if fileLines[idx+1] == "Calling Vacant" {
			calling.Holder = string(fileLines[idx+1])
			idx++
		} else {
			calling.Holder = string(fileLines[idx+1])
			sustained, err := time.Parse("2 Jan 2006", fileLines[idx+2])
			if err == nil {
				calling.Sustained = sustained
			}
			idx += 2
		}

		calling.Org = currentOrganization
		calling.SubOrg = currentSubOrganization
		(*this).CallingMap[currentOrganization] = append((*this).CallingMap[currentOrganization], calling)
		callingCount++
	}

	// if parse issue, keep current contents
	if callingCount == 0 {
		*this = saveCallings.copy()
	}

	return callingCount
}

func (this *Members) ParseMembersFromRawData(data []byte) int {
	saveMembers := this.copy()
	*this = NewMembers(this.initialSize, this.filePath)

	var fileLines []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		fileLines = append(fileLines, strings.TrimSpace(scanner.Text()))
	}

	withinMembers := false
	for idx := 0; idx < len(fileLines); idx++ {
		if !withinMembers {
			if strings.HasPrefix(fileLines[idx], "Name") {
				withinMembers = true
			}
			continue
		}

		if strings.HasPrefix(fileLines[idx], "Count") {
			return len(this.MemberMap)
		}

		memberRecord := strings.Split(fileLines[idx], "\t")
		if len(memberRecord) < 3 {
			continue
		}

		baptized := true
		if memberRecord[0][0] == '*' {
			baptized = false
			memberRecord[0] = memberRecord[0][1:]
		}

		member := Member{
			Name:        memberRecord[0],
			Eligibility: calculateEligibility(memberRecord[3], baptized),
			Type:        calculateType(memberRecord[1]),
		}
		this.MemberMap[member.Name] = member
	}

	// if parse issue, restore current members
	memberCount := len(this.MemberMap)
	if len(this.MemberMap) == 0 {
		*this = saveMembers.copy()
	}

	return memberCount
}

func calculateType(gender string) uint8 {
	switch gender {
	case "M":
		return 1
	case "F":
		return 2
	default:
		return 0
	}
}

func calculateEligibility(birthdate string, baptized bool) uint8 {
	birthday, _ := time.Parse("2 Jan 2006", birthdate)

	if !baptized {
		return CallingNotEligible
	}

	age := ageByEndOfYear(birthday)

	if age >= 18 {
		return CallingAdult
	}

	if age >= 12 {
		return CallingYouth
	}

	return CallingNotEligible
}

func ageByEndOfYear(birthdate time.Time) int {
	age := calculateAge(birthdate)
	_, bm, bd := birthdate.Date()
	_, tm, td := time.Now().Date()

	if bm < tm {
		return age
	}

	if bm > tm {
		return age + 1
	}

	if bd <= td {
		return age
	}

	return age + 1
}

func calculateAge(birthdate time.Time) int {
	today := time.Now()
	today = today.In(birthdate.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := birthdate.Date()
	birthdate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(birthdate) {
		return 0
	}
	age := ty - by
	anniversary := birthdate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}
	return age
}
