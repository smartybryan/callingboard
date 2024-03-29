package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/smartystreets/detour"
	"github.org/smartybryan/callingboard/config"
	"github.org/smartybryan/callingboard/engine"
)

const NotLoggedIn = "Not logged in"

type Controller struct {
	appConfig config.Config
	projects  map[string]*engine.Project
	mutex     sync.Mutex
}

func NewController(appConfig config.Config) *Controller {
	return &Controller{
		appConfig: appConfig,
		projects:  make(map[string]*engine.Project, 25),
	}
}

func (this *Controller) AddProject(handle string, project *engine.Project) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	project.LastAccessed = time.Now()
	this.projects[handle] = project
}

func (this *Controller) RemoveProject(handle string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	delete(this.projects, handle)
}

func (this *Controller) getProject(input *InputModel) *engine.Project {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.cleanProjects()

	project, found := this.projects[input.ProjectHandle]
	if !found {
		return nil
	} else {
		project.LastAccessed = time.Now()
		return project
	}
}

func (this *Controller) cleanProjects() {
	var handles []string
	for handle, _ := range this.projects {
		handles = append(handles, handle)
	}

	for _, handle := range handles {
		// clean projects from controller after 10 days of inactivity
		if time.Now().Sub(this.projects[handle].LastAccessed) > time.Hour*24*10 {
			delete(this.projects, handle)
		}
	}
}

/////////////// LOGIN

func (this *Controller) Login(input *InputModel) detour.Renderer {
	projectHandle := createCookieValue(input)
	input.ProjectHandle = projectHandle
	if project := this.getProject(input); project == nil {
		dataPath := path.Join(this.appConfig.DataPath, input.WardId)
		_ = os.Mkdir(dataPath, 0777)

		members := engine.NewMembers(config.MaxMembers, path.Join(dataPath, this.appConfig.MembersFile))
		logOnError(members.Load())
		callings := engine.NewCallings(config.MaxCallings, path.Join(dataPath, this.appConfig.CallingFile))
		logOnError(callings.Load())

		this.AddProject(projectHandle, engine.NewProject(&callings, &members, dataPath))
	}

	return detour.CookieResult{
		Cookie1: &http.Cookie{
			Name:   config.CookieName,
			Value:  projectHandle,
			Path:   "/",
			MaxAge: 86400 * 10, // 10 days
		},
	}
}

func (this *Controller) Logout(input *InputModel) detour.Renderer {
	this.RemoveProject(input.Username + ":" + input.WardId)
	return detour.JSONResult{
		StatusCode: 200,
		Content:    "Success",
	}
}

func (this *Controller) AuthenticationError() detour.Renderer {
	return detour.ContentResult{
		StatusCode: 401,
		Content:    NotLoggedIn,
	}
}

func createCookieValue(input *InputModel) string {
	return input.Username + ":" + input.WardId
}

func getWardFromCookieValue(value string) string {
	valueSplit := strings.Split(value, ":")
	if len(valueSplit) > 1 {
		return valueSplit[1]
	}
	return ""
}

/////////////// MEMBER

func (this *Controller) AdultsWithoutCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.GetMembersWithType(project.Members.AdultsWithoutACalling(*project.Callings)),
	}
}

func (this *Controller) Members(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.GetMembersWithType(project.Members.GetMembers(engine.AllEligible)),
	}
}

func (this *Controller) NewlyAvailableMembers(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.GetMembersWithType(project.NewlyAvailableMembers()),
	}
}

func (this *Controller) LoadMembers(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.Load(),
	}
}

func (this *Controller) SaveMembers(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	_, err := project.Members.Save()
	return detour.JSONResult{
		StatusCode: 200,
		Content:    err,
	}
}

func (this *Controller) GetMembersWithFocus(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.GetMembersWithFocus(),
	}
}

func (this *Controller) GetFocusMembers(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.GetMembersWithType(project.Members.GetFocusMembers()),
	}
}

func (this *Controller) SetMemberFocus(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	project.Members.SetMemberFocus(input.MemberName, input.Custom)
	return detour.JSONResult{
		StatusCode: 200,
		Content:    "",
	}
}

func (this *Controller) GetMembersWithType(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.GetMembersWithType(strings.Split(input.MemberName, "|")),
	}
}

func (this *Controller) ParseRawMembers(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	numMembers := project.Members.ParseMembersFromRawData(input.RawData)
	if numMembers < 10 {
		return detour.ContentResult{
			StatusCode: 422,
			Content:    "Unable to parse Member data",
		}
	}
	numObjects, err := project.Members.Save()
	msg := fmt.Sprintf("Imported %d members", numObjects)
	if err != nil {
		msg = err.Error()
	}
	return detour.ContentResult{
		StatusCode: 200,
		Content:    msg,
	}
}

