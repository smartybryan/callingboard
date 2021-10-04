package web

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type InputModel struct {
	MemberMinAge int
	MemberMaxAge int
	MemberName   string

	Organization string
	FromOrg      string

	Calling       string
	FromCalling   string
	CustomCalling bool

	TransactionName string
	RawData         []byte
}

func (this *InputModel) Bind(request *http.Request) error {
	this.MemberMinAge = atoi(request.Form.Get("min"))
	this.MemberMaxAge = atoi(request.Form.Get("max"))
	this.MemberName = sanitize(request.Form.Get("member"))

	this.Organization = sanitize(request.Form.Get("org"))
	this.FromOrg = sanitize(request.Form.Get("from-org"))

	this.Calling = sanitize(request.Form.Get("calling"))
	this.FromCalling = sanitize(request.Form.Get("from-calling"))
	this.CustomCalling = atob(request.Form.Get("custom-calling"))

	this.TransactionName = sanitize(request.Form.Get("name"))

	if request.Body != http.NoBody {
		size := atoi(request.Header.Get("Content-Length"))
		if size == 0 {
			size = 128*1024
		}
		this.RawData = make([]byte, size)
		numRead, err := request.Body.Read(this.RawData)
		if numRead == 0 || err != nil && err != io.EOF {
			return errors.New("cannot read body ")
		}
		_ = request.Body.Close()
	}

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

func sanitize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func atoi(value string) int {
	val, _ := strconv.Atoi(sanitize(value))
	return val
}

func atob(value string) bool {
	val, _ := strconv.ParseBool(sanitize(value))
	return val
}
