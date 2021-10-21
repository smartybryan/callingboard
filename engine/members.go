package engine

import (
	"encoding/json"
	"os"
	"sort"
	"time"
)

type MemberName string

type Members struct {
	MemberMap map[MemberName]Member

	initialSize int
	filePath    string
}

func NewMembers(numMembers int, path string) Members {
	return Members{
		MemberMap:   make(map[MemberName]Member, numMembers),
		initialSize: numMembers,
		filePath:    path,
	}
}

func (this *Members) AdultsWithoutACalling(callings Callings) (names []MemberName) {
	return SetDifference(this.GetMembers(18, 120), callings.MembersWithCallings())
}

func (this *Members) AdultsEligibleForACalling() (members []MemberName) {
	return this.GetMembers(18, 99)
}

func (this *Members) GetMemberRecord(name MemberName) Member {
	if member, found := this.MemberMap[name]; found {
		member.Age = member.age()
		member.AgeByEndOfYear = member.ageByEndOfYear()
		return member
	}
	return Member{}
}

func (this *Members) GetMembers(minAge, maxAge int) (names []MemberName) {
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

func (this *Members) YouthEligibleForACalling() (members []MemberName) {
	return this.GetMembers(11, 17)
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

//////////////////////////////////////////////////////

type Member struct {
	Name       MemberName
	Gender     string
	Birthday   time.Time
	Unbaptized bool

	Age            int
	AgeByEndOfYear int
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
