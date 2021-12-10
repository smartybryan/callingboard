package engine

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
	"time"

	"github.org/smartybryan/callingboard/util"
)

type Callings struct {
	CallingMap        map[string][]Calling
	OrganizationOrder []string
	mutex             *sync.Mutex

	initialSize int
	filePath    string
}

func NewCallings(numCallings int, path string) Callings {
	return Callings{
		CallingMap:  make(map[string][]Calling, numCallings),
		initialSize: numCallings,
		filePath:    path,
		mutex:       &sync.Mutex{},
	}
}

func (this *Callings) CallingList(organization string) (callingList []Calling) {
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
	return callingList
}

func (this *Callings) CallingListForMember(member string) (callingList []Calling) {
	allCallings := this.CallingList(ALL_ORGANIZATIONS)
	for _, calling := range allCallings {
		if calling.Holder == member {
			callingList = append(callingList, calling)
		}
	}
	return callingList
}

func (this *Callings) Count() int {
	return len(this.CallingList(ALL_ORGANIZATIONS))
}

func (this *Callings) MembersWithCallings() (names []string) {
	nameMap := make(map[string]struct{}, 200)
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

func (this *Callings) OrganizationList() (organizationList []string) {
	return this.OrganizationOrder
}

func (this *Callings) VacantCallingList(organization string) (callingList []Calling) {
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

func (this *Callings) Save() (numObjects int, err error) {
	if len(this.CallingMap) == 0 {
		return 0, nil
	}
	jsonBytes, err := json.Marshal(this)
	if err != nil {
		return 0, err
	}
	err = os.WriteFile(this.filePath, jsonBytes, 0660)
	return this.Count(), err
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

func (this *Callings) doesMemberHoldCalling(member string, org string, calling string) bool {
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

func (this *Callings) isValidOrganization(org string) bool {
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

///// model modification methods (called by Project) /////

func (this *Callings) addCalling(org string, calling string, custom bool) error {
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

func (this *Callings) removeCalling(org string, calling string) error {
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

func (this *Callings) updateCalling(org string, calling string, custom bool) error {
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

func (this *Callings) addMemberToACalling(member string, org string, calling string) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}

	if this.doesMemberHoldCalling(member, org, calling) {
		return ERROR_MEMBER_INVALID_CALLING
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
	return ERROR_INVALID_TRANSACTION
}

func (this *Callings) moveMemberToAnotherCalling(
	member string, fromOrg string, fromCalling string, toOrg string, toCalling string) error {
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

func (this *Callings) removeMemberFromACalling(member string, org string, calling string) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}
	if !this.doesMemberHoldCalling(member, org, calling) {
		return ERROR_MEMBER_INVALID_CALLING
	}

	callingList := this.CallingMap[org]
	for idx, call := range callingList {
		if call.Holder == member && call.Name == calling {
			call.Holder = VACANT_CALLING
			call.Sustained = time.Time{}
			callingList[idx] = call
			this.CallingMap[org] = callingList
			return nil
		}
	}
	return ERROR_INVALID_TRANSACTION
}

///////////////////////////////////////////////////////

type Calling struct {
	Org           string
	SubOrg        string
	Name          string
	Holder        string
	CustomCalling bool
	Sustained     time.Time

	PrintableSustained     string
	PrintableTimeInCalling string
}

const (
	VACANT_CALLING    = "Calling Vacant"
	ALL_ORGANIZATIONS = "All Organizations"
)

func CallingSetDifference(mainSet, subtractSet []Calling) (callings []Calling) {
	for _, calling := range mainSet {
		if callingInSet(subtractSet, calling) {
			continue
		}
		callings = append(callings, calling)
	}
	return callings
}

func callingInSet(set []Calling, calling Calling) bool {
	for _, call := range set {
		if call.Equal(calling) {
			return true
		}
	}
	return false
}

func (this *Calling) Equal(calling Calling) bool {
	return this.Org == calling.Org &&
		this.SubOrg == calling.SubOrg &&
		this.Name == calling.Name &&
		this.Holder == calling.Holder &&
		this.CustomCalling == calling.CustomCalling &&
		this.Sustained == calling.Sustained
}

func (this *Calling) copy() Calling {
	return Calling{
		Org:           this.Org,
		SubOrg:        this.SubOrg,
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

// OrganizationParseMap
// Used during parsing to determine a change of organization.
// Left values are the sub-organizations listed in the web page data,
// Right values are the organization. If empty, the left value is the
// organization with no sub-org
var OrganizationParseMap = map[string]string{
	"Bishopric":                            "",
	"Elders Quorum Presidency":             "Elders Quorum",
	"Relief Society Presidency":            "Relief Society",
	"Presidency of the Aaronic Priesthood": "Aaronic Priesthood Quorums",
	"Young Women Presidency":               "Young Women",
	"Sunday School Presidency":             "Sunday School",
	"Primary Presidency":                   "Primary",
	"Ward Missionaries":                    "",
	"Full-Time Missionaries":               "",
	"Temple and Family History":            "",
	"Young Single Adult":                   "Other Callings",
}
