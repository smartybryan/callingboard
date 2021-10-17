package engine

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"

	"github.org/smartybryan/callorg/util"
)

type Organization string

type Callings struct {
	CallingMap        map[Organization][]Calling
	OrganizationOrder []Organization

	initialSize int
	filePath    string
}

func NewCallings(numCallings int, path string) Callings {
	return Callings{
		CallingMap:  make(map[Organization][]Calling, numCallings),
		initialSize: numCallings,
		filePath:    path,
	}
}

func (this *Callings) CallingList(organization Organization) (callingList []Calling) {
	if organization == ALL_ORGANIZATIONS {
		for _, org := range this.OrganizationOrder {
			if callings, found := this.CallingMap[org]; found {
				callingList = this.getCallingListByOrganization(callings, callingList)
			}
		}
	} else {
		if callings, found := this.CallingMap[organization]; found {
			callingList = this.getCallingListByOrganization(callings, callingList)
		}
	}
	sort.SliceStable(callingList, func(i, j int) bool {
		return callingList[i].Name < callingList[j].Name
	})
	return callingList
}

func (this *Callings) CallingListForMember(member MemberName) (callingList []Calling) {
	allCallings := this.CallingList(ALL_ORGANIZATIONS)
	for _, calling := range allCallings {
		if calling.Holder == member {
			callingList = append(callingList, calling)
		}
	}
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

func (this *Callings) OrganizationList() (organizationList []Organization) {
	for organization, _ := range this.CallingMap {
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

func (this *Callings) Load() error {
	jsonBytes, err := os.ReadFile(this.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, this)
}

func (this *Callings) Save() error {
	if len(this.CallingMap) == 0 {
		return nil
	}
	jsonBytes, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return os.WriteFile(this.filePath, jsonBytes, 0660)
}

///// private /////

func (this *Callings) copy() Callings {
	newCallings := NewCallings(len(this.CallingMap)*2, "")
	newCallings.initialSize = this.initialSize
	newCallings.filePath = this.filePath

	for organization, callings := range this.CallingMap {
		newCallings.CallingMap[organization] = []Calling{}
		for _, calling := range callings {
			newCallings.CallingMap[organization] = append(newCallings.CallingMap[organization], calling.copy())
		}
	}

	for _, organization := range this.OrganizationOrder {
		newCallings.OrganizationOrder = append(newCallings.OrganizationOrder, organization)
	}

	return newCallings
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

func (this *Callings) getCallingListByOrganization(callings []Calling, callingList []Calling) []Calling {
	for _, calling := range callings {
		calling.PrintableSustained = util.PrintableDate(calling.Sustained)
		calling.PrintableTimeInCalling = util.PrintableTimeInCalling(calling.DaysInCalling())
		callingList = append(callingList, calling)
	}
	return callingList
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

	callingList := this.CallingMap[org]
	newCalling := Calling{
		Name:          calling,
		Holder:        VACANT_CALLING,
		CustomCalling: custom,
		Sustained:     time.Time{},
	}
	callingList = append(callingList, newCalling)
	this.CallingMap[org] = callingList

	return nil
}

func (this *Callings) removeCalling(org Organization, calling string) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}
	callingList := this.CallingMap[org]
	var newCallingList []Calling
	for _, call := range callingList {
		if call.Name != calling {
			newCallingList = append(newCallingList, call)
		}
	}
	if len(callingList) == len(newCallingList) {
		return ERROR_UNKNOWN_CALLING
	}
	this.CallingMap[org] = newCallingList

	return nil
}

func (this *Callings) updateCalling(org Organization, calling string, custom bool) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}
	callingList := this.CallingMap[org]
	for idx, call := range callingList {
		if call.Name == calling {
			call.CustomCalling = custom
			this.CallingMap[org][idx] = call
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
		callingList := this.CallingMap[org]
		for idx, calling := range callingList {
			if calling.Holder == member {
				calling.Holder = VACANT_CALLING
				calling.Sustained = time.Time{}
				callingList[idx] = calling
				this.CallingMap[org] = callingList
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

	PrintableSustained     string
	PrintableTimeInCalling string
}

const (
	VACANT_CALLING = "Calling Vacant"
	ALL_ORGANIZATIONS = "All Organizations"
)

func (this *Calling) copy() Calling {
	return Calling{
		Name:          this.Name,
		Holder:        this.Holder,
		CustomCalling: this.CustomCalling,
		Sustained:     this.Sustained,
	}
}

func (this *Calling) DaysInCalling() int {
	if this.Sustained.IsZero() {
		return 0
	}
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
