package converter

import (
	"github.com/ibm-observability/instanaexporter/internal/converter/model"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type Converter interface {
	AcceptsSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) bool
	ConvertSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) model.Bundle
	Name() string
}
