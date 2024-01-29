package healthCheck

import (
	"boilerplate-service/internal/repository"
	"boilerplate-service/pkg/logger"
	"boilerplate-service/pkg/mySqlExt"

	"boilerplate-service/pkg/redisExt"
	"context"
)

type healthCheck struct {
	logger logger.ILogger
	db     mySqlExt.IMySqlExt
	redis  redisExt.IRedisExt
}

func New(
	logger logger.ILogger,
	db mySqlExt.IMySqlExt,
	redis redisExt.IRedisExt,
) repository.IHealthCheckRepository {
	return &healthCheck{
		logger: logger,
		db:     db,
		redis:  redis,
	}
}

func (r *healthCheck) CheckRedis(ctx context.Context) string {
	return r.redis.Ping(ctx).String()
}

func (r *healthCheck) CheckDB(ctx context.Context) error {
	return r.db.Ping()
}
