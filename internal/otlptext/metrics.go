package otlptext // import "go.opentelemetry.io/collector/internal/otlptext"

import (
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// NewTextMetricsMarshaler returns a serializer.MetricsMarshaler to encode to OTLP text bytes.
func NewTextMetricsMarshaler() pmetric.Marshaler {
	return textMetricsMarshaler{}
}

type textMetricsMarshaler struct{}

// MarshalMetrics pdata.Metrics to OTLP text.
func (textMetricsMarshaler) MarshalMetrics(md pmetric.Metrics) ([]byte, error) {
	buf := dataBuffer{}
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		buf.logEntry("ResourceMetrics #%d", i)
		rm := rms.At(i)
		buf.logEntry("Resource SchemaURL: %s", rm.SchemaUrl())
		buf.logAttributeMap("Resource attributes", rm.Resource().Attributes())
		ilms := rm.ScopeMetrics()
		for j := 0; j < ilms.Len(); j++ {
			buf.logEntry("InstrumentationLibraryMetrics #%d", j)
			ilm := ilms.At(j)
			buf.logEntry("InstrumentationLibraryMetrics SchemaURL: %s", ilm.SchemaUrl())
			buf.logInstrumentationLibrary(ilm.Scope())
			metrics := ilm.Metrics()
			for k := 0; k < metrics.Len(); k++ {
				buf.logEntry("Metric #%d", k)
				metric := metrics.At(k)
				buf.logMetricDescriptor(metric)
				buf.logMetricDataPoints(metric)
			}
		}
	}

	return buf.buf.Bytes(), nil
}
