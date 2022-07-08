package converter

import (
	"fmt"

	"github.com/ibm-observability/instanaexporter/config"
	"github.com/ibm-observability/instanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/model/pdata"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.8.0"
)

const (
	OTEL_SPAN_TYPE = "otel"
)

var _ Converter = (*SpanConverter)(nil)

type SpanConverter struct{}

func (c *SpanConverter) AcceptsMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) bool {
	return false
}

func (c *SpanConverter) ConvertMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) []instanaacceptor.PluginPayload {
	return make([]instanaacceptor.PluginPayload, 0)
}

func (c *SpanConverter) AcceptsSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) bool {

	return true
}

func (c *SpanConverter) ConvertSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) model.Bundle {
	bundle := model.NewBundle()
	spans := make([]model.Span, 0)

	fromS := model.FromS{}

	hostIdValue, ex := attributes.Get(config.AttributeInstanaHostID)
	if !ex {
		fromS.HostID = "unknown-host-id"
	} else {
		fromS.HostID = hostIdValue.AsString()
	}

	processIdValue, ex := attributes.Get(conventions.AttributeProcessPID)
	if !ex {
		fromS.EntityID = "unknown-process-id"
	} else {
		fromS.EntityID = processIdValue.AsString()
	}

	serviceName := ""
	serviceNameValue, ex := attributes.Get(conventions.AttributeServiceName)
	if ex {
		serviceName = serviceNameValue.AsString()
	}

	for i := 0; i < spanSlice.Len(); i++ {
		otelSpan := spanSlice.At(i)

		instanaSpan, err := model.ConvertPDataSpanToInstanaSpan(fromS, otelSpan, serviceName, attributes)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}

		spans = append(spans, instanaSpan)
	}

	bundle.Spans = spans

	return bundle
}

func (c *SpanConverter) Name() string {
	return "SpanConverter"
}
