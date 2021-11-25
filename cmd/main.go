package main

import (
	"fmt"
	"log"
	"net/http"

	"github.org/smartybryan/callingboard/config"
	"github.org/smartybryan/callingboard/engine"
	"github.org/smartybryan/callingboard/web"
)

func main() {
	appConfig := config.ParseConfig()

	members := engine.NewMembers(config.MaxMembers, appConfig.MembersDataPath)
	logOnError(members.Load())
	callings := engine.NewCallings(config.MaxCallings, appConfig.CallingDataPath)
	logOnError(callings.Load())
	project := engine.NewProject(&callings, &members, appConfig.DataPath)

	web.SetupRoutes(appConfig, web.NewController(project))

	fmt.Printf("Listening on port %s\n", appConfig.ListenPort)

	//secPath := path.Join(appConfig.DataPath, "sec")
	//certPath := path.Join(secPath, "server.crt")
	//keyPath := path.Join(secPath, "server.key")
	//secPath := "/var/lib/acme/live/callingboard.org"
	//certPath := path.Join(secPath, "cert")
	//keyPath := path.Join(secPath, "privkey")

	//panicOnError(http.ListenAndServeTLS(appConfig.ListenPort, certPath, keyPath, nil))
	panicOnError(http.ListenAndServe(appConfig.ListenPort,nil))
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func logOnError(err error) {
	if err != nil {
		log.Println(err)
	}
}
