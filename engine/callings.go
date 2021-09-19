package engine

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"
)

type Organization string

type Callings struct {
	callingMap        map[Organization][]Calling
	organizationOrder []Organization
}

func NewCallings(numCallings int) Callings {
	return Callings{
		callingMap: make(map[Organization][]Calling, numCallings),
	}
}

func (this *Callings) CallingList(organization Organization) (callingList []Calling) {
	if callings, found := this.callingMap[organization]; found {
		for _, calling := range callings {
			callingList = append(callingList, calling)
		}
	}
	sort.SliceStable(callingList, func(i, j int) bool {
		return callingList[i].Name < callingList[j].Name
	})
	return callingList
}

func (this *Callings) Copy() Callings {
	newCallings := NewCallings(len(this.callingMap) * 2)

	for organization, callings := range this.callingMap {
		newCallings.callingMap[organization] = callings
	}

	for _, organization := range this.organizationOrder {
		newCallings.organizationOrder = append(newCallings.organizationOrder, organization)
	}

	return newCallings
}

func (this *Callings) MembersWithCallings() (names []MemberName) {
	nameMap := make(map[MemberName]struct{}, 200)
	for _, callings := range (*this).callingMap {
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

func (this *Callings) OrganizationList() (organizationList []Organization) {
	for organization, _ := range this.callingMap {
		organizationList = append(organizationList, organization)
	}
	sort.SliceStable(organizationList, func(i, j int) bool {
		return organizationList[i] < organizationList[j]
	})
	return organizationList
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

func (this *Callings) Load(path string) error {
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, this)
}

func (this *Callings) Save(path string) error {
	jsonBytes, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonBytes, 0660)
}

///// private /////

func (this *Callings) doesMemberHoldCalling(member MemberName, org Organization, calling string) bool {
	if !this.isValidOrganization(org) {
		return false
	}
	callingList := this.callingMap[org]
	for _, call := range callingList {
		if call.Name == calling && call.Holder == member {
			return true
		}
	}
	return false
}

func (this *Callings) isValidOrganization(org Organization) bool {
	_, found := this.callingMap[org]
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

///// model modification methods (called by Project) /////

func (this *Callings) addCalling(org Organization, calling string, custom bool) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}

	callingList := this.callingMap[org]
	newCalling := Calling{
		Name:          calling,
		Holder:        VACANT_CALLING,
		CustomCalling: custom,
		Sustained:     time.Time{},
	}
	callingList = append(callingList, newCalling)
	this.callingMap[org] = callingList

	return nil
}

func (this *Callings) removeCalling(org Organization, calling string) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}
	callingList := this.callingMap[org]
	var newCallingList []Calling
	for _, call := range callingList {
		if call.Name != calling {
			newCallingList = append(newCallingList, call)
		}
	}
	if len(callingList) == len(newCallingList) {
		return ERROR_UNKNOWN_CALLING
	}
	this.callingMap[org] = newCallingList

	return nil
}

func (this *Callings) updateCalling(org Organization, calling string, custom bool) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}
	callingList := this.callingMap[org]
	for idx, call := range callingList {
		if call.Name == calling {
			call.CustomCalling = custom
			this.callingMap[org][idx] = call
			return nil
		}
	}

	return ERROR_UNKNOWN_CALLING
}

func (this *Callings) addMemberToACalling(member MemberName, org Organization, calling string) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}

	if this.doesMemberHoldCalling(member, org, calling) {
		return nil
	}

	callingList := this.callingMap[org]
	for idx, call := range callingList {
		if call.Name == calling && call.Holder == VACANT_CALLING {
			call.Holder = member
			callingList[idx] = call
			this.callingMap[org] = callingList
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
	this.callingMap[org] = callingList

	return nil
}

func (this *Callings) moveMemberToAnotherCalling(
	member MemberName, fromOrg Organization, fromCalling string, toOrg Organization, toCalling string) error {
	err := this.removeMemberFromACalling(member, fromOrg, fromCalling)
	if err != nil {
		return err
	}

	err = this.addMemberToACalling(member, toOrg, toCalling)
	if err != nil {
		return err
	}

	return nil
}

func (this *Callings) removeMemberFromACalling(member MemberName, org Organization, calling string) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}
	if this.doesMemberHoldCalling(member, org, calling) {
		callingList := this.callingMap[org]
		for idx, calling := range callingList {
			if calling.Holder == member {
				calling.Holder = VACANT_CALLING
				calling.Sustained = time.Time{}
				callingList[idx] = calling
				this.callingMap[org] = callingList
				return nil
			}
		}
	}
	return ERROR_MEMBER_INVALID_CALLING
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
