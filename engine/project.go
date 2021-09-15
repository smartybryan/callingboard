package engine

import "fmt"

type Project struct {
	Callings     Callings
	Members      Members
	Transactions []Transaction
	UndoHistory  []Transaction
}

func NewProject(callings Callings, members Members) *Project {
	return &Project{
		Callings:     callings,
		Members:      members,
		Transactions: make([]Transaction, 0, 100),
	}
}

func (this *Project) AddTransaction(operation string, parameters ...interface{}) {
	this.Transactions = append(this.Transactions, Transaction{
		Operation:  operation,
		Parameters: parameters,
	})
}

func (this *Project) PlayTransactions() {
	for _, transaction := range this.Transactions {
		fmt.Printf("Op:%s, Params:%+v\n", transaction.Operation, transaction.Parameters)
	}
}

func (this *Project) UndoTransaction() {
	// TODO check boundaries
	this.UndoHistory = append(this.UndoHistory, this.Transactions[len(this.Transactions)-1])
	this.Transactions = this.Transactions[:len(this.Transactions)-1]
}

func (this *Project) RedoTransaction() {
	// TODO check boundaries
	this.Transactions = append(this.Transactions, this.UndoHistory[len(this.UndoHistory)-1])
	this.UndoHistory = this.UndoHistory[:len(this.UndoHistory)-1]
}
