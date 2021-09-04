package classes

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestCallingsFixture(t *testing.T) {
	gunit.Run(new(CallingsFixture), t)
}

type CallingsFixture struct {
	*gunit.Fixture

	callings Callings
}

func (this *CallingsFixture) TestDaysInCalling() {
	cy, cm, cd := time.Now().Date()
	calling := Calling{
		Name:          "Head Honcho",
		Holder:        "User, Joe",
		CustomCalling: false,
	}

	calling.Sustained = time.Date(cy-1, cm, cd, 0, 0, 0, 0, time.UTC)
	this.So(calling.DaysInCalling(), should.BeBetween, 363, 368)

	calling.Sustained = time.Date(cy-1, cm, cd-5, 0, 0, 0, 0, time.UTC)
	this.So(calling.DaysInCalling(), should.BeBetween, 368, 373)
}

func (this *CallingsFixture) TestMembersWithCallings() {
	callings := NewCallings(5)
	callings.CallingMap["org1"] = []Calling{
		{Holder: "Washington, George"}, {Holder: "Lincoln, Abraham"}, {Holder: "Washington, George"},
	}

	this.So(callings.MembersWithCallings(), should.Resemble, []MemberName{"Lincoln, Abraham", "Washington, George"})
}

func (this *CallingsFixture) TestOrganizationalList() {
	callings := createTestCallings()
	this.So(callings.OrganizationList(), should.Resemble, []Organization{"org1", "org2"})
}

func (this *CallingsFixture) TestCallingList() {
	callings := createTestCallings()
	this.So(len(callings.CallingList("org1")), should.Equal, 3)
}

func (this *CallingsFixture) TestVacantCallingList() {
	callings := createTestCallings()
	this.So(len(callings.VacantCallingList("org1")), should.Equal, 1)
}

func (this *CallingsFixture) TestAddMemberToACalling() {
	callings := createTestCallings()

	err := callings.AddMemberToACalling("Last99, First99", "org99", "bogusCallling")
	this.So(err, should.NotBeNil)

	err = callings.AddMemberToACalling("Last99, First99", "org1", "calling3")
	this.So(err, should.BeNil)
	this.So(callings.doesMemberHoldCalling("Last99, First99", "org1", "calling3"), should.BeTrue)

	err = callings.AddMemberToACalling("Last99, First99", "org1", "calling4")
	this.So(err, should.BeNil)
	this.So(callings.doesMemberHoldCalling("Last99, First99", "org1", "calling4"), should.BeTrue)
}

/*

func (this *Callings) MoveMemberToAnotherCalling(member MemberName, fromOrg Organization, fromCalling string, toOrg Organization, toCalling string) error {
func (this *Callings) RemoveMemberFromACalling(member MemberName, org Organization, calling string) error {

*/
