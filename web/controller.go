package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/smartystreets/detour"
	"github.org/smartybryan/callingboard/config"
	"github.org/smartybryan/callingboard/engine"
)

type Controller struct {
	appConfig config.Config
	projects  map[string]*engine.Project
}

func NewController(appConfig config.Config) *Controller {
	return &Controller{
		appConfig: appConfig,
		projects:  make(map[string]*engine.Project, 25),
	}
}

func (this *Controller) AddProject(handle string, project *engine.Project) {
	this.projects[handle] = project
}

func (this *Controller) RemoveProject(handle string) {
	if _, found := this.projects[handle]; found {
		delete(this.projects, handle)
	}
}

func (this *Controller) getProject(input *InputModel) *engine.Project {
	if project, found := this.projects[input.ProjectHandle]; !found {
		return nil
	} else {
		return project
	}
}

/////////////// LOGIN

// TODO: need a process to clean out old projects.
//       update timestamp on project for each getProject() call
//       then in that same function, delete data not used in two weeks

func (this *Controller) Login(input *InputModel) detour.Renderer {
	projectHandle := createCookieValue(input)
	if project := this.getProject(input); project == nil {
		this.AddProject(projectHandle, engine.NewProject(this.appConfig))
	}

	return detour.CookieResult{
		Cookie1: &http.Cookie{
			Name:       config.CookieName,
			Value:      projectHandle,
			Path:       "/",
			MaxAge:     86400 * 30,
		},
	}
}

func (this *Controller) Logout(input *InputModel) detour.Renderer {
	this.RemoveProject(input.ProjectHandle)
	return detour.JSONResult{
		StatusCode: 200,
		Content:    "Success",
	}
}

func (this *Controller) AuthenticationError() detour.Renderer {
	return detour.JSONResult{
		StatusCode: 401,
		Content:    "Not logged in",
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

func (this *Controller) AdultsEligibleForCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.AdultsEligibleForACalling(),
	}
}

func (this *Controller) AdultsWithoutCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.AdultsWithoutACalling(*project.Callings),
	}
}

func (this *Controller) GetMemberRecord(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.GetMemberRecord(input.MemberName),
	}
}

func (this *Controller) Members(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.GetMembers(input.MemberMinAge, input.MemberMaxAge),
	}
}

func (this *Controller) YouthEligibleForCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.YouthEligibleForACalling(),
	}
}

func (this *Controller) NewlyAvailableMembers(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.NewlyAvailableMembers(),
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
		Content:    project.Members.GetFocusMembers(),
	}
}

func (this *Controller) PutFocusMembers(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.Members.PutFocusMembers(strings.Split(input.MemberName, "|")),
	}
}

func (this *Controller) ParseRawMembers(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	numMembers := project.Members.ParseMembersFromRawData(input.RawData)
	if numMembers < 10 {
		return detour.JSONResult{
			StatusCode: 422,
			Content:    "Unable to parse Member data",
		}
	}
	numObjects, err := project.Members.Save()
	msg := fmt.Sprintf("Imported %d members", numObjects)
	if err != nil {
		msg = err.Error()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    msg,
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
		Content:    project.Callings.MembersWithCallings(),
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

func (this *Controller) LoadCallings(input *InputModel) detour.Renderer {
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
	if numCallings < 10 {
		return detour.JSONResult{
			StatusCode: 422,
			Content:    "Unable to parse Calling data",
		}
	}
	numObjects, err := project.Callings.Save()
	msg := fmt.Sprintf("Imported %d callings", numObjects)
	if err != nil {
		msg = err.Error()
	}
	return detour.JSONResult{
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
		Content:    project.AddCalling(input.Organization, input.Calling, input.CustomCalling),
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
		Content:    project.UpdateCalling(input.Organization, input.Calling, input.CustomCalling),
	}
}

func (this *Controller) AddMemberToCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.AddMemberToACalling(input.MemberName, input.Organization, input.Calling),
	}
}

func (this *Controller) MoveMemberToAnotherCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content: project.MoveMemberToAnotherCalling(input.MemberName, input.FromOrg, input.FromCalling,
			input.Organization, input.Calling),
	}
}

func (this *Controller) RemoveMemberFromCalling(input *InputModel) detour.Renderer {
	project := this.getProject(input)
	if project == nil {
		return this.AuthenticationError()
	}
	return detour.JSONResult{
		StatusCode: 200,
		Content:    project.RemoveMemberFromACalling(input.MemberName, input.Organization, input.Calling),
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
		Content:    project.LoadTransactions(input.TransactionName),
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
