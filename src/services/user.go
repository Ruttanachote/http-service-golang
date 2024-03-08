package services

import (
	"context"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/database"
)

type (
	UserService interface {
		GetUsers(ctx context.Context, paginate database.Pagination, search string) (*database.Pagination, error)
		GetUser(ctx context.Context, id int) (map[string]interface{}, error)
		CreateUser(ctx context.Context, userDto *UserDto) error
		UpdateUser(ctx context.Context, id int, userDto *UserDto) error
		DeleteUser(ctx context.Context, id int) error
	}
	UserDto struct {
		FirstName string `json:"first_name" form:"first_name" query:"first_name" validate:"required,max:50"`
		LastName  string `json:"last_name" form:"last_name" query:"last_name" validate:"required,max:50"`
		Email     string `json:"email" form:"email" query:"email" validate:"required,email,max:100"`
	}
)
