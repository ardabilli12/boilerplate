package response

import (
	"net/http"
)

func HttpStatusErrorCode(errType string) (string, int) {
	switch errType {
	case HttpErrNotFound:
		return HttpStatusErrorNotFound, http.StatusNotFound
	case HttpErrUnauthorized:
		return HttpStatusErrorUnauthorized, http.StatusUnauthorized
	case HttpErrDupCheck:
		return HttpStatusErrorDuplicatedCheck, http.StatusConflict
	case HttpErrRequest:
		return HttpStatusErrorRequest, http.StatusBadRequest
	default:
		return HttpStatusErrorInternal, http.StatusInternalServerError
	}
}
