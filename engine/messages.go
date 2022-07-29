package engine

import (
	"errors"
	"fmt"
)

var (
	ERROR_UNKNOWN_MEMBER         = errors.New(fmt.Sprintf("Member name is unknown"))
	ERROR_UNKNOWN_ORGANIZATION   = errors.New(fmt.Sprintf("string does not exist"))
	ERROR_UNKNOWN_CALLING        = errors.New(fmt.Sprintf("Calling does not exist"))
	ERROR_MEMBER_INVALID_CALLING = errors.New(fmt.Sprintf("Member does not hold that calling"))
	ERROR_MEMBER_HAS_CALLING     = errors.New(fmt.Sprintf("Member already has that calling"))
	ERROR_INVALID_TRANSACTION    = errors.New(fmt.Sprintf("Transaction is invalid"))
	ERROR_UNSUPPORTED_IMAGE      = errors.New(fmt.Sprintf("Unsuppoted image type"))
)
