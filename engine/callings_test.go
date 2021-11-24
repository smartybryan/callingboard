package engine

import (
	"os"
	"reflect"
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
	callings := NewCallings(5,"")
	callings.CallingMap["org1"] = []Calling{
		{Holder: "Washington, George"}, {Holder: "Lincoln, Abraham"}, {Holder: "Washington, George"},
	}

	this.So(callings.MembersWithCallings(), should.Resemble, []string{"Lincoln, Abraham", "Washington, George"})
}

func (this *CallingsFixture) TestOrganizationalList() {
	callings := createTestCallings("")
	this.So(callings.OrganizationList(), should.Resemble, []string{"org1", "org2"})
}

func (this *CallingsFixture) TestCallingList() {
	callings := createTestCallings("")
	this.So(len(callings.CallingList("org1")), should.Equal, 3)
}

func (this *CallingsFixture) TestCallingListForMember() {
	callings := createTestCallings("")
	callingsForMember := callings.CallingListForMember("Last1, First1")
	this.So(len(callingsForMember), should.Equal, 1)
}

func (this *CallingsFixture) TestVacantCallingList() {
	callings := createTestCallings("")
	this.So(len(callings.VacantCallingList("org1")), should.Equal, 1)
}

func (this *CallingsFixture) TestAddMemberToACalling() {
	callings := createTestCallings("")

	// invalid org
	err := callings.addMemberToACalling("Last99, First99", "bogusorg", "calling3")
	this.So(err, should.NotBeNil)

	// update an existing vacant calling
	err = callings.addMemberToACalling("Last99, First99", "org1", "calling3")
	this.So(err, should.BeNil)
	this.So(callings.doesMemberHoldCalling("Last99, First99", "org1", "calling3"), should.BeTrue)

	// create a new calling
	err = callings.addMemberToACalling("Last99, First99", "org1", "calling4")
	this.So(err, should.NotBeNil)
	this.So(callings.doesMemberHoldCalling("Last99, First99", "org1", "calling4"), should.BeFalse)
}

func (this *CallingsFixture) TestRemoveMemberFromACalling() {
	callings := createTestCallings("")

	// invalid org
	err := callings.removeMemberFromACalling("Last2, First2", "bogusorg", "calling3")
	this.So(err, should.NotBeNil)

	// member doesn't hold calling
	err = callings.removeMemberFromACalling("Last2, First2", "org1", "calling3")
	this.So(err, should.NotBeNil)

	// remove from calling happy path
	err = callings.removeMemberFromACalling("Last2, First2", "org1", "calling2")
	this.So(err, should.BeNil)
	this.So(callings.doesMemberHoldCalling("Last2, First2", "org1", "calling2"), should.BeFalse)
	this.So(callings.CallingMap["org1"][1].Holder, should.Equal, VACANT_CALLING)
}

func (this *CallingsFixture) TestMoveMemberToAnotherCalling() {
	callings := createTestCallings("")

	// invalid toOrg
	err := callings.moveMemberToAnotherCalling("Last2, First2", "bogusorg", "calling3", "org1", "calling4")
	this.So(err, should.NotBeNil)

	// invalid fromOrg
	err = callings.moveMemberToAnotherCalling("Last2, First2", "org1", "calling3", "bogusorg", "calling4")
	this.So(err, should.NotBeNil)

	// user doesn't hold fromCalling
	err = callings.moveMemberToAnotherCalling("Last2, First2", "org1", "calling4", "org1", "calling3")
	this.So(err, should.NotBeNil)

	// move calling happy path
	this.So(callings.doesMemberHoldCalling("Last2, First2", "org1", "calling2"), should.BeTrue)
	this.So(callings.doesMemberHoldCalling("Last2, First2", "org2", "calling3"), should.BeFalse)
	err = callings.moveMemberToAnotherCalling("Last2, First2", "org1", "calling2", "org2", "calling3")
	this.So(err, should.BeNil)
	this.So(callings.doesMemberHoldCalling("Last2, First2", "org1", "calling2"), should.BeFalse)
	this.So(callings.doesMemberHoldCalling("Last2, First2", "org2", "calling3"), should.BeTrue)
}

func (this *CallingsFixture) TestAddCalling() {
	callings := createTestCallings("")

	// invalid org
	err := callings.addCalling("bogusorg", "calling4", false)
	this.So(err, should.NotBeNil)

	// happy path
	err = callings.addCalling("org1", "calling4", false)
	this.So(err, should.BeNil)
	this.So(callings.CallingMap["org1"][3].Name, should.Equal, "calling4")
}

func (this *CallingsFixture) TestRemoveCalling() {
	callings := createTestCallings("")

	// invalid org
	err := callings.removeCalling("bogusorg", "calling4")
	this.So(err, should.NotBeNil)

	// invalid calling
	err = callings.removeCalling("org1", "calling4")
	this.So(err, should.NotBeNil)

	// happy path
	err = callings.removeCalling("org1", "calling3")
	this.So(err, should.BeNil)
	this.So(len(callings.CallingMap["org1"]), should.Equal, 2)
}

func (this *CallingsFixture) TestUpdateCalling() {
	callings := createTestCallings("")

	// invalid org
	err := callings.updateCalling("bogusorg", "calling4",true)
	this.So(err, should.NotBeNil)

	// invalid calling
	err = callings.updateCalling("org1", "calling4",true)
	this.So(err, should.NotBeNil)

	// happy path
	err = callings.updateCalling("org1", "calling3",true)
	this.So(err, should.BeNil)
	this.So(callings.CallingMap["org1"][2].CustomCalling, should.BeTrue)
}

func (this *CallingsFixture) TestCopy() {
	callings := createTestCallings("")
	callingsCopy := callings.copy()

	this.So(reflect.DeepEqual(callings, callingsCopy), should.BeTrue)

	err := callings.addCalling("org1", "calling5", false)
	this.So(err, should.BeNil)
	this.So(reflect.DeepEqual(callings, callingsCopy), should.BeFalse)
}

func (this *CallingsFixture) TestSaveLoad() {
	tempFile := "testcallings"
	callings := createTestCallings(tempFile)

	cmLength := len(callings.CallingMap)
	ooLength := len(callings.OrganizationOrder)
	_, err := callings.Save()
	this.So(err, should.BeNil)

	callings = NewCallings(10, tempFile)
	err = callings.Load()
	this.So(err, should.BeNil)
	this.So(len(callings.CallingMap), should.Equal, cmLength)
	this.So(len(callings.OrganizationOrder), should.Equal, ooLength)

	_ = os.Remove(tempFile)
}
