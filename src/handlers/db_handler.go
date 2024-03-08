package handlers

import (
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/tracing"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	DbHandler interface {
		CheckDatabaseConnection(c *fiber.Ctx) error
	}
)

func (h handler) CheckDatabaseConnection(c *fiber.Ctx) error {
	var (
		ctx, span = tracing.TraceStart(c.Context(), h.tracer, "CheckDatabaseConnectionHandler", trace.WithAttributes(attribute.String("handler", "CheckDatabaseConnection")))
	)

	// Call service function
	err := h.dbRepository.CheckDatabaseConnection(ctx)
	if err != nil {
		return err
	}

	tracing.TraceEnd(span)
	return c.SendStatus(fiber.StatusOK)
}
