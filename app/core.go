package app

import (
	"fmt"
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

	r.HandleFunc("/version", handlers.VersionSetHandler)

	r.HandleFunc("/api/sources", handlers.SourceListHandler)
	r.HandleFunc("/api/keys", handlers.AvailableKeysHandler)
	r.HandleFunc("/api/dates", handlers.DateKeysHandler)
	r.HandleFunc("/api/counts", handlers.SubKeysHandler)

	r.HandleFunc("/api/versions", handlers.VersionGetHandler)

	allocatePublicResources(r)

	http.Handle("/", r)

	http.ListenAndServe(listenOn, nil)
}

func allocatePublicResources(router *mux.Router) {
	if len(resource_bundle.AssetKeys) != 0 {
		for _, key := range resource_bundle.AssetKeys {
			route := fmt.Sprintf("/%v", key)
			log.Printf("Allocating route for %v\r\n", route)
			router.HandleFunc(route, handlers.PublicProdHandler)
			if route == "/index.html" {
				router.HandleFunc("/", handlers.PublicProdHandler)
			}
		}
	} else {
		router.HandleFunc("/public/{path:.+}", handlers.PublicDevHandler)
	}
}
