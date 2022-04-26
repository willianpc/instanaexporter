package model

import (
	"encoding/hex"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/model/pdata"
)

func metricNameToCompact(metricName string, attributes pdata.AttributeMap) string {
	if attributes.Len() == 0 {
		return metricName + "{}"
	}

	var labels = []string{}
	attributes.Sort().Range(func(key string, value pdata.AttributeValue) bool {
		labels = append(labels, fmt.Sprintf("%s=\"%s\"", key, value.AsString()))

		return true
	})

	return fmt.Sprintf("%s{%s}", metricName, strings.Join(labels, ","))
}

func convertTraceId(traceId pdata.TraceID) string {
	const byteLength = 16

	bytes := traceId.Bytes()
	traceBytes := make([]byte, 0)

	for (len(traceBytes) + len(bytes)) < byteLength {
		traceBytes = append(traceBytes, 0)
	}

	for _, byte := range bytes {
		traceBytes = append(traceBytes, byte)
	}

	return hex.EncodeToString(traceBytes)
}

func convertSpanId(spanId pdata.SpanID) string {
	const byteLength = 8

	bytes := spanId.Bytes()
	spanBytes := make([]byte, 0)

	for (len(spanBytes) + len(bytes)) < byteLength {
		spanBytes = append(spanBytes, 0)
	}

	for _, byte := range bytes {
		spanBytes = append(spanBytes, byte)
	}

	return hex.EncodeToString(spanBytes)
}

func oTelKindToInstanaKind(otelKind pdata.SpanKind) (string, bool) {
	switch otelKind {
	case pdata.SpanKindServer:
		return INSTANA_SPAN_KIND_SERVER, true
	case pdata.SpanKindClient:
		return INSTANA_SPAN_KIND_CLIENT, false
	case pdata.SpanKindProducer:
		return INSTANA_SPAN_KIND_PRODUCER, false
	case pdata.SpanKindConsumer:
		return INSTANA_SPAN_KIND_CONSUMER, true
	case pdata.SpanKindInternal:
		return INSTANA_SPAN_KIND_INTERNAL, false
	default:
		return "unknown", false
	}
}
