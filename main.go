package main

import (
	"log"
	"time"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/app/http_server"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/cache"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/database"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/exceptions"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/tracing"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/src/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load environment variables
	globalConfig := config.NewConfig()

	// Initialize connection to database
	dbClient := database.Initialize(globalConfig)

	// Initialize connection to cache
	redisClient := cache.Initialize(globalConfig)
	cacher := cache.NewCacher(
		redisClient,
		cache.WithPrefix(globalConfig.Redis.RedisCachePrefix),
		cache.WithExpired(time.Minute*time.Duration(globalConfig.Redis.RedisCacheDuration)),
	)

	// Initialize OpenTelemetry tracing
	traceProvider, tracer := tracing.InitTracer(globalConfig)

	// Initialize Sentry client for error logging and tracing
	exceptions.SentryInitialize(globalConfig)

	// Create microservice instance
	httpServer := http_server.NewHttpServer(
		globalConfig,
		dbClient,
		cacher,
		traceProvider,
		tracer,
	)

	// Start http server
	httpRouting(httpServer)

	// Microservice start up
	httpServer.StartServer()
}

func httpRouting(s *http_server.HttpServer) {
	zone, _ := time.Now().Zone()
	if !fiber.IsChild() {
		log.Println("HTTP service is running")
		log.Println("[Timezone]:", zone)
	}
	routes.HTTPRoutes(s)
}
