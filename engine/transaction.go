package engine

type Transaction struct {
	Operation  string
	Parameters []interface{}
}

// map a web operation to a transaction operation
var TransactionOperationMap = map[string]string{
	"sustainings": "addMemberToACalling",
	"releases":  "removeMemberFromACalling",
}
