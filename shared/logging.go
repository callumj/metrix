package shared

import (
	"github.com/getsentry/raven-go"
	"log"
)

var client *raven.Client

func HandleError(err error) {
	if client == nil {
		if len(Config.Sentry) != 0 {
			log.Println("Starting Sentry logger")
			recClient, innerErr := raven.NewClient(Config.Sentry, nil)
			if innerErr != nil {
				log.Fatalln("Unable to start Sentry logger")
			} else {
				client = recClient
			}
		}
	}

	if client == nil {
		return
	}

	innerError := err.Error()
	exp := raven.NewException(err, nil)

	packet := raven.NewPacket(innerError, exp)
	client.Capture(packet, nil)
}
