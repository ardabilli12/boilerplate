package newRelicExt

import (
	"boilerplate-service/constant"
	"boilerplate-service/pkg/logger"
	"context"
	"errors"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type INewRelicExt interface {
	App() *newrelic.Application

	RecordCustomEvent(eventType string, params map[string]interface{})
	RecordCustomMetric(name string, value float64)
	Shutdown(timeout time.Duration)
	StartTransaction(name string, opts ...newrelic.TraceOption) *newrelic.Transaction
}

type Config struct {
	LicenseKey  string
	Environment string
	ServiceName string

	Logger logger.ILogger
}

type newRelicExt struct {
	app *newrelic.Application
}

const (
	defaultSlowQueryThreshold = 5 * time.Second
)

func New(config Config) (INewRelicExt, error) {
	if config.LicenseKey == "" {
		return nil, errors.New("license key empty")
	}

	options := []newrelic.ConfigOption{}

	options = append(options, newrelic.ConfigLicense(config.LicenseKey))
	options = append(options, newrelic.ConfigAppName(config.ServiceName))

	options = append(options, newrelic.ConfigInfoLogger(os.Stdout))
	options = append(options, newrelic.ConfigEnabled(true))
	options = append(options, newrelic.ConfigDistributedTracerEnabled(true))
	options = append(options, newrelic.ConfigAppLogEnabled(true))
	options = append(options, newrelic.ConfigCodeLevelMetricsEnabled(true))

	if config.Logger != nil {
		options = append(options, nrzap.ConfigLogger(config.Logger.GetLogger().Named("newrelic")))
	}

	// Optional: add additional changes to your configuration via a config function:
	options = append(options, func(cfg *newrelic.Config) {
		cfg.Labels["env"] = config.Environment

		// If not development
		if config.Environment != constant.EnvironmentLocal && config.Environment != constant.EnvironmentDevelopment {
			cfg.SpanEvents.Enabled = true
			cfg.SpanEvents.Attributes.Exclude = excludingAttrSpans()

			cfg.DatastoreTracer.InstanceReporting.Enabled = true
			cfg.DatastoreTracer.DatabaseNameReporting.Enabled = true
			cfg.DatastoreTracer.QueryParameters.Enabled = true
			cfg.DatastoreTracer.SlowQuery.Enabled = true
			cfg.DatastoreTracer.SlowQuery.Threshold = defaultSlowQueryThreshold

			cfg.RuntimeSampler.Enabled = true
		}
	})

	app, err := newrelic.NewApplication(options...)
	return &newRelicExt{app}, err
}

func (n *newRelicExt) App() *newrelic.Application {
	return n.app
}

func (n *newRelicExt) RecordCustomEvent(eventType string, params map[string]interface{}) {
	n.app.RecordCustomEvent(eventType, params)
}

func (n *newRelicExt) RecordCustomMetric(name string, value float64) {
	n.app.RecordCustomMetric(name, value)
}

func (n *newRelicExt) Shutdown(timeout time.Duration) {
	n.app.Shutdown(timeout)
}

func (n *newRelicExt) StartTransaction(name string, opts ...newrelic.TraceOption) *newrelic.Transaction {
	return n.app.StartTransaction(name, opts...)
}

func GetTxnFromCtx(ctx context.Context) *newrelic.Transaction {
	if txn := ctx.Value(constant.CtxNewRelicTxnKey); txn != nil {
		return txn.(*newrelic.Transaction)
	}
	return nil
}

// excludingAttrSpans returns a list of attributes to exclude to shown in NR spans
// By service requirements
//
// No parameters.
// Returns a slice of strings.
func excludingAttrSpans() []string {
	return []string{
		"password",
	}
}
