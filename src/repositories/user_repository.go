package repositories

import (
	"context"
	"fmt"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/database"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/tracing"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/utils"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/src/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type userRepository struct {
	db     *gorm.DB
	tracer trace.Tracer
}

func NewUserRepository(db *gorm.DB, tracer trace.Tracer) UserRepository {
	return userRepository{db: db, tracer: tracer}
}

func (r userRepository) GetUserPaginate(ctx context.Context, pagination database.Pagination, search string) (*database.Pagination, error) {
	var (
		_, childSpan = tracing.TraceStart(ctx, r.tracer, "GetUserPaginateRepository", trace.WithAttributes(attribute.String("repository", "GetUserPaginate"), attribute.String("search", search)))
		users        []models.User
		err          error
	)

	// Pagination query
	if search != "" {
		if err = r.db.Scopes(database.Paginate(users, &pagination, r.db)).
			Where(`email LIKE ?`, fmt.Sprintf(`%%%s%%`, search)).
			Or(`first_name LIKE ?`, fmt.Sprintf(`%%%s%%`, search)).
			Or(`last_name LIKE ?`, fmt.Sprintf(`%%%s%%`, search)).
			Find(&users).Error; err != nil {
			utils.HandleErrors(ctx, err)
			return nil, utils.ErrQueryFailed
		}
	} else {
		if err = r.db.Scopes(database.Paginate(users, &pagination, r.db)).
			Find(&users).Error; err != nil {
			utils.HandleErrors(ctx, err)
			return nil, utils.ErrQueryFailed
		}
	}

	// Set data
	pagination.Data = users

	tracing.TraceEnd(childSpan)

	return &pagination, nil
}

func (r userRepository) GetUserByID(ctx context.Context, id int) (models.User, error) {
	var (
		_, childSpan = tracing.TraceStart(ctx, r.tracer, "GetUserByIDRepository", trace.WithAttributes(attribute.String("repository", "GetUserByID"), attribute.Int("id", id)))
		user         models.User
		err          error
	)

	// Query
	if err = r.db.First(&user, id).Error; err != nil {
		utils.HandleErrors(ctx, err)
		return user, err
	}

	tracing.TraceEnd(childSpan)

	return user, nil
}

func (r userRepository) CreateUser(ctx context.Context, user *models.User) error {
	var (
		_, childSpan = tracing.TraceStart(ctx, r.tracer, "CreateUserRepository", trace.WithAttributes(attribute.String("repository", "CreateUser")))
		err          error
	)

	// Execute
	if err = r.db.Create(&user).Error; err != nil {
		utils.HandleErrors(ctx, err)
		return utils.ErrQueryFailed
	}

	tracing.TraceEnd(childSpan)

	return nil
}

func (r userRepository) UpdateUser(ctx context.Context, id int, user *models.User) error {
	var (
		_, childSpan = tracing.TraceStart(ctx, r.tracer, "UpdateUserRepository", trace.WithAttributes(attribute.String("repository", "UpdateUser"), attribute.Int("id", id)))
		existUser    *models.User
		err          error
	)

	// Get model
	r.db.First(&existUser)

	// Set attributes
	existUser.FirstName = user.FirstName
	existUser.LastName = user.LastName
	existUser.Email = user.Email

	// Execute
	if err = r.db.Save(&existUser).Error; err != nil {
		utils.HandleErrors(ctx, err)
		return utils.ErrQueryFailed
	}

	tracing.TraceEnd(childSpan)

	return nil
}

func (r userRepository) DeleteUser(ctx context.Context, id int) error {
	var (
		_, childSpan = tracing.TraceStart(ctx, r.tracer, "DeleteUserRepository", trace.WithAttributes(attribute.String("repository", "DeleteUser"), attribute.Int("id", id)))
		err          error
	)

	// Execute
	if err = r.db.Delete(&models.User{}, id).Error; err != nil {
		utils.HandleErrors(ctx, err)
		return utils.ErrQueryFailed
	}

	tracing.TraceEnd(childSpan)

	return nil
}
