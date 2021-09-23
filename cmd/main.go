package main

import (
	"fmt"
	"net/http"

	"github.org/smartybryan/callorg/config"
	"github.org/smartybryan/callorg/engine"
	"github.org/smartybryan/callorg/web"
)

func main() {
	appConfig := config.ParseConfig()

	members := engine.NewMembers(config.MaxMembers)
	panicOnError(members.Load(appConfig.MembersDataPath))
	callings := engine.NewCallings(config.MaxCallings)
	panicOnError(callings.Load(appConfig.CallingDataPath))
	project := engine.NewProject(&callings, &members)

	web.SetupRoutes(appConfig, web.NewController(project))

	fmt.Printf("Listening on port %s\n", appConfig.ListenPort)
	_ = http.ListenAndServe(appConfig.ListenPort, nil)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
