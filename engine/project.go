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

///// model modification stubs /////

func (this *Project) AddCalling(org Organization, calling string, custom bool) error {
	this.addTransaction("addCalling", org, calling, custom)
	return this.Callings.addCalling(org, calling, custom)
}

func (this *Project) RemoveCalling(org Organization, calling string) error {
	this.addTransaction("removeCalling", org, calling)
	return this.Callings.removeCalling(org, calling)
}

func (this *Project) UpdateCalling(org Organization, calling string, custom bool) error {
	this.addTransaction("updateCalling", org, calling, custom)
	return this.Callings.updateCalling(org, calling, custom)
}

func (this *Project) AddMemberToACalling(member MemberName, org Organization, calling string) error {
	this.addTransaction("addMemberToACalling", member, org, calling)
	return this.Callings.addMemberToACalling(member, org, calling)
}

func (this *Project) MoveMemberToAnotherCalling(
	member MemberName, fromOrg Organization, fromCalling string, toOrg Organization, toCalling string) error {
	this.addTransaction("moveMemberToAnotherCalling", member, fromOrg, fromCalling, toOrg, toCalling)
	return this.Callings.moveMemberToAnotherCalling(member, fromOrg, fromCalling, toOrg, toCalling)
}

func (this *Project) RemoveMemberFromACalling(member MemberName, org Organization, calling string) error {
	this.addTransaction("removeMemberFromACalling", member, org, calling)
	return this.Callings.removeMemberFromACalling(member, org, calling)
}

///// private /////

func (this *Project) addTransaction(operation string, parameters ...interface{}) {
	this.transactions = append(this.transactions, Transaction{
		Operation:  operation,
		Parameters: parameters,
	})
}

func (this *Project) playTransactions() {
	freshCallings := this.originalCallings.Copy()
	this.Callings = &freshCallings

	for _, transaction := range this.transactions {
		switch transaction.Operation {
		case "addCalling":
			_ = this.Callings.addCalling(
				transaction.Parameters[0].(Organization), transaction.Parameters[1].(string), transaction.Parameters[2].(bool))
		case "removeCalling":
			_ = this.Callings.removeCalling(
				transaction.Parameters[0].(Organization), transaction.Parameters[1].(string))
		case "updateCalling":
			_ = this.Callings.updateCalling(
				transaction.Parameters[0].(Organization), transaction.Parameters[1].(string), transaction.Parameters[2].(bool))
		case "addMemberToACalling":
			_ = this.Callings.addMemberToACalling(
				transaction.Parameters[0].(MemberName), transaction.Parameters[1].(Organization), transaction.Parameters[2].(string))
		case "moveMemberToAnotherCalling":
			_ = this.Callings.moveMemberToAnotherCalling(
				transaction.Parameters[0].(MemberName),
				transaction.Parameters[1].(Organization), transaction.Parameters[2].(string),
				transaction.Parameters[1].(Organization), transaction.Parameters[2].(string))
		case "removeMemberFromACalling":
			_ = this.Callings.removeMemberFromACalling(
				transaction.Parameters[0].(MemberName), transaction.Parameters[1].(Organization), transaction.Parameters[2].(string))
		}
	}
}
