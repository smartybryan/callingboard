package web

import (
	"net/http"
	"strconv"
	"strings"
)

type InputModel struct {
	MemberMinAge int
	MemberMaxAge int

}

func (this *InputModel) Bind(request *http.Request) error {
	this.MemberMinAge = atoi(strings.TrimSpace(request.Form.Get("min")))
	this.MemberMaxAge = atoi(strings.TrimSpace(request.Form.Get("max")))

	return nil
}

func (this *InputModel) Validate() error {
	if this.MemberMaxAge == 0 {
		this.MemberMaxAge = 120
	}
	if this.MemberMinAge > this.MemberMaxAge {
		this.MemberMinAge = 0
	}


	return nil
}

func atoi(value string) int {
	val, _ := strconv.Atoi(value)
	return val
}
