package repositories

import (
	"context"
)

type (
	DbRepository interface {
		CheckDatabaseConnection(ctx context.Context) error
	}
)
