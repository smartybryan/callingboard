package web

import (
	"github.com/smartystreets/detour"
	"github.org/smartybryan/callorg/engine"
)

type Controller struct {
	project *engine.Project
}

func NewController(project *engine.Project) *Controller {
	return &Controller{
		project: project,
	}
}

/////////////// MEMBER

func (this *Controller) AdultsEligibleForCalling() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.AdultsEligibleForACalling(),
	}
}

func (this *Controller) AdultsWithoutCalling() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.AdultsWithoutACalling(*this.project.Callings),
	}
}

func (this *Controller) GetMemberRecord(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.GetMemberRecord(engine.MemberName(input.MemberName)),
	}
}

func (this *Controller) Members(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.GetMembers(input.MemberMinAge, input.MemberMaxAge),
	}
}

func (this *Controller) YouthEligibleForCalling() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.YouthEligibleForACalling(),
	}
}

func (this *Controller) LoadMembers() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.Load(),
	}
}

func (this *Controller) SaveMembers() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.Save(),
	}
}

func (this *Controller) ParseRawMembers(input *InputModel) detour.Renderer {
	this.project.Members.ParseMembersFromRawData(input.RawData)
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.Save(),
	}
}

///////////// CALLINGS

func (this *Controller) CallingList(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Callings.CallingList(engine.Organization(input.Organization)),
	}
}

func (this *Controller) MembersWithCallings() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Callings.MembersWithCallings(),
	}
}

func (this *Controller) OrganizationList() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Callings.OrganizationList(),
	}
}

func (this *Controller) VacantCallingList(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Callings.VacantCallingList(engine.Organization(input.Organization)),
	}
}

func (this *Controller) LoadCallings() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Callings.Load(),
	}
}

func (this *Controller) SaveCallings() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Callings.Save(),
	}
}

func (this *Controller) ParseRawCallings(input *InputModel) detour.Renderer {
	this.project.Callings.ParseCallingsFromRawData(input.RawData)
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Callings.Save(),
	}
}

///////////// PROJECT

func (this *Controller) AddCalling(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.AddCalling(engine.Organization(input.Organization), input.Calling, input.CustomCalling),
	}
}

func (this *Controller) RemoveCalling(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.RemoveCalling(engine.Organization(input.Organization), input.Calling),
	}
}

func (this *Controller) UpdateCalling(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.UpdateCalling(
			engine.Organization(input.Organization), input.Calling, input.CustomCalling),
	}
}

func (this *Controller) AddMemberToCalling(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.AddMemberToACalling(
			engine.MemberName(input.MemberName), engine.Organization(input.Organization), input.Calling),
	}
}

func (this *Controller) MoveMemberToAnotherCalling(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.MoveMemberToAnotherCalling(
			engine.MemberName(input.MemberName),
			engine.Organization(input.FromOrg), input.FromCalling,
			engine.Organization(input.Organization), input.Calling),
	}
}

func (this *Controller) RemoveMemberFromCalling(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.RemoveMemberFromACalling(
			engine.MemberName(input.MemberName), engine.Organization(input.Organization), input.Calling),
	}
}

func (this *Controller) Diff() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Diff(),
	}
}

func (this *Controller) ListTransactionFiles() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.ListTransactionFiles(),
	}
}

func (this *Controller) LoadTransactions(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.LoadTransactions(input.TransactionName),
	}
}

func (this *Controller) SaveTransactions(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.SaveTransactions(input.TransactionName),
	}
}

func (this *Controller) UndoTransaction() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.UndoTransaction(),
	}
}

func (this *Controller) RedoTransaction() detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.RedoTransaction(),
	}
}
