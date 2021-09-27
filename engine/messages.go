package engine

import (
	"errors"
	"fmt"
)

var (
	ERROR_UNKNOWN_MEMBER = errors.New(fmt.Sprintf("Member name is unknown"))
	ERROR_UNKNOWN_ORGANIZATION =  errors.New(fmt.Sprintf("Organization does not exist"))
	ERROR_UNKNOWN_CALLING        = errors.New(fmt.Sprintf("Calling does not exist"))
	ERROR_MEMBER_INVALID_CALLING = errors.New(fmt.Sprintf("Member does not hold that calling"))
)
