package httputil_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	httputil "boilerplate-service/pkg/util/http"
)

func TestRequestHitAPI(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		uri                string
		data               interface{}
		header             map[string]string
		mockResponseStatus int
		mockResponseBody   string
		wantCode           int
		wantErr            bool
	}{
		{
			name:               "Successful GET Request",
			method:             "GET",
			uri:                "/success",
			data:               nil,
			header:             map[string]string{"Authorization": "Bearer token123", "Custom-Header": "Value123"},
			mockResponseStatus: http.StatusOK,
			mockResponseBody:   `{"status":"ok"}`,
			wantCode:           http.StatusOK,
			wantErr:            false,
		},
		{
			name:               "Successful POST Request",
			method:             "POST",
			uri:                "/create",
			data:               map[string]string{"Authorization": "Bearer token123", "Custom-Header": "Value123"},
			header:             nil,
			mockResponseStatus: http.StatusCreated,
			mockResponseBody:   `{"status":"created"}`,
			wantCode:           http.StatusCreated,
			wantErr:            false,
		},
		{
			name:               "Network Failure",
			method:             "GET",
			uri:                "/network-failure",
			data:               nil,
			header:             nil,
			mockResponseStatus: 0,
			mockResponseBody:   "",
			wantCode:           0,
			wantErr:            true,
		},
		{
			name:               "Error Unmarshalling JSON",
			method:             "GET",
			uri:                "/json-error",
			data:               nil,
			header:             nil,
			mockResponseStatus: http.StatusBadRequest,
			mockResponseBody:   "{invalid-json}",
			wantCode:           http.StatusBadRequest,
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockResponseStatus)
				w.Write([]byte(tt.mockResponseBody))
			}))
			defer mockServer.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			res, code, err := httputil.RequestHitAPI(ctx, tt.method, mockServer.URL+tt.uri, tt.data, tt.header)

			if (err != nil) != tt.wantErr {
				t.Errorf("%s: RequestHitAPI() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if code != tt.wantCode {
				t.Errorf("%s: RequestHitAPI() code = %v, want %v", tt.name, code, tt.wantCode)
			}
			if tt.mockResponseBody != "" && !strings.Contains(string(res), tt.mockResponseBody) {
				t.Errorf("%s: RequestHitAPI() response = %v, want contains %v", tt.name, string(res), tt.mockResponseBody)
			}
		})
	}
}
