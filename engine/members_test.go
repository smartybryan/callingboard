package engine

import (
	"os"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestMembersFixture(t *testing.T) {
	gunit.Run(new(MembersFixture), t)
}

type MembersFixture struct {
	*gunit.Fixture
}

func (this *MembersFixture) TestAgeFunctions() {
	member := createMember("User, Joe", 0)

	member.Birthday = setDate(-20, 0, 0)
	this.So(member.age(), should.Equal, 20)

	member.Birthday = setDate(-20, 0, 2)
	this.So(member.age(), should.Equal, 19)

	member.Birthday = setDate(-20, 0, 2)
	this.So(member.ageByEndOfYear(), should.Equal, 20)
}

func (this *MembersFixture) TestGetMembers() {
	members := createTestMembers("")

	this.So(len(members.GetMembers(11, 17)), should.Equal, 1)
	this.So(len(members.GetMembers(18, 99)), should.Equal, 3)
	this.So(members.GetMembers(18, 99), should.Resemble,
		[]string{"Last2, First2", "Last3, First3", "Last4, First4"})
}

func (this *MembersFixture) TestFocusMembers() {
	members := createTestMembers("")

	data := []string{"Last2, First2", "Last4, First4"}
	expected := []MemberWithFocus{
		{Name: "Last2, First2", Focus: true},
		{Name: "Last3, First3", Focus: false},
		{Name: "Last4, First4", Focus: true},
	}
	_ = members.PutFocusMembers(data)
	this.So(members.GetMembersWithFocus(), should.Resemble, expected)
}

func (this *MembersFixture) TestAdultsWithoutACalling() {
	members := createTestMembers("")
	callings := createTestCallings("")

	this.So(members.AdultsWithoutACalling(callings), should.Resemble, []string{"Last3, First3", "Last4, First4"})
}

func (this *MembersFixture) TestSaveLoad() {
	tempFile := "testmembers"
	members := createTestMembers(tempFile)
	_ = members.PutFocusMembers([]string{"Focus, One;Focus, Two"})
	mLength := len(members.MemberMap)
	fLength := len(members.FocusMembers)
	_, err := members.Save()
	this.So(err, should.BeNil)

	members = NewMembers(10, tempFile)
	err = members.Load()
	this.So(err, should.BeNil)
	this.So(len(members.MemberMap), should.Equal, mLength)
	this.So(len(members.FocusMembers), should.Equal, fLength)
	_ = os.Remove(tempFile)
}

////////////////////////////////////////////////////////

func createTestMembers(path string) Members {
	members := NewMembers(5, path)
	members.MemberMap["Last1, First1"] = createMember("Last1, First1", 15)
	members.MemberMap["Last2, First2"] = createMember("Last2, First2", 20)
	members.MemberMap["Last3, First3"] = createMember("Last3, First3", 55)
	members.MemberMap["Last4, First4"] = createMember("Last4, First4", 30)
	return members
}

func createTestCallings(path string) Callings {
	callings := NewCallings(5, path)
	calling1 := createCalling("calling1", "Last1, First1", 2, 6)
	calling2 := createCalling("calling2", "Last2, First2", 1, 6)
	calling3 := createCalling("calling3", VACANT_CALLING, 0, 6)
	callings.CallingMap["org1"] = []Calling{calling1, calling2, calling3}
	callings.CallingMap["org2"] = []Calling{calling3}
	callings.OrganizationOrder = append(callings.OrganizationOrder, "org1")
	callings.OrganizationOrder = append(callings.OrganizationOrder, "org2")
	return callings
}

func createMember(name string, age int) Member {
	return Member{Name: string(name), Birthday: setDate(-age, 0, 0)}
}

func createCalling(name, memberName string, years, months int) Calling {
	return Calling{Name: name, Holder: string(memberName), Sustained: setDate(-years, -months, 0)}
}

func setDate(yearOffset, monthOffset, dayOffset int) time.Time {
	cy, cm, cd := time.Now().Date()
	return time.Date(cy+yearOffset, cm+time.Month(monthOffset), cd+dayOffset, 0, 0, 0, 0, time.UTC)
}
