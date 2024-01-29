package healthCheck

import (
	"boilerplate-service/internal/service"
	"boilerplate-service/pkg/newRelicExt"
	"boilerplate-service/pkg/util/response"
	"boilerplate-service/port/http/controller"
	"net/http"
)

type healthCheck struct {
	healthCheckSvc service.IHealthCheckService
}

func New(
	healthCheckSvc service.IHealthCheckService,
) controller.V1HealthCheckController {
	return &healthCheck{
		healthCheckSvc: healthCheckSvc,
	}

}

func (c *healthCheck) Check(w http.ResponseWriter, r *http.Request) {
	txn := newRelicExt.GetTxnFromCtx(r.Context())
	if txn != nil {
		if segment := txn.StartSegment("port/http/controller/v1/healthCheck/healthCheckController.go"); segment != nil {
			defer segment.End()
		}
	}

	check := c.healthCheckSvc.Check(r.Context())

	response.SendResponseOK(w, check)
}
