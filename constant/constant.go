package constant

type constantKey string

const (
	// ctxTraceIdKey is the context key for trace id
	CtxTraceIdKey constantKey = "trace_id"

	EnvironmentDevelopment = "development"
	EnvironmentLocal       = "local"
	EnvironmentStaging     = "staging"
)

const (
	DefaultCurrencyIDR string = "IDR"
)
