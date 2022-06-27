package engine

import (
	"encoding/json"
	"os"
	"sort"
	"time"
)

const (
	Adult = iota
	Youth
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
	return MemberSetDifference(this.GetMembers(18, 120), callings.MembersWithCallings())
}

func (this *Members) GetMemberRecord(name string) Member {
	if member, found := this.MemberMap[name]; found {
		member.Age = member.age()
		member.AgeByEndOfYear = member.ageByEndOfYear()
		return member
	}
	return Member{}
}

func (this *Members) GetMembers(minAge, maxAge int) (names []string) {
	for name, member := range this.MemberMap {
		if !member.Unbaptized && member.age() >= minAge && member.age() <= maxAge {
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

func (this *Members) GetMembersWithFocus() (focusMembers []MemberWithFocus) {
	members := this.GetMembers(18, 99)

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

///// private /////

func (this *Members) copy() Members {
	newMembers := NewMembers(this.initialSize, this.filePath)
	newMembers.initialSize = this.initialSize
	newMembers.filePath = this.filePath

	for name, member := range this.MemberMap {
		newMembers.MemberMap[name] = Member{
			Name:           member.Name,
			Gender:         member.Gender,
			Birthday:       member.Birthday,
			Unbaptized:     member.Unbaptized,
			Age:            0,
			AgeByEndOfYear: 0,
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
	Name            string
	CallingEligible uint8
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

func (this *Member) age() int {
	today := time.Now()
	birthdate := this.Birthday
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

func (this *Member) ageByEndOfYear() int {
	age := this.age()
	_, bm, bd := this.Birthday.Date()
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
