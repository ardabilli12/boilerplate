package middleware

import (
	"boilerplate-service/constant"
	"boilerplate-service/pkg/logger"
	"boilerplate-service/pkg/newRelicExt"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

// ResponseWriter is a wrapper around http.ResponseWriter that captures the response body
type ResponseWriter struct {
	http.ResponseWriter
	body   *bytes.Buffer
	Status int
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteHeader overrides the default WriteHeader method to capture the response status
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.Status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func getDurationInMilliseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}

func LoggerMiddleware(app newRelicExt.INewRelicExt, logger logger.ILogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			txn := app.App().StartTransaction("HTTP Request")
			defer txn.End()

			// Adding the Method and URL Path as the segment name
			segmentName := r.Method + " " + r.URL.Path
			segment := txn.StartSegment(segmentName)
			defer segment.End()

			// Add the transaction to the context
			r = newrelic.RequestWithTransactionContext(r, txn)

			// Get or Set X-Request-Id
			requestId := r.Header.Get("X-Request-Id")
			if requestId == "" {
				requestId = uuid.New().String()
			}
			ctx := context.WithValue(r.Context(), "trace_id", requestId)
			r = r.WithContext(ctx)

			w.Header().Set("X-Request-Id", requestId)

			// Start timer
			start := time.Now()
			reqBody, _ := io.ReadAll(r.Body)
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			// Add Request Log
			requestStringify := strings.ReplaceAll(string(reqBody), " ", "")
			regexpNewLineAndTab := regexp.MustCompile(`\r?\n`)
			requestStringify = regexpNewLineAndTab.ReplaceAllString(requestStringify, "")
			requestLog := map[string]interface{}{
				"request_id":      requestId,
				"request_url":     r.URL.Path,
				"request_method":  r.Method,
				"request_headers": r.Header,
				"request_payload": requestStringify,
			}
			logger.Info(ctx, fmt.Sprintf("%s %s Request Log", r.Method, r.URL.Path), zap.Any("data", requestLog))

			// Wrap the original writer with our custom writer
			wrappedBody := &bytes.Buffer{}
			wrappedWriter := &ResponseWriter{ResponseWriter: w, body: wrappedBody}

			next.ServeHTTP(wrappedWriter, r)

			// Stop timer
			duration := getDurationInMilliseconds(start)

			go func(requestId string, duration float64, wrappedWriter *ResponseWriter) {
				ctx := context.WithValue(context.Background(), constant.CtxTraceIdKey, requestId)

				// Add Response Log
				responseStringify := strings.ReplaceAll(wrappedWriter.body.String(), " ", "")
				regexpNewLineAndTab := regexp.MustCompile(`\r?\n`)
				responseStringify = regexpNewLineAndTab.ReplaceAllString(responseStringify, "")
				responseLog := map[string]interface{}{
					"response_status": wrappedWriter.Status,
					"response_data":   responseStringify,
					"duration":        duration,
				}
				logger.Info(ctx, fmt.Sprintf("%s %s Response Log", r.Method, r.URL.Path), zap.Any("data", responseLog))
			}(requestId, duration, wrappedWriter)
		})
	}
}
