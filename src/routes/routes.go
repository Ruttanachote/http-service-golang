package routes

import (
	"github.com/Stream-I-T-Consulting/stream-http-service-go/app/http_server"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/src/handlers"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/src/repositories"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/src/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func HTTPRootMiddleware(s *http_server.HttpServer) {
	// Default middleware configs
	fiberRequestID := requestid.New()
	fiberETag := etag.New(s.Config.Fiber.Middleware.ETag)
	fiberCors := cors.New(s.Config.Fiber.Middleware.Cors)
	fiberLogger := logger.New(s.Config.Fiber.Middleware.Logger)
	fiberFavicon := favicon.New(s.Config.Fiber.Middleware.Favicon)
	fiberLimiter := limiter.New(s.Config.Fiber.Middleware.Limiter)
	fiberRecover := recover.New()

	s.Use(fiberRequestID)
	s.Use(fiberETag)
	s.Use(fiberCors)
	s.Use(fiberLogger)
	s.Use(fiberFavicon)
	s.Use(fiberLimiter)
	s.Use(fiberRecover)
}

func HTTPRootRoute(s *http_server.HttpServer) {
	HTTPRootMiddleware(s)

	s.GET("/", func(c *fiber.Ctx) error {
		return handlers.GetRootPath(c)
	})
	s.GET("/monitor", monitor.New(monitor.Config{Title: "Fiber Monitoring"}))
}

func HTTPRoutes(s *http_server.HttpServer) {
	HTTPRootRoute(s)

	// Initialize repositories, services, and handlers
	dbRepo := repositories.NewDbRepository(s.DbClient, s.Tracer)
	userRepo := repositories.NewUserRepository(s.DbClient, s.Tracer)

	// Initialize services
	userService := services.NewUserService(s.Tracer, userRepo)

	// Initialize handlers
	handler := handlers.NewHandler(
		s.Cacher,
		s.Tracer,
		dbRepo,
		userService,
	)

	// REST API endpoint ------------------------------------------------------------------
	s.GET("/health", func(c *fiber.Ctx) error { return handler.CheckDatabaseConnection(c) })

	// User service routes
	s.GET("/users", func(c *fiber.Ctx) error { return handler.GetUsers(c) })
	s.GET("/users/:id", func(c *fiber.Ctx) error { return handler.GetUser(c) })
	s.POST("/users", func(c *fiber.Ctx) error { return handler.CreateUser(c) })
	s.PUT("/users/:id", func(c *fiber.Ctx) error { return handler.UpdateUser(c) })
	s.DELETE("/users/:id", func(c *fiber.Ctx) error { return handler.DeleteUser(c) })
}
