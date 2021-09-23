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

func (this *Controller) AdultsEligibleForCalling(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.AdultsEligibleForACalling(),
	}
}

func (this *Controller) AdultsWithoutCalling(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.AdultsWithoutACalling(*this.project.Callings),
	}
}

func (this *Controller) Members(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.GetMembers(input.MemberMinAge, input.MemberMaxAge),
	}
}

func (this *Controller) YouthEligibleForCalling(input *InputModel) detour.Renderer {
	return detour.JSONResult{
		StatusCode:  200,
		Content:     this.project.Members.YouthEligibleForACalling(),
	}
}
