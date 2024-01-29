package response

import (
	"encoding/json"
	"net/http"
)

func SendResponseOK(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := Response{
		Code: HttpStatusOK,
		Data: data,
	}

	return json.NewEncoder(w).Encode(resp)
}

func SendResponseError(w http.ResponseWriter, errType string, errMessage error) error {
	code, statusCode := HttpStatusErrorCode(errType)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response{
		Code:  code,
		Error: errMessage.Error(),
	}

	return json.NewEncoder(w).Encode(resp)
}
