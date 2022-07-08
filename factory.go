package instanaexporter

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"

	instanaConfig "github.com/ibm-observability/instanaexporter/config"
)

const (
	// The value of "type" key in configuration.
	typeStr = "instana"
)

func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesExporter(createTracesExporter),
		component.WithMetricsExporter(createMetricsExporter),
	)
}

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

func createTracesExporter(_ context.Context, set component.ExporterCreateSettings, config config.Exporter) (component.TracesExporter, error) {
	cfg := config.(*instanaConfig.Config)

	exporterLogger, err := createLogger(cfg)
	if err != nil {
		return nil, err
	}

	return newTracesExporter(config, cfg, exporterLogger, set)
}

func createMetricsExporter(_ context.Context, set component.ExporterCreateSettings, config config.Exporter) (component.MetricsExporter, error) {
	cfg := config.(*instanaConfig.Config)

	exporterLogger, err := createLogger(cfg)
	if err != nil {
		return nil, err
	}

	return newMetricsExporter(config, cfg, exporterLogger, set)
}

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
