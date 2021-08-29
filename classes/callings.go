package classes

import (
	"encoding/json"
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

