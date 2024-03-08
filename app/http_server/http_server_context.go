package http_server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/cache"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/tracing"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type (
	IHttpServer interface {
		StartServer() error
		StopServer()
		Cleanup() error
		Log(tag string, message string)

		startHTTP(exitChannel chan bool) error
		stopHTTP()

		// HTTP Services
		Use(args ...interface{})
		Group(prefix string, h ...func(ctx *fiber.Ctx) error) fiber.Router

		// HTTP Methods
		GET(path string, h ...func(ctx *fiber.Ctx) error) fiber.Router
		POST(path string, h ...func(ctx *fiber.Ctx) error) fiber.Router
		PUT(path string, h ...func(ctx *fiber.Ctx) error) fiber.Router
		PATCH(path string, h ...func(ctx *fiber.Ctx) error) fiber.Router
		DELETE(path string, h ...func(ctx *fiber.Ctx) error) fiber.Router
	}
	// HttpServer implement IHttpServer it is context for HTTP service
	HttpServer struct {
		Config        *config.Config
		fiber         *fiber.App
		DbClient      *gorm.DB
		Cacher        *cache.Cache
		traceProvider *sdktrace.TracerProvider
		Tracer        trace.Tracer
		exitChannel   chan bool
	}
)

// NewHttpServer is the constructor function for HttpServer
func NewHttpServer(
	config *config.Config,
	dbClient *gorm.DB,
	cacher *cache.Cache,
	traceProvider *sdktrace.TracerProvider,
	tracer trace.Tracer,
) *HttpServer {
	return &HttpServer{
		Config:        config,
		fiber:         fiber.New(config.Fiber.Config),
		DbClient:      dbClient,
		Cacher:        cacher,
		traceProvider: traceProvider,
		Tracer:        tracer,
	}
}

// Start start all registered services
func (s *HttpServer) StartServer() error {

	httpN := len(s.fiber.Stack())
	var exitHTTP chan bool
	if httpN > 0 {
		exitHTTP = make(chan bool, 1)
		go func() {
			s.startHTTP(exitHTTP)
		}()
	}

	// There are 2 ways to exit from Microservices
	// 1. The SigTerm can be send from outside program such as from k8s
	// 2. Send true to ms.exitChannel
	osQuit := make(chan os.Signal, 1)
	s.exitChannel = make(chan bool, 1)
	signal.Notify(osQuit, syscall.SIGTERM, syscall.SIGINT)
	exit := false
	for {
		if exit {
			break
		}
		select {
		case <-osQuit:
			// Exit from HTTP as well
			if exitHTTP != nil {
				exitHTTP <- true
			}
			exit = true
		case <-s.exitChannel:
			// Exit from HTTP as well
			if exitHTTP != nil {
				exitHTTP <- true
			}
			exit = true
		}
	}

	// Run cleanup
	s.Cleanup()
	return nil
}

// Stop server the services
func (s *HttpServer) StopServer() {
	if s.exitChannel == nil {
		return
	}
	s.exitChannel <- true
}

// Cleanup clean resources up from every registered services before exit
func (s *HttpServer) Cleanup() error {
	if !fiber.IsChild() {
		s.Log("HttpServer", "Start cleanup, close all client connections...")
	}

	// Close all database connection
	if s.DbClient != nil {
		sqlDB, _ := s.DbClient.DB()
		sqlDB.Close()
	}

	// Close all redis connection
	if s.Cacher != nil {
		s.Cacher.Close()
	}

	// Sentry: Flush buffered events before the program terminates.
	if s.Config.Sentry.SentryDSN != "" {
		sentry.Flush(2 * time.Second)
		sentry.Recover()
	}

	// Close all tracer connection
	if s.traceProvider != nil {
		tracing.Cleanup(s.traceProvider)
	}

	return nil
}

// Log message to console
func (s *HttpServer) Log(tag string, message string) {
	// _, fn, line, _ := runtime.Caller(1)
	// fns := strings.Split(fn, "/")
	log.Println(
		// fmt.Sprintf("[%s]", time.Now().Format("2006-01-02 15:04:05")),
		tag,
		// fns[len(fns)-1],
		// line,
		"|",
		message,
	)
}

// startHTTP will start HTTP service, this function will block thread
func (s *HttpServer) startHTTP(exitChannel chan bool) error {
	// Caller can exit by sending value to exitChannel
	go func() {
		<-exitChannel
		s.stopHTTP()
	}()

	return s.fiber.Listen(":" + s.Config.App.HttpPort)
}

// stopHTTP will stop HTTP service graceful shutdown
func (s *HttpServer) stopHTTP() {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if !fiber.IsChild() {
		s.Log("HttpServer", "Shutting down...")
	}
	s.fiber.Shutdown()
}
