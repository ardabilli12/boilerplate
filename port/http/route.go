package http

import (
	"boilerplate-service/pkg/logger"
	"boilerplate-service/pkg/newRelicExt"
	"boilerplate-service/port/http/controller"
	customMiddleware "boilerplate-service/port/http/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func HttpRoute(
	app newRelicExt.INewRelicExt,
	logger logger.ILogger,
	v1HealthCheckController controller.V1HealthCheckController,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Custom Middleware e.g idempotency etc
	r.Use(customMiddleware.LoggerMiddleware(app, logger))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health-check", v1HealthCheckController.Check)
	})

	return r
}
