package engine

type Project struct {
	Callings         *Callings
	Members          *Members
	originalCallings Callings
	transactions     []Transaction
	undoHistory      []Transaction
}

func NewProject(callings *Callings, members *Members) *Project {
	return &Project{
		Callings:         callings,
		originalCallings: callings.Copy(),
		Members:          members,
		transactions:     make([]Transaction, 0, 100),
	}
}

func (this *Project) AddTransaction(operation string, parameters ...interface{}) {
	this.transactions = append(this.transactions, Transaction{
		Operation:  operation,
		Parameters: parameters,
	})
}

func (this *Project) RedoTransaction() {
	if len(this.undoHistory) == 0 {
		return
	}
	this.transactions = append(this.transactions, this.undoHistory[len(this.undoHistory)-1])
	this.undoHistory = this.undoHistory[:len(this.undoHistory)-1]
	this.playTransactions()
}

func (this *Project) UndoTransaction() {
	if len(this.transactions) == 0 {
		return
	}
	this.undoHistory = append(this.undoHistory, this.transactions[len(this.transactions)-1])
	this.transactions = this.transactions[:len(this.transactions)-1]
	this.playTransactions()
}

func (this *Project) playTransactions() {
	freshCallings := this.originalCallings.Copy()
	this.Callings = &freshCallings

	for _, transaction := range this.transactions {
		switch transaction.Operation {
		case "AddCalling":
			_ = this.Callings.AddCalling(
				transaction.Parameters[0].(Organization), transaction.Parameters[1].(string), transaction.Parameters[2].(bool))
		case "AddMemberToACalling":
			_ = this.Callings.AddMemberToACalling(
				transaction.Parameters[0].(MemberName), transaction.Parameters[1].(Organization), transaction.Parameters[2].(string))
		}
	}
}
