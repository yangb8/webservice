package sentry

import (
	"log"

	"github.com/getsentry/raven-go"
)

var SentryClient *raven.Client

// StartupSentry initializes sentry client
func StartupSentry(enabled bool, dsn string) {
	if enabled {
		var err error
		if SentryClient, err = raven.NewClient(dsn, map[string]string{}); err != nil {
			log.Fatalf("Unable to initialize Sentry: %s", err)
		}
		log.Printf("Sentry is enabled, project ID is: %s", SentryClient.ProjectID())
	} else {
		log.Printf("Sentry is disabled")
	}
}

// ShutdownSentry closes sentry client (if it was created)
func ShutdownSentry() {
	if SentryClient != nil {
		SentryClient.Close()
	}
}
