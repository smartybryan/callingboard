package classes

import (
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
	this.So(member.Age(), should.Equal, 20)

	member.Birthday = setDate(-20, 0, 2)
	this.So(member.Age(), should.Equal, 19)

	member.Birthday = setDate(-20, 0, 2)
	this.So(member.AgeByEndOfYear(), should.Equal, 20)
}

func (this *MembersFixture) TestGetMembers() {
	members := NewMembers(5)
	members["Last1, First1"] = createMember("Last1, First1", 15)
	members["Last2, First2"] = createMember("Last2, First2", 20)
	members["Last3, First3"] = createMember("Last3, First3", 55)

	this.So(len(members.GetMembers(11, 17)), should.Equal, 1)
	this.So(len(members.GetMembers(18, 99)), should.Equal, 2)
	this.So(members.GetMembers(18, 99), should.Resemble, []MemberName{"Last2, First2", "Last3, First3"})
}

/*
func (this *Members) AdultsWithoutACallng(callings Callings) (members []Member) {
func (this *Members) AdultsEligibleForACalling() (members []Member) {
func (this *Members) YouthEligibleForACalling() (members []Member) {
*/

////////////////////////////////////////////////////////

func createMember(name string, age int) Member {
	return Member{Name: MemberName(name), Birthday: setDate(-age, 0,0)}
}

func createCalling(name, memberName string, years, months int) Calling {
	return Calling{Name: name, Holder: MemberName(memberName), Sustained: setDate(-years, -months,0)}
}

func setDate(yearOffset, monthOffset, dayOffset int) time.Time {
	cy, cm, cd := time.Now().Date()
	return time.Date(cy+yearOffset, cm+time.Month(monthOffset), cd+dayOffset, 0, 0, 0, 0, time.UTC)
}
