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

func (this *ProjectFixture) Setup() {
}

func (this *ProjectFixture) TestUndoRedoTransactions() {
	callings := createTestCallings()
	members := createTestMembers()

	project := NewProject(callings, members)
	project.AddTransaction("op1", "p1", "p2")
	project.AddTransaction("op2", "p3", "p4", true)
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
	
	project.PlayTransactions()
}
