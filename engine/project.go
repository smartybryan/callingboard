package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	TransactionFileSuffix = ".txf"

	OpAddCalling                 = "addCalling"
	OpRemoveCalling              = "removeCalling"
	OpUpdateCalling              = "updateCalling"
	OpAddMemberToACalling        = "addMemberToACalling"
	OpRemoveMemberFromACalling   = "removeMemberFromACalling"
	OpMoveMemberToAnotherCalling = "moveMemberToAnotherCalling"
)

type Project struct {
	Callings         *Callings
	Members          *Members
	FocusMembers     *Members
	originalCallings Callings
	transactions     []Transaction
	undoHistory      []Transaction
	dataPath         string
	LastAccessed     time.Time

	diff DiffResult
}

type DiffResult struct {
	Sustainings  []Calling
	Releases     []Calling
	NewVacancies []Calling
	ModelName    string
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
	this.diff.NewVacancies = this.diff.NewVacancies[:0]

	modelVacancies := make([]Calling, 0, 20)
	originalVacancies := make([]Calling, 0, 20)

	for _, organization := range this.Callings.OrganizationOrder {
		// sustainings
		modelCallings := this.Callings.CallingList(organization)
		for _, modelCalling := range modelCallings {
			if modelCalling.Holder == VACANT_CALLING {
				modelVacancies = append(modelVacancies, modelCalling)
				continue
			}
			if this.originalCallings.doesMemberHoldCalling(modelCalling.Holder, organization, modelCalling.SubOrg, modelCalling.Name) {
				continue
			}
			this.diff.Sustainings = append(this.diff.Sustainings, modelCalling)
		}

		// releases
		originalCallings := this.originalCallings.CallingList(organization)
		for _, originalCalling := range originalCallings {
			if originalCalling.Holder == VACANT_CALLING {
				originalVacancies = append(originalVacancies, originalCalling)
				continue
			}
			if this.Callings.doesMemberHoldCalling(originalCalling.Holder, organization, originalCalling.SubOrg, originalCalling.Name) {
				continue
			}
			this.diff.Releases = append(this.diff.Releases, originalCalling)
		}
	}

	this.diff.NewVacancies = CallingSetDifference(modelVacancies, originalVacancies)

	return this.diff
}

func (this *Project) NewlyAvailableMembers() []string {
	return this.newlyAvailableMembers()
}

func (this *Project) RedoTransaction() bool {
	if len(this.undoHistory) == 0 {
		return false
	}
	this.transactions = append(this.transactions, this.undoHistory[len(this.undoHistory)-1])

	// calling focus undo
	transactionToRedo := this.transactions[len(this.transactions)-1]
	parameters := transactionToRedo.Parameters
	if transactionToRedo.Operation == OpRemoveMemberFromACalling && len(parameters) >= 4 {
		this.Callings.changeFocus(parameters[1], parameters[2], parameters[3], parameters[0], VACANT_CALLING)
		this.SaveCallingFocusList()
	}

	this.undoHistory = this.undoHistory[:len(this.undoHistory)-1]
	this.playTransactions()
	return true
}

