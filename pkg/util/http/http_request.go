package httputil

import (
	"boilerplate-service/pkg/newRelicExt"
	"bytes"
	"context"

	"encoding/json"
	"io"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func RequestHitAPI(
	ctx context.Context,
	method string,
	uri string,
	data interface{},
	header map[string]string,
) (
	res []byte,
	code int,
	err error,
) {
	segment := newRelicExt.
		GetTxnFromCtx(ctx).
		StartSegment("util/http_request.go/RequestHitAPI")
	defer segment.End()

	httpClient := &http.Client{}
	httpClient.Transport = newrelic.NewRoundTripper(httpClient.Transport)
	request, err := assertTypeRequest(data, method, uri)
	if err != nil {
		return res, code, err
	}

	for k, v := range header {
		request.Header.Add(k, v)
	}

	request.Header.Set("Content-type", "application/json")

	response, err := httpClient.Do(request)
	if err != nil {
		return res, code, err
	}

	defer response.Body.Close()

	code = response.StatusCode

	res, err = io.ReadAll(response.Body)
	if err != nil {
		return res, code, err
	}

	if isHttpError := code != http.StatusOK && code != http.StatusCreated; isHttpError {
		var errRes map[string]interface{}
		err := json.Unmarshal(res, &errRes)
		return res, code, err
	}

	return res, code, err
}

func assertTypeRequest(data interface{}, method string, uri string) (request *http.Request, err error) {
	if data == nil {
		request, err = http.NewRequest(method, uri, nil)
		return
	}

	paramReq, _ := json.Marshal(data)
	request, err = http.NewRequest(method, uri, bytes.NewBuffer(paramReq))

	return
}
