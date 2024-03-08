package http_server

import (
	"github.com/gofiber/fiber/v2"
)

// For Fiber middlewares
func (s *HttpServer) Use(args ...interface{}) {
	s.fiber.Use(args...)
}

// For Fiber route grouping
func (s *HttpServer) Group(prefix string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return s.fiber.Group(prefix, handlers...)
}

// GET register service endpoint for HTTP GET
func (s *HttpServer) GET(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return s.fiber.Get(path, handlers...)
}

// POST register service endpoint for HTTP POST
func (s *HttpServer) POST(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return s.fiber.Post(path, handlers...)
}

// PUT register service endpoint for HTTP PUT
func (s *HttpServer) PUT(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return s.fiber.Put(path, handlers...)
}

// PATCH register service endpoint for HTTP PATCH
func (s *HttpServer) PATCH(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return s.fiber.Patch(path, handlers...)
}

// DELETE register service endpoint for HTTP DELETE
func (s *HttpServer) DELETE(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return s.fiber.Delete(path, handlers...)
}
