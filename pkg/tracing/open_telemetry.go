package tracing

import (
	"context"
	"log"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/utils/color"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func InitTracer(config *config.Config) (*sdktrace.TracerProvider, trace.Tracer) {
	if config.OpenTelemetry.OtelExporterOTLPEndpoint == "" {
		return nil, nil
	}

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if config.OpenTelemetry.OtelInsecureMode {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(config.OpenTelemetry.OtelExporterOTLPEndpoint),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", config.App.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Println("Could not set resources:", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Check if OpenTelemetry Endpoint is already set
	if config.OpenTelemetry.OtelExporterOTLPEndpoint != "" {
		if !fiber.IsChild() {
			log.Println("OpenTelemetry: Tracing is", color.Format(color.GREEN, "on!"))
		}
	}

	// Set main tracer
	tracer := otel.Tracer(config.App.ServiceName)

	return tp, tracer
}

func TraceStart(ctx context.Context, tracer trace.Tracer, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if tracer == nil {
		return ctx, nil
	}

	return tracer.Start(ctx, spanName, opts...)
}

func TraceEnd(span trace.Span) {
	if span != nil {
		span.End()
	}
}

func Cleanup(traceProvider *sdktrace.TracerProvider) {
	if err := traceProvider.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}
}
