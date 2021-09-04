package classes

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Organization string

type Callings struct {
	CallingMap        map[Organization][]Calling
	OrganizationOrder []Organization
}

func NewCallings(numCallings int) Callings {
	return Callings{CallingMap: make(map[Organization][]Calling, numCallings)}
}

func (this *Callings) AddMemberToACalling(member MemberName, org Organization, calling string) error {
	if !this.isValidOrganization(org) {
		return errors.New(fmt.Sprintf("Calling organization %s does not exist", org))
	}

	if this.doesMemberHoldCalling(member, org, calling) {
		return nil
	}

	callingList := this.CallingMap[org]
	for idx, call := range callingList {
		if call.Name == calling && call.Holder == VACANT_CALLING {
			call.Holder = member
			callingList[idx] = call
			this.CallingMap[org] = callingList
			return nil
		}
	}
	newCalling := Calling{
		Name:          calling,
		Holder:        member,
		CustomCalling: false,
		Sustained:     time.Time{},
	}
	callingList = append(callingList, newCalling)
	this.CallingMap[org] = callingList
	return nil
}

func (this *Callings) CallingList(organization Organization) (callingList []Calling) {
	if callings, found := this.CallingMap[organization]; found {
		for _, calling := range callings {
			callingList = append(callingList, calling)
		}
	}
	sort.SliceStable(callingList, func(i, j int) bool {
		return callingList[i].Name < callingList[j].Name
	})
	return callingList
}

func (this *Callings) MembersWithCallings() (names []MemberName) {
	nameMap := make(map[MemberName]struct{}, 200)
	for _, callings := range (*this).CallingMap {
		for _, calling := range callings {
			nameMap[calling.Holder] = struct{}{}
		}
	}

	for name, _ := range nameMap {
		names = append(names, name)
	}
	sort.SliceStable(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}

func (this *Callings) MoveMemberToAnotherCalling(
	member MemberName, fromOrg Organization, fromCalling string, toOrg Organization, toCalling string) error {

	return nil
}

func (this *Callings) OrganizationList() (organizationList []Organization) {
	for organization, _ := range this.CallingMap {
		organizationList = append(organizationList, organization)
	}
	sort.SliceStable(organizationList, func(i, j int) bool {
		return organizationList[i] < organizationList[j]
	})
	return organizationList
}

func (this *Callings) RemoveMemberFromACalling(member MemberName, org Organization, calling string) error {

	return nil
}

func (this *Callings) VacantCallingList(organization Organization) (callingList []Calling) {
	allCallings := this.CallingList(organization)
	for _, calling := range allCallings {
		if calling.Holder == VACANT_CALLING {
			callingList = append(callingList, calling)
		}
	}
	return callingList
}

func (this *Callings) Save(path string) error {
	jsonBytes, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonBytes, 0660)
}

func (this *Callings) Load(path string) error {
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, this)
}

func (this *Callings) doesMemberHoldCalling(member MemberName, org Organization, calling string) bool {
	if !this.isValidOrganization(org) {
		return false
	}
	callingList := this.CallingMap[org]
	for _, call := range callingList {
		if call.Name == calling && call.Holder == member {
			return true
		}
	}
	return false
}

func (this *Callings) isValidOrganization(org Organization) bool {
	_, found := this.CallingMap[org]
	return found
}

func getOrganizationPrefixFromCalling(callingName string) Organization {
	for _, organization := range SharedCallingOrganizations {
		if strings.HasPrefix(callingName, string(organization)) {
			return organization
		}
	}
	return ""
}

///////////////////////////////////////////////////////

type Calling struct {
	Name          string
	Holder        MemberName
	CustomCalling bool
	Sustained     time.Time
}

const (
	VACANT_CALLING = "Calling Vacant"
)

func (this *Calling) DaysInCalling() int {
	return int(time.Now().Sub(this.Sustained).Hours() / 24)
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
	"Primary",
	"Relief Society",
	"Sunday School",
}
