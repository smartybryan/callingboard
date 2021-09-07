package engine

import (
	"encoding/json"
	"os"
	"sort"
	"time"
)

type MemberName string
type Members map[MemberName]Member

func NewMembers(numMembers int) Members {
	return make(map[MemberName]Member, numMembers)
}

func (this *Members) GetMembers(minAge, maxAge int) (names []MemberName) {
	for name, member := range *this {
		if !member.Unbaptized && member.Age() >= minAge && member.Age() <= maxAge {
			names = append(names, name)
		}
	}
	sort.SliceStable(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}

func (this *Members) AdultsWithoutACalling(callings Callings) (names []MemberName) {
	return SetDifference(this.GetMembers(18, 99), callings.MembersWithCallings())
}

func (this *Members) AdultsEligibleForACalling() (members []MemberName) {
	return this.GetMembers(18, 99)
}

func (this *Members) YouthEligibleForACalling() (members []MemberName) {
	return this.GetMembers(11, 17)
}

func (this *Members) Save(path string) error {
	jsonBytes, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonBytes, 0660)
}

func (this *Members) Load(path string) error {
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, this)
}

func (this *Members) SortedKeys() []MemberName {
	var names []MemberName

	for name, _ := range *this {
		names = append(names, name)
	}
	sort.SliceStable(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	return names
}

//////////////////////////////////////////////////////

type Member struct {
	Name       MemberName
	Gender     string
	Birthday   time.Time
	Unbaptized bool
}

func (this *Member) Age() int {
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

func (this *Member) AgeByEndOfYear() int {
	age := this.Age()
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
