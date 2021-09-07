package engine

import (
	"bufio"
	"os"
	"strings"
	"time"
)

func (this *Callings) ParseCallingsFromRawData(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	var fileLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileLines = append(fileLines, strings.TrimSpace(scanner.Text()))
	}

	var currentOrganization Organization
	withinOrganization := false

	for idx := 0; idx < len(fileLines); idx++ {
		if strings.HasPrefix(fileLines[idx], "Position") {
			currentOrganization = Organization(fileLines[idx-1])
			this.OrganizationOrder = append(this.OrganizationOrder, currentOrganization)
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
			calling.Holder = MemberName(fileLines[idx+1])
			idx++
		} else {
			calling.Holder = MemberName(fileLines[idx+1])
			sustained, err := time.Parse("2 Jan 2006", fileLines[idx+2])
			if err == nil {
				calling.Sustained = sustained
			}
			idx += 2
		}

		// for organization names that are shared by multiple organizations,
		// prepend the actual organization to make them specific
		if _, found := MultiUseOrganizations[currentOrganization]; found {
			if prefix := getOrganizationPrefixFromCalling(calling.Name); len(prefix) > 0 {
				currentOrganization = prefix + " " + currentOrganization
				this.OrganizationOrder[len(this.OrganizationOrder)-1] = currentOrganization
			}
		}
		(*this).CallingMap[currentOrganization] = append((*this).CallingMap[currentOrganization], calling)
	}

	return nil
}


func (this *Members) ParseMembersFromRawData(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	var fileLines []string
	scanner := bufio.NewScanner(file)
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
			return nil
		}

		memberRecord := strings.Split(fileLines[idx], "\t")
		if len(memberRecord) < 3 {
			continue
		}

		unbaptized :=  false
		if memberRecord[0][0] == '*' {
			unbaptized = true
			memberRecord[0] = memberRecord[0][1:]
		}
		birthday, _ := time.Parse("2 Jan 2006", memberRecord[3])
		member := Member{
			Name:     MemberName(memberRecord[0]),
			Gender:   memberRecord[1],
			Birthday: birthday,
			Unbaptized: unbaptized,
		}

		(*this)[MemberName(member.Name)] = member
	}

	return nil
}