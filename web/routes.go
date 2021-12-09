package web

import (
	"net/http"

	"github.com/smartystreets/detour"
	"github.org/smartybryan/callingboard/config"
)

func SetupRoutes(appConfig config.Config, controller *Controller) {
	http.Handle("/", http.FileServer(http.Dir(appConfig.HtmlServerPath)))

	// authentication
	http.Handle("/v1/login", detour.New(controller.Login)) // username, wardid
	http.Handle("/v1/logout", detour.New(controller.Logout)) // username, wardid

	// members
	http.Handle("/v1/members", detour.New(controller.Members)) // min, max
	http.Handle("/v1/adults-without-calling", detour.New(controller.AdultsWithoutCalling))
	http.Handle("/v1/eligible-adults", detour.New(controller.AdultsEligibleForCalling))
	http.Handle("/v1/eligible-youth", detour.New(controller.YouthEligibleForCalling))
	http.Handle("/v1/newly-available", detour.New(controller.NewlyAvailableMembers))
	http.Handle("/v1/members-with-focus", detour.New(controller.GetMembersWithFocus))
	http.Handle("/v1/focus-members", detour.New(controller.GetFocusMembers))
	http.Handle("/v1/put-focus-members", detour.New(controller.PutFocusMembers)) // member (list sep by ;)
	http.Handle("/v1/get-member", detour.New(controller.GetMemberRecord)) // member
	http.Handle("/v1/load-members", detour.New(controller.LoadMembers))
	http.Handle("/v1/save-members", detour.New(controller.SaveMembers))
	http.Handle("/v1/parse-raw-members", detour.New(controller.ParseRawMembers))

	// callings
	http.Handle("/v1/organizations", detour.New(controller.OrganizationList))
	http.Handle("/v1/callings", detour.New(controller.CallingList)) // org
	http.Handle("/v1/callings-for-member", detour.New(controller.CallingListForMember)) // member
	http.Handle("/v1/vacant-calling-list", detour.New(controller.VacantCallingList)) // org
	http.Handle("/v1/members-with-callings", detour.New(controller.MembersWithCallings))
	http.Handle("/v1/add-calling", detour.New(controller.AddCalling)) // org, calling, custom-calling
	http.Handle("/v1/remove-calling", detour.New(controller.RemoveCalling)) // org, calling
	http.Handle("/v1/update-calling", detour.New(controller.UpdateCalling)) // org, calling, custom-calling
	http.Handle("/v1/add-member-calling", detour.New(controller.AddMemberToCalling)) // member, org, calling
	http.Handle("/v1/remove-member-calling", detour.New(controller.RemoveMemberFromCalling)) // member, org, calling
	http.Handle("/v1/move-member-calling", detour.New(controller.MoveMemberToAnotherCalling)) // member, from-org, from-calling, org, calling
	http.Handle("/v1/backout-transaction", detour.New(controller.RemoveTransaction)) // name, params (sep by semi-colon)
	http.Handle("/v1/load-callings", detour.New(controller.LoadCallings))
	http.Handle("/v1/save-callings", detour.New(controller.SaveCallings))
	http.Handle("/v1/parse-raw-callings", detour.New(controller.ParseRawCallings))

	// project
	http.Handle("/v1/diff", detour.New(controller.Diff))
	http.Handle("/v1/list-trans-files", detour.New(controller.ListTransactionFiles))
	http.Handle("/v1/load-trans", detour.New(controller.LoadTransactions)) // name
	http.Handle("/v1/save-trans", detour.New(controller.SaveTransactions)) // name
	http.Handle("/v1/delete-trans", detour.New(controller.DeleteTransactions)) // name
	http.Handle("/v1/reset-model", detour.New(controller.ResetModel))
	http.Handle("/v1/undo-trans", detour.New(controller.UndoTransaction))
	http.Handle("/v1/redo-trans", detour.New(controller.RedoTransaction))
}