func (this *Project) UndoTransaction() bool {
	if len(this.transactions) == 0 {
		return false
	}
	this.undoHistory = append(this.undoHistory, this.transactions[len(this.transactions)-1])

	// calling focus undo
	transactionToUndo := this.transactions[len(this.transactions)-1]
	parameters := transactionToUndo.Parameters
	if transactionToUndo.Operation == OpRemoveMemberFromACalling && len(parameters) >= 4 {
		this.Callings.changeFocus(parameters[1], parameters[2], parameters[3], VACANT_CALLING, parameters[0])
		this.SaveCallingFocusList()
	}
	if transactionToUndo.Operation == OpAddMemberToACalling && len(parameters) >= 4 {
		this.Callings.changeFocus(parameters[1], parameters[2], parameters[3], parameters[0], VACANT_CALLING)
		this.SaveCallingFocusList()
	}

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

func (this *Project) LoadTransactions(name string, overwrite bool) error {
	currentTransactions := len(this.transactions)
	if overwrite {
		currentTransactions = 0
	}
	_ = this.ResetModel()
	this.transactions = this.transactions[:currentTransactions]

	dataPath := filepath.Join(this.dataPath, name+TransactionFileSuffix)
	jsonBytes, err := os.ReadFile(dataPath)
	if err != nil {
		return err
	}
	var transactions []Transaction
	err = json.Unmarshal(jsonBytes, &transactions)
	if err != nil {
		return err
	}
	this.diff.ModelName = name

	if overwrite {
		this.transactions = this.transactions[:0]
	}
	for _, trans := range transactions {
		this.transactions = append(this.transactions, trans)
		// redundant transactions will be eliminated in playTransactions
	}
	this.playTransactions()
	return nil
}

func (this *Project) SaveTransactions(name string) error {
	jsonBytes, err := json.Marshal(this.transactions)
	if err != nil {
		return err
	}
	this.diff.ModelName = name
	dataPath := filepath.Join(this.dataPath, name+TransactionFileSuffix)
	return os.WriteFile(dataPath, jsonBytes, 0660)
}

func (this *Project) DeleteTransactions(name string) error {
	dataPath := filepath.Join(this.dataPath, name+TransactionFileSuffix)
	return os.Remove(dataPath)
}

func (this *Project) ResetModel() error {
	_ = this.Callings.Load()
	this.originalCallings = this.Callings.copy()
	this.transactions = this.transactions[:0]
	this.undoHistory = this.undoHistory[:0]
	this.diff.ModelName = ""
	return nil
}

func (this *Project) SetCallingFocus(org, suborg, callingName, member string, focus bool) {
	this.Callings.SetCallingFocus(org, suborg, callingName, member, focus)
	this.SaveCallingFocusList()
}

func (this *Project) SaveCallingFocusList() {
	saveCallingMap := this.Callings.CallingMap
	this.Callings.CallingMap = this.originalCallings.copy().CallingMap
	_, _ = this.Callings.Save()
	this.Callings.CallingMap = saveCallingMap
}

///// model modification stubs /////

func (this *Project) AddCalling(org string, calling string, custom bool) error {
	this.addTransaction(OpAddCalling, org, calling, boolToString(custom))
	return this.Callings.addCalling(org, calling, custom)
}

func (this *Project) RemoveCalling(org string, calling string) error {
	this.addTransaction(OpRemoveCalling, org, calling)
	return this.Callings.removeCalling(org, calling)
}

func (this *Project) UpdateCalling(org string, calling string, custom bool) error {
	this.addTransaction(OpUpdateCalling, org, calling, boolToString(custom))
	return this.Callings.updateCalling(org, calling, custom)
}

func (this *Project) AddMemberToACalling(member string, org string, suborg string, calling string) error {
	this.addTransaction(OpAddMemberToACalling, member, org, suborg, calling)
	err := this.Callings.addMemberToACalling(member, org, suborg, calling)
	this.SaveCallingFocusList()
	return err
}

func (this *Project) MoveMemberToAnotherCalling(
	member string, fromOrg string, fromSubOrg string, fromCalling string, toOrg string, toSubOrg string, toCalling string) error {
	this.addTransaction(OpMoveMemberToAnotherCalling, member, fromOrg, fromSubOrg, fromCalling, toOrg, toSubOrg, toCalling)
	return this.Callings.moveMemberToAnotherCalling(member, fromOrg, fromSubOrg, fromCalling, toOrg, toSubOrg, toCalling)
}

func (this *Project) RemoveMemberFromACalling(member string, org string, suborg string, calling string) error {
	this.addTransaction(OpRemoveMemberFromACalling, member, org, suborg, calling)
	err := this.Callings.removeMemberFromACalling(member, org, suborg, calling)
	this.SaveCallingFocusList()
	return err
}

func (this *Project) RemoveTransaction(operation string, parameters []string) error {
	return this.removeTransaction(operation, parameters)
}

func (this *Project) GetImagePath() string {
	imagePath := path.Join(path.Dir(this.dataPath), "html", path.Base(this.dataPath))
	_ = os.Mkdir(imagePath, os.FileMode(0777))
	return imagePath
}

///// private /////

func (this *Project) addTransaction(operation string, parameters ...string) {
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
				if functionParameter == transactionParameter {
					paramsMatched++
				}
			}
		}
		if paramsMatched == len(transaction.Parameters) {
			var transactionToRemove = this.transactions[i]
			if transactionToRemove.Operation == OpRemoveMemberFromACalling && len(parameters) >= 4 {
				this.Callings.changeFocus(parameters[1], parameters[2], parameters[3], VACANT_CALLING, parameters[0])
			}
			if transactionToRemove.Operation == OpAddMemberToACalling && len(parameters) >= 4 {
				this.Callings.changeFocus(parameters[1], parameters[2], parameters[3], parameters[0], VACANT_CALLING)
			}

			this.transactions = append(this.transactions[:i], this.transactions[i+1:]...)
		}
	}

	this.playTransactions()
	return nil
}

