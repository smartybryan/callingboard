package engine

import (
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.org/smartybryan/callingboard/util"
)

var CallingIdCounter map[string]int

type Callings struct {
	CallingMap        map[string][]Calling
	FocusCallings     map[string]struct{}
	OrganizationOrder []string

	initialSize int
	filePath    string
}

func NewCallings(numCallings int, path string) Callings {
	ResetCallingIdCounter()

	return Callings{
		CallingMap:    make(map[string][]Calling, numCallings),
		FocusCallings: make(map[string]struct{}, numCallings),
		initialSize:   numCallings,
		filePath:      path,
	}
}

func ResetCallingIdCounter() {
	CallingIdCounter = make(map[string]int, 30)
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
	return this.setFocusOnList(callingList)
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

func (this *Callings) SetCallingFocus(callingId string, focus bool) {
	if focus {
		this.insertFocusCalling(callingId)
	} else {
		this.removeFocusCalling(callingId)
	}
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
	return this.setFocusOnList(callingList)
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

func (this *Callings) doesMemberHoldCalling(member string, org string, suborg string, calling string) bool {
	if !this.isValidOrganization(org) {
		return false
	}
	callingList := this.CallingMap[org]
	for _, call := range callingList {
		if call.Name == calling && call.SubOrg == suborg && call.Holder == member {
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

func (this *Callings) addMemberToACalling(member string, org string, suborg string, callingName string) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}

	if this.doesMemberHoldCalling(member, org, suborg, callingName) {
		return ERROR_MEMBER_HAS_CALLING
	}

	callingList := this.CallingMap[org]
	for idx, call := range callingList {
		if call.Name == callingName && call.SubOrg == suborg && call.Holder == VACANT_CALLING {
			call.Holder = member
			callingList[idx] = call
			this.CallingMap[org] = callingList
			return nil
		}
	}
	return ERROR_INVALID_TRANSACTION
}

func (this *Callings) moveMemberToAnotherCalling(
	member string, fromOrg string, fromSubOrg string, fromCalling string, toOrg string, toSubOrg string, toCalling string) error {
	err := this.removeMemberFromACalling(member, fromOrg, fromSubOrg, fromCalling)
	if err != nil {
		return err
	}

	err = this.addMemberToACalling(member, toOrg, toSubOrg, toCalling)
	if err != nil {
		return err
	}

	return nil
}

func (this *Callings) removeMemberFromACalling(member string, org string, suborg string, callingName string) error {
	if !this.isValidOrganization(org) {
		return ERROR_UNKNOWN_ORGANIZATION
	}
	if !this.doesMemberHoldCalling(member, org, suborg, callingName) {
		return ERROR_MEMBER_INVALID_CALLING
	}

	callingList := this.CallingMap[org]
	for idx, call := range callingList {
		if call.Holder == member && call.SubOrg == suborg && call.Name == callingName {
			call.Holder = VACANT_CALLING
			call.Sustained = time.Time{}
			callingList[idx] = call
			this.CallingMap[org] = callingList

			return nil
		}
	}
	return ERROR_INVALID_TRANSACTION
}

func (this *Callings) isCallingFocused(calling Calling) bool {
	_, found := this.FocusCallings[calling.FocusKey()]
	return found
}

func (this *Callings) insertFocusCalling(callingId string) {
	this.FocusCallings[callingId] = struct{}{}
}

func (this *Callings) removeFocusCalling(callingId string) {
	delete(this.FocusCallings, callingId)
}

func (this *Callings) setFocusOnList(callingList []Calling) []Calling {
	for i, calling := range callingList {
		if _, found := this.FocusCallings[calling.Id]; found {
			callingList[i].Focus = true
		}
	}
	return callingList
}

///////////////////////////////////////////////////////

type Calling struct {
	Id            string
	Org           string
	SubOrg        string
	Name          string
	Holder        string
	CustomCalling bool
	Focus         bool
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
		Id:            this.Id,
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

func (this *Calling) GenerateId() (id string) {
	id = getInitials(this.Org) + getInitials(this.SubOrg) + getInitials(this.Name) + strconv.Itoa(CallingIdCounter[this.Org])
	CallingIdCounter[this.Org] = CallingIdCounter[this.Org] + 1
	return id
}
func getInitials(data string) (initials string) {
	if len(data) == 0 {
		return initials
	}

	if strings.Index(data, " ") > -1 {
		parts := strings.Split(data, " ")
		initials = parts[0][:1] + parts[1][:1]
	} else {
		initials = data[:1] + data[len(data)-1:]
	}
	return strings.ToUpper(initials)
}

func (this *Calling) FocusKey() string {
	return this.Id
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

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
