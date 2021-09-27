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


}
