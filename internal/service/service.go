package service

import (
	healthCheckModel "boilerplate-service/internal/model/healthCheck"
	"context"
)

type IHealthCheckService interface {
	Check(ctx context.Context) healthCheckModel.HttpResponseHealthCheck
}
