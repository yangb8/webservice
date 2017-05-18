package statsd

import (
	"log"
	"time"

	"github.com/Unix4ever/statsd"
)

var StatsdClient *statsd.StatsdClient

// StartupStatsd initializes sentry client
func StartupStatsd(enabled bool, addr string) {
	if enabled {
		// TODO: make parameters configurable
		StatsdClient = statsd.NewStatsdClient(addr, "", 1400, time.Second*10, time.Minute*5)
		if err := StatsdClient.CreateSocket(); err != nil {
			log.Fatalf("Unable to initialize Statsd: %s", err)
		}
		log.Printf("Statsd is enabled")
	} else {
		log.Printf("Statsd is disabled")
	}
}

// ShutdownStatsd closes sentry client (if it was created)
func ShutdownStatsd() {
	if StatsdClient != nil {
		StatsdClient.Close()
	}
}
