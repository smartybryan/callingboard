package engine

import (
	"os"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestProjectFixture(t *testing.T) {
	gunit.Run(new(ProjectFixture), t, gunit.Options.SequentialFixture()) //sequential for load/save/list tests
}

type ProjectFixture struct {
	*gunit.Fixture
}

func (this *ProjectFixture) TestUndoRedoTransactions() {
	callings := createTestCallings("")
	members := createTestMembers("")

	project := NewProject(&callings, &members, "")
	project.addTransaction("op1", "p1", "p2")
	project.addTransaction("op2", "p3", "p4", boolToString(true))
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

func (this *ProjectFixture) TestRemoveTransaction() {
	callings := createTestCallings("")
	members := createTestMembers("")

	project := NewProject(&callings, &members, "")
	project.addTransaction("removeMemberFromACalling", "p1", "p2", "p3")
	project.addTransaction("addMemberToACalling", "p3", "p4", "p5")
	this.So(len(project.transactions), should.Equal, 2)
}

func (this *ProjectFixture) TestPlayTransactions() {
	callings := createTestCallings("")
	members := createTestMembers("")
	project := NewProject(&callings, &members, "")

	_ = project.AddCalling("org1", "calling4", false)
	this.So(project.Callings.CallingList("org1")[3].Name, should.Equal, "calling4")
	this.So(project.Callings.CallingList("org1")[3].Holder, should.Equal, VACANT_CALLING)

	_ = project.AddMemberToACalling("Last3, First3","org1", "calling4")
	this.So(project.Callings.CallingList("org1")[3].Holder, should.Equal, "Last3, First3")

	project.UndoTransaction()
	this.So(project.Callings.CallingList("org1")[3].Holder, should.Equal, VACANT_CALLING)

	project.RedoTransaction()
	this.So(project.Callings.CallingList("org1")[3].Holder, should.Equal, "Last3, First3")

	this.So(len(project.transactions), should.Equal, 2)
	_ = project.removeTransaction("sustainings", []string{"Last3, First3", "org1", "calling4"})
	this.So(len(project.transactions), should.Equal, 1)

	_ = project.RemoveMemberFromACalling("Last3, First3","org1", "calling4")
	_ = project.removeTransaction("releases", []string{"Last3, First3", "org1", "calling4"})
	this.So(len(project.transactions), should.Equal, 1)
}

func (this *ProjectFixture) TestDiff() {
	callings := createTestCallings("")
	members := createTestMembers("")
	project := NewProject(&callings, &members, "")

	_ = project.AddCalling("org2", "calling10", true)
	_ = project.AddMemberToACalling("Last10, First10","org2", "calling10")
	_ = project.RemoveMemberFromACalling("Last1, First1","org1", "calling1")
	_ = project.RemoveCalling("org1", "calling1")
	_ = project.MoveMemberToAnotherCalling("Last2, First2","org1", "calling2", "org2", "calling22")

	diff := project.Diff()
	this.So(len(diff.Sustainings), should.Equal, 1)
	this.So(len(diff.Releases), should.Equal, 2)

	this.So(len(diff.NewVacancies), should.Equal, 1)
	this.So(diff.NewVacancies[0].Name, should.Equal, "calling2")
}

func (this *ProjectFixture) TestSaveLoadTransactions() {
	tempFile := "testtransactions"
	callings := createTestCallings("")
	members := createTestMembers("")
	project := NewProject(&callings, &members, "./")

	_ = project.AddCalling("org2", "calling10", true)
	_ = project.AddMemberToACalling("Last10, First10","org2", "calling10")
	_ = project.RemoveMemberFromACalling("Last1, First1","org1", "calling1")
	_ = project.RemoveCalling("org1", "calling1")
	_ = project.MoveMemberToAnotherCalling("Last2, First2","org1", "calling2", "org2", "calling22")

	tLastOp := project.transactions[len(project.transactions)-1].Operation
	tLastPLength := len(project.transactions[len(project.transactions)-1].Parameters)
	err := project.SaveTransactions(tempFile)
	this.So(err, should.BeNil)

	project.transactions = project.transactions[:]
	err = project.LoadTransactions(tempFile, true)
	this.So(err, should.BeNil)
	this.So(len(project.transactions), should.Equal, 3)
	this.So(project.transactions[len(project.transactions)-1].Operation, should.Equal, tLastOp)
	this.So(len(project.transactions[len(project.transactions)-1].Parameters), should.Equal, tLastPLength)

	os.Remove(tempFile+TransactionFileSuffix)
}

func (this *ProjectFixture) TestLoadTransactionsNoOverwrite() {
	tempFile, tempFile1 := "testtransactions1", "testtransactions2"
	callings := createTestCallings("")
	members := createTestMembers("")
	project := NewProject(&callings, &members, "./")

	_ = project.AddCalling("org2", "calling10", true)
	_ = project.AddCalling("org2", "calling11", true)
	_ = project.AddMemberToACalling("Last10, First10","org2", "calling10")
	_ = project.RemoveMemberFromACalling("Last1, First1","org1", "calling1")
	_ = project.RemoveCalling("org1", "calling1")
	_ = project.MoveMemberToAnotherCalling("Last2, First2","org1", "calling2", "org2", "calling22")
	err := project.SaveTransactions(tempFile)
	this.So(err, should.BeNil)

	project.transactions = project.transactions[:0]
	_ = project.AddMemberToACalling("Last11, First11","org2", "calling11")
	tLastOp := project.transactions[len(project.transactions)-1].Operation
	tLastPLength := len(project.transactions[len(project.transactions)-1].Parameters)
	err = project.SaveTransactions(tempFile1)
	this.So(err, should.BeNil)

	err = project.LoadTransactions(tempFile, true)
	this.So(err, should.BeNil)
	err = project.LoadTransactions(tempFile1, false)
	this.So(err, should.BeNil)
	this.So(len(project.transactions), should.Equal, 4)
	this.So(project.transactions[len(project.transactions)-1].Operation, should.Equal, tLastOp)
	this.So(len(project.transactions[len(project.transactions)-1].Parameters), should.Equal, tLastPLength)

	os.Remove(tempFile+TransactionFileSuffix)
	os.Remove(tempFile1+TransactionFileSuffix)
}

func (this *ProjectFixture) TestListTransactionFiles() {
	tempFile1, tempFile2 := "testtrans1", "testtrans2"
	callings := createTestCallings("")
	members := createTestMembers("")
	project := NewProject(&callings, &members, "./")

	project.SaveTransactions(tempFile1)
	project.SaveTransactions(tempFile2)
	files := project.ListTransactionFiles()
	this.So(len(files), should.BeGreaterThanOrEqualTo, 2)

	os.Remove(tempFile1+TransactionFileSuffix)
	os.Remove(tempFile2+TransactionFileSuffix)
}
