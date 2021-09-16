package engine

import "fmt"

type Project struct {
	callings         Callings
	originalCallings Callings
	members          Members
	transactions     []Transaction
	undoHistory      []Transaction
}

func NewProject(callings Callings, members Members) *Project {
	return &Project{
		callings:         callings,
		originalCallings: callings.Copy(),
		members:          members,
		transactions:     make([]Transaction, 0, 100),
	}
}

func (this *Project) AddTransaction(operation string, parameters ...interface{}) {
	this.transactions = append(this.transactions, Transaction{
		Operation:  operation,
		Parameters: parameters,
	})
}

func (this *Project) PlayTransactions() {
	this.callings = this.originalCallings.Copy()

	for _, transaction := range this.transactions {
		fmt.Printf("Op:%s, Params:%+v\n", transaction.Operation, transaction.Parameters)
	}
}

func (this *Project) RedoTransaction() {
	if len(this.undoHistory) == 0 {
		return
	}
	this.transactions = append(this.transactions, this.undoHistory[len(this.undoHistory)-1])
	this.undoHistory = this.undoHistory[:len(this.undoHistory)-1]
}

func (this *Project) UndoTransaction() {
	if len(this.transactions) == 0 {
		return
	}
	this.undoHistory = append(this.undoHistory, this.transactions[len(this.transactions)-1])
	this.transactions = this.transactions[:len(this.transactions)-1]
}
