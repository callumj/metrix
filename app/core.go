package app

import (
	"github.com/callumj/metrix/handlers"
	"github.com/callumj/metrix/shared"
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

	http.Handle("/ping", http.HandlerFunc(handlers.PingHandler))
	http.Handle("/test", http.HandlerFunc(handlers.TestHandler))
	http.Handle("/metric/increment", http.HandlerFunc(handlers.IncrementMetricHandler))

	http.Handle("/api/sources", http.HandlerFunc(handlers.SourceListHandler))
	http.Handle("/api/keys", http.HandlerFunc(handlers.AvailableKeysHandler))
	http.Handle("/api/dates", http.HandlerFunc(handlers.DateKeysHandler))
	http.Handle("/api/counts", http.HandlerFunc(handlers.SubKeysHandler))

	http.ListenAndServe(listenOn, nil)
}
