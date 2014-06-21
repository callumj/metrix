package app

import (
	"github.com/callumj/metrix/handlers"
	"github.com/callumj/metrix/resource_bundle"
	"github.com/callumj/metrix/shared"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Run(args []string) {
	if len(args) > 1 {
		shared.LoadConfig(args[1])
		shared.ExplainConfig()
	}

	shared.InitializeRedisPool()

	listenOn := shared.Config.Listen
	if len(listenOn) == 0 {
		listenOn = ":8080"
	}

	log.Printf("Starting web server on %v", listenOn)

	resource_bundle.FetchFilesFromSelf()

	r := mux.NewRouter()

	r.HandleFunc("/ping", handlers.PingHandler)
	r.HandleFunc("/test", handlers.TestHandler)
	r.HandleFunc("/metric/increment", handlers.IncrementMetricHandler)

	r.HandleFunc("/api/sources", handlers.SourceListHandler)
	r.HandleFunc("/api/keys", handlers.AvailableKeysHandler)
	r.HandleFunc("/api/dates", handlers.DateKeysHandler)
	r.HandleFunc("/api/counts", handlers.SubKeysHandler)

	r.HandleFunc("/public/{path:.+}", handlers.PublicHandler)

	http.Handle("/", r)

	http.ListenAndServe(listenOn, nil)
}
