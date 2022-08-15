package converter

import (
	"github.com/ibm-observability/instanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type Converter interface {
	AcceptsMetrics(attributes pcommon.Map, metricSlice pmetric.MetricSlice) bool
	ConvertMetrics(attributes pcommon.Map, metricSlice pmetric.MetricSlice) []instanaacceptor.PluginPayload
	AcceptsSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) bool
	ConvertSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) model.Bundle
	Name() string
}
