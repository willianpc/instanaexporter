package config

import (
	"errors"

	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
)

const (
	// AttributeInstanaHostID can be used to distinguish multiple hosts' data
	// being processed by a single collector (in a chained scenario)
	AttributeInstanaHostID = "instana.host.id"

	HeaderKey  = "x-instana-key"
	HeaderHost = "x-instana-host"
	HeaderTime = "x-instana-time"
)

// Config defines configuration for logging exporter.
type Config struct {
	config.ExporterSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	AgentEndpoint string `mapstructure:"agent_endpoint"`

	AgentKey string `mapstructure:"agent_key"`

	CustomZone string `mapstructure:"zone"`

	confighttp.HTTPClientSettings `mapstructure:",squash"`

	// LogLevel defines log level of the logging exporter; options are debug, info, warn, error.
	LogLevel zapcore.Level `mapstructure:"loglevel"`
}

var _ config.Exporter = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {

	if cfg.AgentEndpoint == "" {
		return errors.New("no Instana Agent Endpoint set")
	}

	if cfg.AgentKey == "" {
		return errors.New("no Instana Agent key set")
	}

	if cfg.CustomZone == "" {
		return errors.New("no zone set")
	}

	return nil
}