func (this *Project) deleteIrrelevantTransactions(invalidTransactions []int) {
	for i := len(invalidTransactions) - 1; i > -1; i-- {
		toDelete := invalidTransactions[i]
		this.transactions = append(this.transactions[:toDelete], this.transactions[toDelete+1:]...)
	}
}

func (this *Project) playTransactions() {
	var err error
	freshCallings := this.originalCallings.copy()
	this.Callings.CallingMap = freshCallings.CallingMap
	var irrelevantTransactions []int

	for idx, transaction := range this.transactions {
		switch transaction.Operation {
		case OpAddCalling:
			_ = this.Callings.addCalling(
				transaction.Parameters[0], transaction.Parameters[1], parseBool(transaction.Parameters[2]))
		case OpRemoveCalling:
			_ = this.Callings.removeCalling(
				transaction.Parameters[0], transaction.Parameters[1])
		case OpUpdateCalling:
			_ = this.Callings.updateCalling(
				transaction.Parameters[0], transaction.Parameters[1], parseBool(transaction.Parameters[2]))
		case OpAddMemberToACalling:
			err = this.Callings.addMemberToACalling(
				transaction.Parameters[0], transaction.Parameters[1], transaction.Parameters[2], transaction.Parameters[3])
		case OpMoveMemberToAnotherCalling:
			_ = this.Callings.moveMemberToAnotherCalling(
				transaction.Parameters[0],
				transaction.Parameters[1], transaction.Parameters[2], transaction.Parameters[3],
				transaction.Parameters[4], transaction.Parameters[5], transaction.Parameters[6])
		case OpRemoveMemberFromACalling:
			err = this.Callings.removeMemberFromACalling(
				transaction.Parameters[0], transaction.Parameters[1], transaction.Parameters[2], transaction.Parameters[3])
		}
		if err != nil {
			irrelevantTransactions = append(irrelevantTransactions, idx)
			err = nil
		}
	}

	// this can occur when the transaction is no longer relevant
	// such as when a dependent event was removed
	this.deleteIrrelevantTransactions(irrelevantTransactions)
}

func (this *Project) newlyAvailableMembers() []string {
	var releasedMembers, sustainedMembers []string
	for _, transaction := range this.transactions {
		if transaction.Operation == OpAddMemberToACalling {
			sustainedMembers = append(sustainedMembers, transaction.Parameters[0])
		}
		if transaction.Operation == OpRemoveMemberFromACalling {
			releasedMembers = append(releasedMembers, transaction.Parameters[0])
		}
	}
	availMembers := MemberSetDifference(releasedMembers, sustainedMembers)
	availMembers = MemberSetDifference(availMembers, this.Callings.MembersWithCallings())

	sort.Strings(availMembers)
	return availMembers
}

func boolToString(value bool) string {
	return fmt.Sprintf("%t", value)
}

func parseBool(value string) (val bool) {
	val, _ = strconv.ParseBool(value)
	return val
}
