package instanaexporter

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	instanaConfig "github.com/ibm-observability/instanaexporter/config"
)

const (
	// The value of "type" key in configuration.
	typeStr = "instana"
	// The stability level of the exporter.
	stability = component.StabilityLevelBeta
)

//NewFactory creates an Instana exporter factory
func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesExporter(createTracesExporter, stability),
		component.WithMetricsExporter(createMetricsExporter, stability),
	)
}

// createDefaultConfig creates the default exporter configuration
func createDefaultConfig() config.Exporter {
	return &instanaConfig.Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		LogLevel:         zapcore.InfoLevel,
		HTTPClientSettings: confighttp.HTTPClientSettings{
			Endpoint: "",
			Timeout:  30 * time.Second,
			Headers:  map[string]string{},
			// We almost read 0 bytes, so no need to tune ReadBufferSize.
			WriteBufferSize: 512 * 1024,
		},
	}
}

// createTracesExporter creates a trace exporter based on this configuration
func createTracesExporter(ctx context.Context, set component.ExporterCreateSettings, config config.Exporter) (component.TracesExporter, error) {
	cfg := config.(*instanaConfig.Config)

	exporterLogger, err := createLogger(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	instanaExporter, err := newInstanaExporter(exporterLogger, cfg, set)
	if err != nil {
		cancel()
		return nil, err
	}

	return exporterhelper.NewTracesExporterWithContext(
		ctx,
		set,
		config,
		instanaExporter.pushConvertedTraces,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(instanaExporter.start),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
		exporterhelper.WithShutdown(func(context.Context) error {
			cancel()
			return nil
		}),
	)
}

// createMetricsExporter creates a metrics exporter based on this configuration
func createMetricsExporter(ctx context.Context, set component.ExporterCreateSettings, config config.Exporter) (component.MetricsExporter, error) {
	cfg := config.(*instanaConfig.Config)

	exporterLogger, err := createLogger(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	instanaExporter, err := newInstanaExporter(exporterLogger, cfg, set)
	if err != nil {
		cancel()
		return nil, err
	}

	return exporterhelper.NewMetricsExporterWithContext(
		ctx,
		set,
		config,
		instanaExporter.pushConvertedMetrics,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(instanaExporter.start),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
		exporterhelper.WithShutdown(func(context.Context) error {
			cancel()
			return nil
		}),
	)
}

// createLogger creates a logger for logging trace and errors
func createLogger(cfg *instanaConfig.Config) (*zap.Logger, error) {
	// We take development config as the base since it matches the purpose
	// of logging exporter being used for debugging reasons (so e.g. console encoder)
	conf := zap.NewDevelopmentConfig()
	conf.Level = zap.NewAtomicLevelAt(cfg.LogLevel)

	logginglogger, err := conf.Build()
	if err != nil {
		return nil, err
	}
	return logginglogger, nil
}
