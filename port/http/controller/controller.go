package controller

import "net/http"

type V1HealthCheckController interface {
	Check(w http.ResponseWriter, r *http.Request)
}
