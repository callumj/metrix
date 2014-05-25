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
			client, err = raven.NewClient(Config.Sentry, nil)
			if err != nil {
				log.Fatalln("Unable to start Sentry logger")
			}
		}
	}

	if client == nil {
		return
	}

	packet := raven.NewPacket(err.Error(), raven.NewException(err, nil))
	client.Capture(packet, nil)
}
