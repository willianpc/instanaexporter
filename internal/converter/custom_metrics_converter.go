package converter

import (
	"github.com/ibm-observability/instanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/model/pdata"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.8.0"
)

var _ Converter = (*CustomMetricsConverter)(nil)

type CustomMetricsConverter struct{}

func (c *CustomMetricsConverter) AcceptsMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) bool {
	return true
}

func (c *CustomMetricsConverter) ConvertMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) []instanaacceptor.PluginPayload {
	// return early if no metrics contained
	if metricSlice.Len() == 0 {
		return make([]instanaacceptor.PluginPayload, 0)
	}

	entityID := "h"
	if pidAttribute, ex := attributes.Get(conventions.AttributeProcessPID); ex {
		entityID = pidAttribute.AsString()
	}

	metricData := model.NewOpenTelemetryCustomMetricsData()

	for i := 0; i < metricSlice.Len(); i++ {
		metric := metricSlice.At(i)

		metricData.AppendMetric(metric)
	}

	metricsPayload := model.NewOpenTelemetryMetricsPluginPayload(entityID, metricData)

	return []instanaacceptor.PluginPayload{
		metricsPayload,
	}
}

func (c *CustomMetricsConverter) AcceptsSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) bool {

	return false
}

func (c *CustomMetricsConverter) ConvertSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) model.Bundle {

	return model.NewBundle()
}

func (c *CustomMetricsConverter) Name() string {
	return "CustomMetricsConverter"
}
