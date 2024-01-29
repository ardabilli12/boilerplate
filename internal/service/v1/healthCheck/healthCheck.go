package healthCheck

import (
	"boilerplate-service/config"
	"boilerplate-service/internal/repository"
)

type healthCheck struct {
	config          *config.Config
	healthCheckRepo repository.IHealthCheckRepository
}

func New(config *config.Config, healthCheckRepo repository.IHealthCheckRepository) *healthCheck {
	return &healthCheck{
		config:          config,
		healthCheckRepo: healthCheckRepo,
	}
}
