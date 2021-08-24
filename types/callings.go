package types

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type Organization string

type Callings struct {
	CallingMap        map[Organization][]Calling
	OrganizationOrder []Organization
}

type Calling struct {
	Name          string
	Holder        string
	CustomCalling bool
	Sustained     time.Time
}

func NewCallings(numCallings int) Callings {
	return Callings{CallingMap: make(map[Organization][]Calling, numCallings)}
}

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
		if strings.HasPrefix(fileLines[idx], "Position\tName") {
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
			calling.Name = calling.Name[2:]
			calling.CustomCalling = true
		}

		if fileLines[idx+1] == "Calling Vacant" {
			calling.Holder = fileLines[idx+1]
			idx++
		} else {
			calling.Holder = fileLines[idx+1]
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

func getOrganizationPrefixFromCalling(callingName string) Organization {
	for _, organization := range SharedCallingOrganizations {
		if strings.HasPrefix(callingName, string(organization)) {
			return organization
		}
	}
	return ""
}

var MultiUseOrganizations = map[Organization]struct{}{
	"Activities":          struct{}{},
	"Ministering":         struct{}{},
	"Music":               struct{}{},
	"Service":             struct{}{},
	"Unassigned Teachers": struct{}{},
}

var SharedCallingOrganizations = []Organization{
	"Elders Quorum",
	"Priesthood",
	"Primary",
	"Relief Society",
	"Sunday School",
}

//TODO: fix music - see how music is currrently broken down and determine why Priesthood music is way down at the bottoms