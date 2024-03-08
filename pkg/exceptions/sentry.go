package exceptions

import (
	"log"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/utils/color"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

func SentryInitialize(config *config.Config) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.Sentry.SentryDSN,
		EnableTracing:    config.Sentry.SentryEnableTracing,
		TracesSampleRate: config.Sentry.SentryTracesSampleRate,
	})

	// If there is an error, do not continue.
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// Check if Sentry DSN is already set
	if config.Sentry.SentryDSN != "" {
		if !fiber.IsChild() {
			log.Println("Sentry: Error handling is", color.Format(color.GREEN, "on!"))
			if config.Sentry.SentryEnableTracing {
				log.Println("Sentry: Tracing is", color.Format(color.GREEN, "on!"))
				log.Println("Sentry: Tracing sample rate is", config.Sentry.SentryTracesSampleRate)
			}
		}
	}
}
