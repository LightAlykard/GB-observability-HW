package s

import (
	"log"

	"github.com/getsentry/sentry-go"
)

type SentryLogger struct {
	//l sentry
}

func NewSentryLogger() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://заменить.ingest.sentry.io/6640945",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	log.Println("sentry is initialized")
}
