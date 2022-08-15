package converter

import (
	"fmt"

	"github.com/ibm-observability/instanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

var _ Converter = (*ConvertAllConverter)(nil)

type ConvertAllConverter struct {
	converters []Converter
	logger     *zap.Logger
}

func (c *ConvertAllConverter) AcceptsMetrics(attributes pcommon.Map, metricSlice pmetric.MetricSlice) bool {
	return true
}

func (c *ConvertAllConverter) ConvertMetrics(attributes pcommon.Map, metricSlice pmetric.MetricSlice) []instanaacceptor.PluginPayload {
	plugins := make([]instanaacceptor.PluginPayload, 0)

	for i := 0; i < len(c.converters); i++ {
		if !c.converters[i].AcceptsMetrics(attributes, metricSlice) {
			c.logger.Debug(fmt.Sprintf("Converter %s didnt Accept", c.converters[i].Name()))

			continue
		}

		plugins = append(plugins, c.converters[i].ConvertMetrics(attributes, metricSlice)...)
	}

	return plugins
}

func (c *ConvertAllConverter) AcceptsSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) bool {
	return true
}

func (c *ConvertAllConverter) ConvertSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) model.Bundle {
	bundle := model.NewBundle()

	for i := 0; i < len(c.converters); i++ {
		if !c.converters[i].AcceptsSpans(attributes, spanSlice) {
			c.logger.Debug(fmt.Sprintf("Converter %s didnt Accept", c.converters[i].Name()))

			continue
		}

		converterBundle := c.converters[i].ConvertSpans(attributes, spanSlice)
		if len(converterBundle.Metrics.Plugins) > 0 {
			bundle.Metrics.Plugins = append(bundle.Metrics.Plugins, converterBundle.Metrics.Plugins...)
		}

		if len(converterBundle.Spans) > 0 {
			bundle.Spans = append(bundle.Spans, converterBundle.Spans...)
		}
	}

	return bundle
}

func (c *ConvertAllConverter) Name() string {
	return "ConvertAllConverter"
}

func NewConvertAllConverter(logger *zap.Logger) Converter {

	return &ConvertAllConverter{
		converters: []Converter{
			&DockerContainerMetricConverter{},
			&HostMetricConverter{},
			&ProcessMetricConverter{},
			&CustomMetricsConverter{},
			&CollectorMetricsConverter{},

			// Runtimes
			&RuntimeGoConverter{},
			&RuntimeJavaConverter{},
			&RuntimePythonConverter{},

			&SpanConverter{},
		},
		logger: logger,
	}
}
