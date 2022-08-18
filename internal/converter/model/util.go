package model

import (
	"encoding/hex"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func convertTraceId(traceId pcommon.TraceID) string {
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

func convertSpanId(spanId pcommon.SpanID) string {
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

func oTelKindToInstanaKind(otelKind ptrace.SpanKind) (string, bool) {
	switch otelKind {
	case ptrace.SpanKindServer:
		return INSTANA_SPAN_KIND_SERVER, true
	case ptrace.SpanKindClient:
		return INSTANA_SPAN_KIND_CLIENT, false
	case ptrace.SpanKindProducer:
		return INSTANA_SPAN_KIND_PRODUCER, false
	case ptrace.SpanKindConsumer:
		return INSTANA_SPAN_KIND_CONSUMER, true
	case ptrace.SpanKindInternal:
		return INSTANA_SPAN_KIND_INTERNAL, false
	default:
		return "unknown", false
	}
}
