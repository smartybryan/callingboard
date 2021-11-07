package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	TransactionFileSuffix = ".txf"
)

type Project struct {
	Callings         *Callings
	Members          *Members
	originalCallings Callings
	transactions     []Transaction
	undoHistory      []Transaction
	dataPath         string

	diff DiffResult
}

type DiffResult struct {
	Sustainings []Calling
	Releases    []Calling
}

func NewProject(callings *Callings, members *Members, dataPath string) *Project {
	return &Project{
		Callings:         callings,
		originalCallings: callings.copy(),
		Members:          members,
		transactions:     make([]Transaction, 0, 100),
		dataPath:         dataPath,
		diff:             NewDiff(),
	}
}

func NewDiff() DiffResult {
	return DiffResult{
		Sustainings: make([]Calling, 0, 20),
		Releases:    make([]Calling, 0, 20),
	}
}

func (this *Project) Diff() DiffResult {
	this.diff.Releases = this.diff.Releases[:0]
	this.diff.Sustainings = this.diff.Sustainings[:0]

	for _, organization := range this.Callings.OrganizationOrder {
		// sustainings
		modelCallings := this.Callings.CallingList(organization)
		for _, modelCalling := range modelCallings {
			if modelCalling.Holder == VACANT_CALLING {
				continue
			}
			if this.originalCallings.doesMemberHoldCalling(modelCalling.Holder, organization, modelCalling.Name) {
				continue
			}
			this.diff.Sustainings = append(this.diff.Sustainings, modelCalling)
		}

		// releases
		originalCallings := this.originalCallings.CallingList(organization)
		for _, originalCalling := range originalCallings {
			if originalCalling.Holder == VACANT_CALLING {
				continue
			}
			if this.Callings.doesMemberHoldCalling(originalCalling.Holder, organization, originalCalling.Name) {
				continue
			}
			this.diff.Releases = append(this.diff.Releases, originalCalling)
		}
	}

	sort.SliceStable(this.diff.Releases, func(i, j int) bool {
		return this.diff.Releases[i].Name < this.diff.Releases[j].Name
	})
	sort.SliceStable(this.diff.Sustainings, func(i, j int) bool {
		return this.diff.Sustainings[i].Name < this.diff.Sustainings[j].Name
	})

	return this.diff
}

func (this *Project) RedoTransaction() bool {
	if len(this.undoHistory) == 0 {
		return false
	}
	this.transactions = append(this.transactions, this.undoHistory[len(this.undoHistory)-1])
	this.undoHistory = this.undoHistory[:len(this.undoHistory)-1]
	this.playTransactions()
	return true
}

func (this *Project) UndoTransaction() bool {
	if len(this.transactions) == 0 {
		return false
	}
	this.undoHistory = append(this.undoHistory, this.transactions[len(this.transactions)-1])
	this.transactions = this.transactions[:len(this.transactions)-1]
	this.playTransactions()

	return true
}

func (this *Project) ListTransactionFiles() (transactionFiles []string) {
	files, _ := ioutil.ReadDir(this.dataPath)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), TransactionFileSuffix) {
			transactionFiles = append(transactionFiles, strings.TrimSuffix(filepath.Base(file.Name()), TransactionFileSuffix))
		}
	}

	return transactionFiles
}

func (this *Project) LoadTransactions(name string) error {
	path := filepath.Join(this.dataPath, name+TransactionFileSuffix)
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, &this.transactions)
}

func (this *Project) SaveTransactions(name string) error {
	jsonBytes, err := json.Marshal(this.transactions)
	if err != nil {
		return err
	}
	path := filepath.Join(this.dataPath, name+TransactionFileSuffix)
	return os.WriteFile(path, jsonBytes, 0660)
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

func (this *Project) RemoveTransaction(operation string, parameters []string) error {
	return this.removeTransaction(operation, parameters)
}

///// private /////

func (this *Project) addTransaction(operation string, parameters ...interface{}) {
	this.transactions = append(this.transactions, Transaction{
		Operation:  operation,
		Parameters: parameters,
	})
}

func (this *Project) removeTransaction(operation string, parameters []string) error {
	for i, transaction := range this.transactions {
		if _, found := TransactionOperationMap[operation]; !found {
			continue
		}
		if transaction.Operation != TransactionOperationMap[operation] {
			continue
		}
		paramsMatched := 0
		for _, transactionParameter := range transaction.Parameters {
			for _, functionParameter := range parameters {
				testResult := false
				switch transactionParameter.(type) {
				case string:
					testResult = functionParameter == transactionParameter
				case bool:
					boolVal, _ := strconv.ParseBool(functionParameter)
					testResult = boolVal == transactionParameter
				case MemberName:
					testResult = MemberName(functionParameter) == transactionParameter
				case Organization:
					testResult = Organization(functionParameter) == transactionParameter
				}
				if testResult {
					paramsMatched++
				}
			}
		}
		if paramsMatched == len(transaction.Parameters) {
			this.transactions = append(this.transactions[:i], this.transactions[i+1:]...)
		}
	}

	this.playTransactions()
	return nil
}

func (this *Project) playTransactions() {
	freshCallings := this.originalCallings.copy()
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
