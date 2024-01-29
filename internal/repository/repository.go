package repository

import (
	"context"
)

type IHealthCheckRepository interface {
	CheckDB(ctx context.Context) error
	CheckRedis(ctx context.Context) string
}
