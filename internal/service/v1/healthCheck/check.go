package healthCheck

import (
	healthCheckModel "boilerplate-service/internal/model/healthCheck"
	"boilerplate-service/pkg/newRelicExt"
	"context"
	"strings"
)

func (s *healthCheck) Check(ctx context.Context) healthCheckModel.HttpResponseHealthCheck {
	txn := newRelicExt.GetTxnFromCtx(ctx)
	if txn != nil {
		if segment := txn.StartSegment("internal/service/v1/healthCheck/check.go"); segment != nil {
			defer segment.End()
		}
	}

	var result healthCheckModel.HttpResponseHealthCheck

	redisCheck := s.healthCheckRepo.CheckRedis(ctx)
	if strings.ToLower(redisCheck) == "pong" {
		result.RedisAvailable = true
	}

	errMysql := s.healthCheckRepo.CheckDB(ctx)
	if errMysql == nil {
		result.MysqlAvailable = true
	}

	result.ServiceAvailable = true

	return result

}
