package main

import (
	"log"
	"net/http"
	"path"

	"github.org/smartybryan/callingboard/config"
	"github.org/smartybryan/callingboard/web"
)

func main() {
	appConfig := config.ParseConfig()
	web.SetupRoutes(appConfig, web.NewController(appConfig))

	switch appConfig.ListenPort {
	case config.ListenPortClearDefault: // for local debugging
		log.Printf("Listening on port %s\n", appConfig.ListenPort)
		panicOnError(http.ListenAndServe(appConfig.ListenPort, nil))
	case config.ListenPortDefault:
		log.Printf("Listening on port %s\n", appConfig.ListenPort)
		secPath := "/var/lib/acme/live/callingboard.org"
		certPath := path.Join(secPath, "cert")
		keyPath := path.Join(secPath, "privkey")
		go listenAndRedirectToTLS()
		panicOnError(http.ListenAndServeTLS(appConfig.ListenPort, certPath, keyPath, nil))
	}
}

func listenAndRedirectToTLS() {
	panicOnError(http.ListenAndServe(config.ListenPortClearDefault, http.HandlerFunc(redirectTLS)))
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	log.Printf("Redirect to: %s\n", "https://"+r.Host+r.RequestURI)
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
