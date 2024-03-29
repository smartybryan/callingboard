package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"sort"

	"github.com/nfnt/resize"
)

const (
	CallingNotEligible = iota
	CallingYouth
	CallingAdult
	AllEligible
	AllMembers
)

type Members struct {
	MemberMap    map[string]Member
	FocusMembers []string

	initialSize int
	filePath    string
}

type MemberWithFocus struct {
	Name  string
	Focus bool
}

func NewMembers(numMembers int, path string) Members {
	return Members{
		MemberMap:   make(map[string]Member, numMembers),
		initialSize: numMembers,
		filePath:    path,
	}
}

func (this *Members) AdultsWithoutACalling(callings Callings) (names []string) {
	return MemberSetDifference(this.GetMembers(CallingAdult), callings.MembersWithCallings())
}

func (this *Members) GetMemberRecord(name string) Member {
	if member, found := this.MemberMap[name]; found {
		return member
	}
	return Member{}
}

func (this *Members) GetMembers(eligibility uint8) (names []string) {
	for name, member := range this.MemberMap {
		if eligibility == AllMembers ||
			(eligibility == AllEligible && (member.Eligibility == CallingYouth || member.Eligibility == CallingAdult)) ||
			(eligibility == CallingYouth && member.Eligibility == CallingYouth) ||
			(eligibility == CallingAdult && member.Eligibility == CallingAdult) {

			names = append(names, name)
		}
	}
	sort.SliceStable(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	return names
}

func (this *Members) Load() error {
	jsonBytes, err := os.ReadFile(this.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, this)
}

func (this *Members) Save() (numObjects int, err error) {
	if len(this.MemberMap) == 0 {
		return 0, nil
	}
	jsonBytes, err := json.Marshal(this)
	if err != nil {
		return 0, err
	}
	err = os.WriteFile(this.filePath, jsonBytes, 0660)
	return len(this.MemberMap), err
}

//
//func (this *Members) GetMembersWithType(memberNames []string) (names []string) {
//	for _, name := range memberNames {
//		names = append(names, this.GetMemberRecord(name).BuildMemberName())
//	}
//	return names
//}

func (this *Members) GetMembersWithType(memberNames []string) (focusMembers []MemberWithFocus) {
	for _, name := range memberNames {
		focusMembers = append(focusMembers, MemberWithFocus{
			Name:  this.GetMemberRecord(name).BuildMemberName(),
			Focus: this.isMemberFocused(name),
		})
	}
	return focusMembers
}

func (this *Members) GetMembersWithFocus() (focusMembers []MemberWithFocus) {
	members := this.GetMembers(CallingAdult)

	for _, member := range members {
		focusMembers = append(focusMembers, MemberWithFocus{
			Name:  member,
			Focus: this.isMemberFocused(member),
		})
	}
	return focusMembers
}

func (this *Members) GetFocusMembers() []string {
	return this.FocusMembers
}

func (this *Members) SetMemberFocus(member string, focus bool) {
	for i := 0; i < len(this.FocusMembers); i++ {
		if this.FocusMembers[i] == member {
			this.FocusMembers = append(this.FocusMembers[:i], this.FocusMembers[i+1:]...)
			_, _ = this.Save()
			return
		}
	}
	if focus && !this.isMemberFocused(member) {
		this.FocusMembers = append(this.FocusMembers, member)
		_, _ = this.Save()
	}
}

func (this *Members) UploadMemberImage(imagePath, originalFile string, imageBytes []byte) error {
	imageType := path.Ext(originalFile)

	switch imageType {
	case ".jpg", ".jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return err
		}
		return this.resizeAndWriteImage(imagePath, img, err)
	case ".png":
		img, err := png.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return err
		}
		return this.resizeAndWriteImage(imagePath, img, err)
	default:
		return ERROR_UNSUPPORTED_IMAGE
	}
}

func (this *Members) resizeAndWriteImage(imagePath string, img image.Image, err error) error {
	newImg := resize.Resize(0, 300, img, resize.Lanczos3)
	var outBuffer bytes.Buffer
	err = jpeg.Encode(&outBuffer, newImg, nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(imagePath, outBuffer.Bytes(), os.FileMode(0666))
}

func (this *Members) DeleteMemberImage(path string) error {
	return os.Remove(path)
}

///// private /////

func (this *Members) copy() Members {
	newMembers := NewMembers(this.initialSize, this.filePath)
	newMembers.initialSize = this.initialSize
	newMembers.filePath = this.filePath

	for name, member := range this.MemberMap {
		newMembers.MemberMap[name] = Member{
			Name:        member.Name,
			Eligibility: member.Eligibility,
		}
	}

	return newMembers
}

func (this *Members) isMemberFocused(member string) bool {
	for _, focusedMember := range this.FocusMembers {
		if focusedMember == member {
			return true
		}
	}
	return false
}

//////////////////////////////////////////////////////

type Member struct {
	Name        string
	Eligibility uint8
	Type        uint8
}

func (this Member) BuildMemberName() string {
	return fmt.Sprintf("%s;%d", this.Name, this.Type)
}

func MemberSetDifference(mainSet, subtractSet []string) (names []string) {
	for _, name := range mainSet {
		if memberInSet(subtractSet, name) {
			continue
		}
		if !memberInSet(names, name) {
			names = append(names, name)
		}
	}
	return names
}

func memberInSet(set []string, name string) bool {
	for _, value := range set {
		if value == name {
			return true
		}
	}
	return false
}