func (this *Controller) ImageUpload(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}

	err := project.Members.UploadMemberImage(path.Join(project.GetImagePath(), input.MemberName+".jpg"), input.ImageFileName, input.RawData)
	if err != nil {
		return detour.ContentResult{
			StatusCode: 422,
			Content:    "Unable to upload Member Image: " + err.Error(),
		}
	}

	return detour.ContentResult{
		StatusCode: 200,
		Content:    "success",
	}
}

func (this *Controller) ImageDelete(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}

	err := project.Members.DeleteMemberImage(path.Join(project.GetImagePath(), input.MemberName+".jpg"))
	if err != nil {
		return detour.ContentResult{
			StatusCode: 422,
			Content:    "Unable to delete Member Image: " + err.Error(),
		}
	}

	return detour.ContentResult{
		StatusCode: 200,
		Content:    "success",
	}
}

///////////// CALLINGS

func (this *Controller) CallingList(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Callings.CallingList(input.Organization),
	}
}

func (this *Controller) CallingListForMember(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Callings.CallingListForMember(input.MemberName),
	}
}

func (this *Controller) MembersWithCallings(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.GetMembersWithType(project.Callings.MembersWithCallings()),
	}
}

func (this *Controller) OrganizationList(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Callings.OrganizationList(),
	}
}

func (this *Controller) VacantCallingList(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Callings.VacantCallingList(input.Organization),
	}
}

func (this *Controller) SetCallingFocus(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	project.SetCallingFocus(input.Calling, input.Custom)
	return detour.JSONResult{
		StatusCode: 200,
		Content:    "",
	}
}

func (this *Controller) LoadCallings(input *InputModel) detour.Renderer {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Callings.Load(),
	}
}

func (this *Controller) SaveCallings(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	_, err := project.Callings.Save()
	return detour.JSONResult{
		StatusCode: 200,
		Content:    err,
	}
}

func (this *Controller) ParseRawCallings(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	numCallings := project.Callings.ParseCallingsFromRawData(input.RawData)
	project.FinishCallingImport()
	if numCallings < 10 {
		return detour.ContentResult{
			StatusCode: 422,
			Content:    "Unable to parse Calling data",
		}
	}
	numObjects, err := project.Callings.Save()
	msg := fmt.Sprintf("Imported %d callings", numObjects)
	if err != nil {
		msg = err.Error()
	}
	return detour.ContentResult{
		StatusCode: 200,
		Content:    msg,
	}
}

///////////// PROJECT

func (this *Controller) AddCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.AddCalling(input.Organization, input.Calling, input.Custom),
	}
}

func (this *Controller) RemoveCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.RemoveCalling(input.Organization, input.Calling),
	}
}

func (this *Controller) UpdateCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.UpdateCalling(input.Organization, input.Calling, input.Custom),
	}
}

func (this *Controller) AddMemberToCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.AddMemberToACalling(input.MemberName, input.Organization, input.SubOrganization, input.Calling),
	}
}

func (this *Controller) MoveMemberToAnotherCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content: project.MoveMemberToAnotherCalling(input.MemberName, input.FromOrg, input.FromSubOrg, input.FromCalling,
			input.Organization, input.SubOrganization, input.Calling),
	}
}

func (this *Controller) RemoveMemberFromCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.RemoveMemberFromACalling(input.MemberName, input.Organization, input.SubOrganization, input.Calling),
	}
}

func (this *Controller) RemoveTransaction(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content: project.RemoveTransaction(
			input.TransactionName, strings.Split(input.TransactionParams, ":")),
	}
}

func (this *Controller) Diff(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Diff(),
	}
}

func (this *Controller) ListTransactionFiles(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.ListTransactionFiles(),
	}
}

func (this *Controller) LoadTransactions(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.LoadTransactions(input.TransactionName, true),
	}
}

func (this *Controller) MergeTransactions(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.LoadTransactions(input.TransactionName, false),
	}
}

func (this *Controller) SaveTransactions(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.SaveTransactions(input.TransactionName),
	}
}

func (this *Controller) DeleteTransactions(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.DeleteTransactions(input.TransactionName),
	}
}

func (this *Controller) ResetModel(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.ResetModel(),
	}
}

func (this *Controller) UndoTransaction(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.UndoTransaction(),
	}
}

func (this *Controller) RedoTransaction(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.RedoTransaction(),
	}
}

func logOnError(err error) {
	if err != nil {
		log.Println(err)
	}
}
