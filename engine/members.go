package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

const (
	CallingNotEligible = iota
	CallingYouth
	CallingAdult
	AllEligible
	AllMembers
)

type Members struct {
	MemberMap    map[string]Member
	FocusMembers []string

	initialSize int
	filePath    string
}

type MemberWithFocus struct {
	Name  string
	Focus bool
}

func NewMembers(numMembers int, path string) Members {
	return Members{
		MemberMap:   make(map[string]Member, numMembers),
		initialSize: numMembers,
		filePath:    path,
	}
}

func (this *Members) AdultsWithoutACalling(callings Callings) (names []string) {
	return MemberSetDifference(this.GetMembers(CallingAdult), callings.MembersWithCallings())
}

func (this *Members) GetMemberRecord(name string) Member {
	if member, found := this.MemberMap[name]; found {
		return member
	}
	return Member{}
}

func (this *Members) GetMembers(eligibility uint8) (names []string) {
	for name, member := range this.MemberMap {
		if eligibility == AllMembers ||
			(eligibility == AllEligible && (member.Eligibility == CallingYouth || member.Eligibility == CallingAdult)) ||
			(eligibility == CallingYouth && member.Eligibility == CallingYouth) ||
			(eligibility == CallingAdult && member.Eligibility == CallingAdult) {

			names = append(names, name)
		}
	}
	sort.SliceStable(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	return names
}

func (this *Members) Load() error {
	jsonBytes, err := os.ReadFile(this.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, this)
}

func (this *Members) Save() (numObjects int, err error) {
	if len(this.MemberMap) == 0 {
		return 0, nil
	}
	jsonBytes, err := json.Marshal(this)
	if err != nil {
		return 0, err
	}
	err = os.WriteFile(this.filePath, jsonBytes, 0660)
	return len(this.MemberMap), err
}

func (this *Members) GetMembersWithType(memberNames []string) (names []string) {
	for _, name := range memberNames {
		names = append(names, this.GetMemberRecord(name).BuildMemberName())
	}
	return names
}

func (this *Members) GetMembersWithFocus() (focusMembers []MemberWithFocus) {
	members := this.GetMembers(CallingAdult)

	for _, member := range members {
		focusMembers = append(focusMembers, MemberWithFocus{
			Name:  member,
			Focus: this.isMemberFocused(member),
		})
	}
	return focusMembers
}

func (this *Members) GetFocusMembers() []string {
	return this.FocusMembers
}

func (this *Members) PutFocusMembers(names []string) error {
	this.FocusMembers = names
	_, err := this.Save()
	return err
}

func (this *Members) UploadMemberImage(path, originalFile string, image []byte) error {
	// TODO: convert the original file to jpg
	return ioutil.WriteFile(path, image, os.FileMode(0666))
}

func (this *Members) DeleteMemberImage(path string) error {
	return os.Remove(path)
}

///// private /////

func (this *Members) copy() Members {
	newMembers := NewMembers(this.initialSize, this.filePath)
	newMembers.initialSize = this.initialSize
	newMembers.filePath = this.filePath

	for name, member := range this.MemberMap {
		newMembers.MemberMap[name] = Member{
			Name:        member.Name,
			Eligibility: member.Eligibility,
		}
	}

	return newMembers
}

func (this *Members) isMemberFocused(member string) bool {
	for _, focusedMember := range this.FocusMembers {
		if focusedMember == member {
			return true
		}
	}
	return false
}

//////////////////////////////////////////////////////

type Member struct {
	Name        string
	Eligibility uint8
	Type        uint8
}

func (this Member) BuildMemberName() string {
	return fmt.Sprintf("%s;%d", this.Name, this.Type)
}

func MemberSetDifference(mainSet, subtractSet []string) (names []string) {
	for _, name := range mainSet {
		if memberInSet(subtractSet, name) {
			continue
		}
		if !memberInSet(names, name) {
			names = append(names, name)
		}
	}
	return names
}

func memberInSet(set []string, name string) bool {
	for _, value := range set {
		if value == name {
			return true
		}
	}
	return false
}
