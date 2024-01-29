package response

type Response struct {
	Code  string      `json:"code,omitempty"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
