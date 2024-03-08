package services

import (
	"context"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/database"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/tracing"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/src/models"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/src/repositories"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	userService struct {
		tracer         trace.Tracer
		userRepository repositories.UserRepository
	}
)

func NewUserService(
	tracer trace.Tracer,
	userRepo repositories.UserRepository,
) UserService {
	return &userService{
		tracer:         tracer,
		userRepository: userRepo,
	}
}

func (s userService) GetUsers(ctx context.Context, paginate database.Pagination, search string) (*database.Pagination, error) {
	ctx, childSpan := tracing.TraceStart(ctx, s.tracer, "GetUsersService", trace.WithAttributes(attribute.String("service", "GetUsers")))
	result, err := s.userRepository.GetUserPaginate(ctx, paginate, search)
	tracing.TraceEnd(childSpan)

	return result, err
}

func (s userService) GetUser(ctx context.Context, id int) (map[string]interface{}, error) {
	ctx, childSpan := tracing.TraceStart(ctx, s.tracer, "GetUserService", trace.WithAttributes(attribute.String("service", "GetUser")))
	user, err := s.userRepository.GetUserByID(ctx, id)
	tracing.TraceEnd(childSpan)

	return map[string]interface{}{"data": user}, err
}

func (s userService) CreateUser(ctx context.Context, userDto *UserDto) error {
	ctx, childSpan := tracing.TraceStart(ctx, s.tracer, "CreateUserService", trace.WithAttributes(attribute.String("service", "CreateUser")))
	user := new(models.User)

	user.FirstName = userDto.FirstName
	user.LastName = userDto.LastName
	user.Email = userDto.Email

	tracing.TraceEnd(childSpan)

	return s.userRepository.CreateUser(ctx, user)
}

func (s userService) UpdateUser(ctx context.Context, id int, userDto *UserDto) error {
	ctx, childSpan := tracing.TraceStart(ctx, s.tracer, "UpdateUserService", trace.WithAttributes(attribute.String("service", "UpdateUser")))
	user := new(models.User)

	user.FirstName = userDto.FirstName
	user.LastName = userDto.LastName
	user.Email = userDto.Email

	tracing.TraceEnd(childSpan)

	return s.userRepository.UpdateUser(ctx, id, user)
}

func (s userService) DeleteUser(ctx context.Context, id int) error {
	ctx, childSpan := tracing.TraceStart(ctx, s.tracer, "DeleteUserService", trace.WithAttributes(attribute.String("service", "DeleteUser")))
	err := s.userRepository.DeleteUser(ctx, id)
	tracing.TraceEnd(childSpan)

	return err
}
