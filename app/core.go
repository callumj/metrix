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
	http.Handle("/metric/increment", http.HandlerFunc(handlers.IncrementMetricHandler))
	http.ListenAndServe(listenOn, nil)
}
