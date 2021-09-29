package web

import (
	"net/http"

	"github.com/smartystreets/detour"
	"github.org/smartybryan/callorg/config"
)

func SetupRoutes(appConfig config.Config, controller *Controller) {
	http.Handle("/", http.FileServer(http.Dir(appConfig.HtmlServerPath)))
	http.Handle("/v1/members", detour.New(controller.Members)) // min, max
	http.Handle("/v1/adults-without-calling", detour.New(controller.AdultsWithoutCalling))
	http.Handle("/v1/eligible-adults", detour.New(controller.AdultsEligibleForCalling))
	http.Handle("/v1/eligible-youth", detour.New(controller.YouthEligibleForCalling))
	http.Handle("/v1/get-member", detour.New(controller.GetMemberRecord)) // member
	http.Handle("/v1/organizations", detour.New(controller.OrganizationList))
	http.Handle("/v1/callings", detour.New(controller.CallingList)) // org


}

/*
func (this *Controller) MembersWithCallings() detour.Renderer {
func (this *Controller) VacantCallingList(input *InputModel) detour.Renderer {
func (this *Controller) AddCalling(input *InputModel) detour.Renderer {
func (this *Controller) RemoveCalling(input *InputModel) detour.Renderer {
func (this *Controller) UpdateCalling(input *InputModel) detour.Renderer {
func (this *Controller) AddMemberToCalling(input *InputModel) detour.Renderer {
func (this *Controller) MoveMemberToAnotherCalling(input *InputModel) detour.Renderer {
func (this *Controller) RemoveMemberFromCalling(input *InputModel) detour.Renderer {
func (this *Controller) Diff() detour.Renderer {
func (this *Controller) ListTransactionFiles() detour.Renderer {
func (this *Controller) LoadTransactions(input *InputModel) detour.Renderer {
func (this *Controller) SaveTransactions(input *InputModel) detour.Renderer {
func (this *Controller) UndoTransaction() detour.Renderer {
func (this *Controller) RedoTransaction() detour.Renderer {

func (this *Controller) LoadCallings() detour.Renderer {
func (this *Controller) SaveCallings() detour.Renderer {
func (this *Controller) LoadMembers() detour.Renderer {
func (this *Controller) SaveMembers() detour.Renderer {
parseMembers
parseCallings
 */