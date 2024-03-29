package web

import (
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.org/smartybryan/callingboard/config"
)

const maxImageFileSize = 1024 * 1024 * 10

type InputModel struct {
	ProjectHandle string
	Username      string
	Password      string
	WardId        string

	MemberMinAge  int
	MemberMaxAge  int
	MemberName    string
	ImageFileName string

	Organization    string
	SubOrganization string
	FromOrg         string
	FromSubOrg      string

	Calling     string
	FromCalling string
	Custom      bool

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
	this.MemberName = unescape(sanitize(request.Form.Get("member")))
	this.ImageFileName = unescape(sanitize(request.Form.Get("file")))

	this.Organization = sanitize(request.Form.Get("org"))
	this.SubOrganization = sanitize(request.Form.Get("suborg"))
	this.FromOrg = sanitize(request.Form.Get("from-org"))
	this.FromSubOrg = sanitize(request.Form.Get("from-suborg"))

	this.Calling = sanitize(request.Form.Get("calling"))
	this.FromCalling = sanitize(request.Form.Get("from-calling"))
	this.Custom = atob(request.Form.Get("custom"))

	this.TransactionName = sanitize(request.Form.Get("name"))
	this.TransactionParams = sanitize(request.Form.Get("params"))

	if request.Body != http.NoBody && request.Method == "POST" {
		fileUploadRequested := false
		err = request.ParseMultipartForm(maxImageFileSize)
		var file multipart.File
		if err == nil {
			file, _, err = request.FormFile("imageFile")
			if err == nil {
				fileUploadRequested = true
			}
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

func unescape(value string) string {
	val, _ := url.QueryUnescape(value)
	return val
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
