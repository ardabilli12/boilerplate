package response

const (
	HttpStatusOK string = "00"

	HttpStatusErrorInternal        string = "99"
	HttpStatusErrorDatabase        string = "98"
	HttpStatusErrorThirdParty      string = "50"
	HttpStatusErrorRequest         string = "40"
	HttpStatusErrorUnauthorized    string = "41"
	HttpStatusErrorNotFound        string = "44"
	HttpStatusErrorDuplicatedCheck string = "49"
)

const (
	HttpErrInternal     string = "ERROR_INTERNAL"
	HttpErrDatabase     string = "ERROR_DATABASE"
	HttpErrThirdParty   string = "ERROR_THIRD_PARTY"
	HttpErrRequest      string = "ERROR_REQUEST"
	HttpErrUnauthorized string = "ERROR_UNAUTHORIZED"
	HttpErrNotFound     string = "ERROR_NOT_FOUND"
	HttpErrDupCheck     string = "ERROR_DUPLICATE_CHECK"
)
