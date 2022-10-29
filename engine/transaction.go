package engine

type Transaction struct {
	Operation  string
	Parameters []string
}

// map a web operation to a transaction operation
var TransactionOperationMap = map[string]string{
	"sustainings": OpAddMemberToACalling,
	"releases":    OpRemoveMemberFromACalling,
}
