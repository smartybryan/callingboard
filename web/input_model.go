package web

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.org/smartybryan/callingboard/config"
)

type InputModel struct {
	ProjectHandle string
	Username      string
	Password      string
	WardId        string

	MemberMinAge int
	MemberMaxAge int
	MemberName   string

	Organization string
	FromOrg      string

	Calling       string
	FromCalling   string
	CustomCalling bool

	TransactionName   string
	TransactionParams string
	RawData           []byte
}

func (this *InputModel) Bind(request *http.Request) error {
	handle, err := request.Cookie(config.CookieName)
	if err == nil {
		this.ProjectHandle = handle.Value
	}

	this.Username = sanitize(request.Form.Get("username"))
	this.WardId = sanitize(request.Form.Get("wardid"))

	this.MemberMinAge = atoi(request.Form.Get("min"))
	this.MemberMaxAge = atoi(request.Form.Get("max"))
	this.MemberName = sanitize(request.Form.Get("member"))

	this.Organization = sanitize(request.Form.Get("org"))
	this.FromOrg = sanitize(request.Form.Get("from-org"))

	this.Calling = sanitize(request.Form.Get("calling"))
	this.FromCalling = sanitize(request.Form.Get("from-calling"))
	this.CustomCalling = atob(request.Form.Get("custom-calling"))

	this.TransactionName = sanitize(request.Form.Get("name"))
	this.TransactionParams = sanitize(request.Form.Get("params"))

	if request.Body != http.NoBody && request.Method == "POST" {
		fileUploadRequested := false
		err = request.ParseMultipartForm(5 * 1024 * 1024) // 5mb
		if err != nil {
			return err
		}
		file, _, err2 := request.FormFile("imageFile")
		if err2 == nil {
			fileUploadRequested = true
		}

		// upload image file
		if fileUploadRequested {
			defer func() { _ = file.Close() }()

			this.RawData, err = ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			return nil
		}

		// upload import data
		size := atoi(request.Header.Get("Content-Length"))
		if size == 0 {
			size = 128 * 1024
		}
		this.RawData = make([]byte, size)
		buf := make([]byte, 4096)

		var err error
		var numRead, currentLoc int
		for err != io.EOF {
			numRead, err = request.Body.Read(buf)
			copy(this.RawData[currentLoc:currentLoc+numRead], buf)
			currentLoc += numRead
		}
		if (numRead == 0 && currentLoc == 0) || err != nil && err != io.EOF {
			return errors.New("cannot read request body: " + err.Error())
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
	return strings.TrimSpace(value)
}

func atoi(value string) int {
	val, _ := strconv.Atoi(sanitize(value))
	return val
}

func atob(value string) bool {
	val, _ := strconv.ParseBool(sanitize(value))
	return val
}
