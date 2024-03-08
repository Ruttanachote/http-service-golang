package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// OAuth config
var OAuthConfig *certificateConfig

type Config struct {
	App           *appConfig
	Fiber         *fiberConfig
	Database      *databaseConfig
	Redis         *redisConfig
	Sentry        *sentryConfig
	OpenTelemetry *openTelemetryConfig
}

type appConfig struct {
	Env         string
	HttpPort    string
	GrpcPort    string
	AppName     string
	ServiceName string
}

type fiberConfig struct {
	Config     fiber.Config
	Middleware *fiberMiddlewareConfig
}

type fiberMiddlewareConfig struct {
	ETag    etag.Config
	Cors    cors.Config
	Logger  logger.Config
	Favicon favicon.Config
	Limiter limiter.Config
}

type databaseConfig struct {
	DatabaseDSN string
	// Additional
	DatabaseMaxIdleConns int
	DatabaseMaxOpenConns int
}

type redisConfig struct {
	RedisHost          string
	RedisPort          string
	RedisUsername      string
	RedisPassword      string
	RedisCachePrefix   string
	RedisCacheDuration int
}

type sentryConfig struct {
	SentryDSN              string
	SentryEnableTracing    bool
	SentryTracesSampleRate float64
}

type openTelemetryConfig struct {
	OtelExporterOTLPEndpoint string
	OtelInsecureMode         bool
}

type certificateConfig struct {
	PublicKey  string
	PrivateKey string
}

func NewConfig() *Config {
	// Load .env file
	godotenv.Load()

	// Set secrets
	OAuthConfig = &certificateConfig{
		PublicKey:  os.Getenv("OAUTH_PUBLIC_KEY"),
		PrivateKey: os.Getenv("OAUTH_PRIVATE_KEY"),
	}

	return &Config{
		App: &appConfig{
			Env:         os.Getenv("ENV"),
			HttpPort:    os.Getenv("HTTP_PORT"),
			GrpcPort:    os.Getenv("GRPC_PORT"),
			AppName:     os.Getenv("APP_NAME"),
			ServiceName: os.Getenv("SERVICE_NAME"),
		},
		Fiber: NewFiberConfig(),
		Database: &databaseConfig{
			DatabaseDSN: func() string {
				return fmt.Sprintf(
					`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s`,
					os.Getenv("DATABASE_HOST"),
					os.Getenv("DATABASE_USER"),
					os.Getenv("DATABASE_PASSWORD"),
					os.Getenv("DATABASE_NAME"),
					os.Getenv("DATABASE_PORT"),
					os.Getenv("DATABASE_TIMEZONE"),
				)
			}(),
			DatabaseMaxIdleConns: func() int {
				// Default max idle conns is 2
				databaseMaxIdleConns := 2
				envDatabaseMaxIdleConns, err := strconv.Atoi(os.Getenv("DATABASE_MAX_IDLE_CONNS"))
				if err == nil {
					databaseMaxIdleConns = envDatabaseMaxIdleConns
				}
				return databaseMaxIdleConns
			}(),
			DatabaseMaxOpenConns: func() int {
				// Default max open conns is 3
				databaseMaxOpenConns := 3
				envDatabaseMaxOpenConns, err := strconv.Atoi(os.Getenv("DATABASE_MAX_OPEN_CONNS"))
				if err == nil {
					databaseMaxOpenConns = envDatabaseMaxOpenConns
				}
				return databaseMaxOpenConns
			}(),
		},
		Redis: &redisConfig{
			RedisHost:        os.Getenv("REDIS_HOST"),
			RedisPort:        os.Getenv("REDIS_PORT"),
			RedisUsername:    os.Getenv("REDIS_USERNAME"),
			RedisPassword:    os.Getenv("REDIS_PASSWORD"),
			RedisCachePrefix: os.Getenv("REDIS_CACHE_PREFIX"),
			RedisCacheDuration: func() int {
				// Default max cache duration is 5
				redisCacheDuration := 5
				envRedisCacheDuration, err := strconv.Atoi(os.Getenv("REDIS_CACHE_DURATION"))
				if err == nil {
					redisCacheDuration = envRedisCacheDuration
				}
				return redisCacheDuration
			}(),
		},
		Sentry: &sentryConfig{
			SentryDSN: os.Getenv("SENTRY_DSN"),
			SentryEnableTracing: func() bool {
				// Default is false
				sentryEnableTracing := false
				envSentryEnableTracing, err := strconv.ParseBool(os.Getenv("SENTRY_ENABLE_TRACING"))
				if err == nil {
					sentryEnableTracing = envSentryEnableTracing
				}
				return sentryEnableTracing
			}(),
			SentryTracesSampleRate: func() float64 {
				// Default traces sample rate is 0.2
				sentryTracesSampleRate := 0.2
				envSentryTracesSampleRate, err := strconv.ParseFloat(os.Getenv("SENTRY_TRACES_SAMPLE_RATE"), 64)
				if err == nil {
					sentryTracesSampleRate = envSentryTracesSampleRate
				}
				return sentryTracesSampleRate
			}(),
		},
		OpenTelemetry: &openTelemetryConfig{
			OtelExporterOTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
			OtelInsecureMode: func() bool {
				// Default insecure mode is true
				otelInsecureMode := true
				envOtelInsecureMode, err := strconv.ParseBool(os.Getenv("OTEL_INSECURE_MODE"))
				if err == nil {
					otelInsecureMode = envOtelInsecureMode
				}
				return otelInsecureMode
			}(),
		},
	}
}
