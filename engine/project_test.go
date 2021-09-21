package engine

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestProjectFixture(t *testing.T) {
	gunit.Run(new(ProjectFixture), t)
}

type ProjectFixture struct {
	*gunit.Fixture
}

func (this *ProjectFixture) TestUndoRedoTransactions() {
	callings := createTestCallings()
	members := createTestMembers()

	project := NewProject(&callings, &members)
	project.addTransaction("op1", "p1", "p2")
	project.addTransaction("op2", "p3", "p4", true)
	this.So(len(project.transactions), should.Equal, 2)

	project.UndoTransaction()
	this.So(len(project.transactions), should.Equal, 1)
	this.So(project.transactions[0].Operation, should.Equal, "op1")
	this.So(len(project.undoHistory), should.Equal, 1)

	project.UndoTransaction()
	this.So(len(project.transactions), should.Equal, 0)
	this.So(len(project.undoHistory), should.Equal, 2)

	project.UndoTransaction()
	this.So(len(project.transactions), should.Equal, 0)

	project.RedoTransaction()
	this.So(len(project.transactions), should.Equal, 1)
	this.So(project.transactions[0].Operation, should.Equal, "op1")

	project.RedoTransaction()
	this.So(len(project.transactions), should.Equal, 2)
	this.So(len(project.undoHistory), should.Equal, 0)

	project.RedoTransaction()
	this.So(len(project.transactions), should.Equal, 2)
}

func (this *ProjectFixture) TestPlayTransactions() {
	callings := createTestCallings()
	members := createTestMembers()
	project := NewProject(&callings, &members)

	_ = project.AddCalling("org1", "calling4", false)
	this.So(project.Callings.CallingList("org1")[3].Name, should.Equal, "calling4")
	this.So(project.Callings.CallingList("org1")[3].Holder, should.Equal, VACANT_CALLING)

	_ = project.AddMemberToACalling("Last3, First3","org1", "calling4")
	this.So(project.Callings.CallingList("org1")[3].Holder, should.Equal, "Last3, First3")

	project.UndoTransaction()
	this.So(project.Callings.CallingList("org1")[3].Holder, should.Equal, VACANT_CALLING)

	project.RedoTransaction()
	this.So(project.Callings.CallingList("org1")[3].Holder, should.Equal, "Last3, First3")
}

func (this *ProjectFixture) TestDiff() {
	callings := createTestCallings()
	members := createTestMembers()
	project := NewProject(&callings, &members)

	_ = project.AddCalling("org2", "calling10", true)
	_ = project.AddMemberToACalling("Last10, First10","org2", "calling10")
	_ = project.RemoveMemberFromACalling("Last1, First1","org1", "calling1")
	_ = project.RemoveCalling("org1", "calling1")
	_ = project.MoveMemberToAnotherCalling("Last2, First2","org1", "calling2", "org2", "calling22")

	releases, sustainings := project.Diff()
	this.So(len(sustainings), should.Equal, 2)
	this.So(len(releases), should.Equal, 2)
}
