package classes

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"
)

type MemberName string
type Members map[MemberName]Member

func NewMembers(numMembers int) Members {
	return make(map[MemberName]Member, numMembers)
}

func (this *Members) ParseMembersFromRawData(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	var fileLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileLines = append(fileLines, strings.TrimSpace(scanner.Text()))
	}

	withinMembers := false
	for idx := 0; idx < len(fileLines); idx++ {
		if !withinMembers {
			if strings.HasPrefix(fileLines[idx], "Name") {
				withinMembers = true
			}
			continue
		}

		if strings.HasPrefix(fileLines[idx], "Count") {
			return nil
		}

		memberRecord := strings.Split(fileLines[idx], "\t")
		if len(memberRecord) < 3 {
			continue
		}

		unbaptized :=  false
		if memberRecord[0][0] == '*' {
			unbaptized = true
			memberRecord[0] = memberRecord[0][1:]
		}
		birthday, _ := time.Parse("2 Jan 2006", memberRecord[3])
		member := Member{
			Name:     memberRecord[0],
			Gender:   memberRecord[1],
			Birthday: birthday,
			Unbaptized: unbaptized,
		}

		(*this)[MemberName(member.Name)] = member
	}

	return nil
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
	Name       string
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
