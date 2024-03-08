package handlers

import (
	"context"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/cache"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/database"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/src/repositories"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/src/services"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

type (
	// Register handler services
	handler struct {
		cacher       *cache.Cache
		tracer       trace.Tracer
		userService  services.UserService
		dbRepository repositories.DbRepository
	}
	// Register handler interfaces
	Handler interface {
		DbHandler
		UserHandler
	}
)

func NewHandler(
	cacher *cache.Cache,
	tracer trace.Tracer,
	dbRepository repositories.DbRepository,
	userService services.UserService,
) handler {
	return handler{
		cacher:       cacher,
		tracer:       tracer,
		dbRepository: dbRepository,
		userService:  userService,
	}
}

type ServicePaginationFunc func(ctx context.Context, paginate database.Pagination, search string) (*database.Pagination, error)
type ServiceQueryFunc func(ctx context.Context, id int) (map[string]interface{}, error)

func (h handler) PaginationCache(ctx context.Context, key string, tags []string, paginate database.Pagination, search string, f ServicePaginationFunc) (*database.Pagination, error) {
	var (
		responseData *database.Pagination
		err          error
	)

	// Get the cached attributes object
	err = h.cacher.Get(ctx, key, &responseData)
	if err != nil {
		return nil, err
	}

	if responseData == nil {
		// Call service function
		responseData, err = f(ctx, paginate, search)
		if err != nil {
			return nil, err
		}

		// Set cache
		err = h.cacher.Tag(tags...).Set(ctx, key, &responseData)
		if err != nil {
			return nil, err
		}
	}

	return responseData, nil
}

func (h handler) QueryCache(ctx context.Context, key string, tags []string, id int, f ServiceQueryFunc) (map[string]interface{}, error) {
	var (
		responseData map[string]interface{}
		err          error
	)

	// Get the cached attributes object
	err = h.cacher.Get(ctx, key, &responseData)
	if err != nil {
		return nil, err
	}

	if responseData == nil {
		// Call service function
		responseData, err = f(ctx, id)
		if err != nil {
			return nil, err
		}

		// Set cache
		err = h.cacher.Tag(tags...).Set(ctx, key, &responseData)
		if err != nil {
			return nil, err
		}
	}

	return responseData, nil
}

// Root handlers  ------------------------------------------------------------------

func GetRootPath(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString(c.App().Config().AppName)
}
