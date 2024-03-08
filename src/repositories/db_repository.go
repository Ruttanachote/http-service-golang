package repositories

import (
	"context"
	"fmt"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/tracing"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type dbRepository struct {
	db     *gorm.DB
	tracer trace.Tracer
}

func NewDbRepository(db *gorm.DB, tracer trace.Tracer) DbRepository {
	return dbRepository{db: db, tracer: tracer}
}

func (r dbRepository) CheckDatabaseConnection(ctx context.Context) error {
	var (
		_, childSpan = tracing.TraceStart(ctx, r.tracer, "CheckDatabaseConnectionHandler", trace.WithAttributes(attribute.String("handler", "CheckDatabaseConnection")))
	)

	sqlDB, _ := r.db.DB()
	if err := sqlDB.Ping(); err != nil {
		utils.HandleErrors(ctx, fmt.Errorf("failed to connect to database server: connection refused"))
		return fiber.ErrServiceUnavailable
	}

	tracing.TraceEnd(childSpan)

	return nil
}
