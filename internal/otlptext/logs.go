package otlptext // import "go.opentelemetry.io/collector/internal/otlptext"

import (
	"go.opentelemetry.io/collector/pdata/plog"
)

// NewTextLogsMarshaler returns a serializer.LogsMarshaler to encode to OTLP text bytes.
func NewTextLogsMarshaler() plog.Marshaler {
	return textLogsMarshaler{}
}

type textLogsMarshaler struct{}

// MarshalLogs pdata.Logs to OTLP text.
func (textLogsMarshaler) MarshalLogs(ld plog.Logs) ([]byte, error) {
	buf := dataBuffer{}
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		buf.logEntry("ResourceLog #%d", i)
		rl := rls.At(i)
		buf.logEntry("Resource SchemaURL: %s", rl.SchemaUrl())
		buf.logAttributeMap("Resource labels", rl.Resource().Attributes())
		ills := rl.ScopeLogs()
		for j := 0; j < ills.Len(); j++ {
			buf.logEntry("InstrumentationLibraryLogs #%d", j)
			ils := ills.At(j)
			buf.logEntry("InstrumentationLibraryLogs SchemaURL: %s", ils.SchemaUrl())
			buf.logInstrumentationLibrary(ils.Scope())

			logs := ils.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				buf.logEntry("LogRecord #%d", k)
				lr := logs.At(k)
				buf.logEntry("Timestamp: %s", lr.Timestamp())
				buf.logEntry("Severity: %s", lr.SeverityText())

				// TODO (hickeyma): Name() removed between 0.48 and 0.58 of Otel API
				// Need to refactor with latest means to get the name
				//buf.logEntry("ShortName: %s", lr.Name())

				buf.logEntry("Body: %s", attributeValueToString(lr.Body()))
				buf.logAttributeMap("Attributes", lr.Attributes())
				buf.logEntry("Trace ID: %s", lr.TraceID().HexString())
				buf.logEntry("Span ID: %s", lr.SpanID().HexString())
				buf.logEntry("Flags: %d", lr.Flags())
			}
		}
	}

	return buf.buf.Bytes(), nil
}
