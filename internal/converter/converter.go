package converter

import (
	"github.com/ibm-observability/ibminstanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/model/pdata"
)

type Converter interface {
	AcceptsMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) bool
	ConvertMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) []instanaacceptor.PluginPayload
	AcceptsSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) bool
	ConvertSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) model.Bundle
	Name() string
}
