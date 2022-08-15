package otlptext

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type dataBuffer struct {
	buf bytes.Buffer
}

func (b *dataBuffer) logEntry(format string, a ...interface{}) {
	b.buf.WriteString(fmt.Sprintf(format, a...))
	b.buf.WriteString("\n")
}

func (b *dataBuffer) logAttr(label string, value string) {
	b.logEntry("    %-15s: %s", label, value)
}

func (b *dataBuffer) logAttributeMap(label string, am pcommon.Map) {
	if am.Len() == 0 {
		return
	}

	b.logEntry("%s:", label)
	am.Range(func(k string, v pcommon.Value) bool {
		b.logEntry("     -> %s: %s(%s)", k, v.Type().String(), attributeValueToString(v))
		return true
	})
}

func (b *dataBuffer) logInstrumentationLibrary(il pcommon.InstrumentationScope) {
	b.logEntry(
		"InstrumentationLibrary %s %s",
		il.Name(),
		il.Version())
}

func (b *dataBuffer) logMetricDescriptor(md pmetric.Metric) {
	b.logEntry("Descriptor:")
	b.logEntry("     -> Name: %s", md.Name())
	b.logEntry("     -> Description: %s", md.Description())
	b.logEntry("     -> Unit: %s", md.Unit())
	b.logEntry("     -> DataType: %s", md.DataType().String())
}

func (b *dataBuffer) logMetricDataPoints(m pmetric.Metric) {
	switch m.DataType() {
	case pmetric.MetricDataTypeNone:
		return
	case pmetric.MetricDataTypeGauge:
		b.logNumberDataPoints(m.Gauge().DataPoints())
	case pmetric.MetricDataTypeSum:
		data := m.Sum()
		b.logEntry("     -> IsMonotonic: %t", data.IsMonotonic())
		b.logEntry("     -> AggregationTemporality: %s", data.AggregationTemporality().String())
		b.logNumberDataPoints(data.DataPoints())
	case pmetric.MetricDataTypeHistogram:
		data := m.Histogram()
		b.logEntry("     -> AggregationTemporality: %s", data.AggregationTemporality().String())
		b.logDoubleHistogramDataPoints(data.DataPoints())
	case pmetric.MetricDataTypeSummary:
		data := m.Summary()
		b.logDoubleSummaryDataPoints(data.DataPoints())
	}
}

func (b *dataBuffer) logNumberDataPoints(ps pmetric.NumberDataPointSlice) {
	for i := 0; i < ps.Len(); i++ {
		p := ps.At(i)
		b.logEntry("NumberDataPoints #%d", i)
		b.logDataPointAttributes(p.Attributes())

		b.logEntry("StartTimestamp: %s", p.StartTimestamp())
		b.logEntry("Timestamp: %s", p.Timestamp())
		switch p.ValueType() {
		case pmetric.NumberDataPointValueTypeInt:
			b.logEntry("Value: %d", p.IntVal())
		case pmetric.NumberDataPointValueTypeDouble:
			b.logEntry("Value: %f", p.DoubleVal())
		}
	}
}

func (b *dataBuffer) logDoubleHistogramDataPoints(ps pmetric.HistogramDataPointSlice) {
	for i := 0; i < ps.Len(); i++ {
		p := ps.At(i)
		b.logEntry("HistogramDataPoints #%d", i)
		b.logDataPointAttributes(p.Attributes())

		b.logEntry("StartTimestamp: %s", p.StartTimestamp())
		b.logEntry("Timestamp: %s", p.Timestamp())
		b.logEntry("Count: %d", p.Count())
		b.logEntry("Sum: %f", p.Sum())

		bounds := p.ExplicitBounds()
		for i := 0; i < bounds.Len(); i++ {
			b.logEntry("ExplicitBounds #%d: %f", i, bounds.At(i))
		}

		buckets := p.BucketCounts()
		for j := 0; j < buckets.Len(); j++ {
			b.logEntry("Buckets #%d, Count: %d", j, buckets.At(i))
		}
	}
}

func (b *dataBuffer) logDoubleSummaryDataPoints(ps pmetric.SummaryDataPointSlice) {
	for i := 0; i < ps.Len(); i++ {
		p := ps.At(i)
		b.logEntry("SummaryDataPoints #%d", i)
		b.logDataPointAttributes(p.Attributes())

		b.logEntry("StartTimestamp: %s", p.StartTimestamp())
		b.logEntry("Timestamp: %s", p.Timestamp())
		b.logEntry("Count: %d", p.Count())
		b.logEntry("Sum: %f", p.Sum())

		quantiles := p.QuantileValues()
		for i := 0; i < quantiles.Len(); i++ {
			quantile := quantiles.At(i)
			b.logEntry("QuantileValue #%d: Quantile %f, Value %f", i, quantile.Quantile(), quantile.Value())
		}
	}
}

func (b *dataBuffer) logDataPointAttributes(labels pcommon.Map) {
	b.logAttributeMap("Data point attributes", labels)
}

func (b *dataBuffer) logEvents(description string, se ptrace.SpanEventSlice) {
	if se.Len() == 0 {
		return
	}

	b.logEntry("%s:", description)
	for i := 0; i < se.Len(); i++ {
		e := se.At(i)
		b.logEntry("SpanEvent #%d", i)
		b.logEntry("     -> Name: %s", e.Name())
		b.logEntry("     -> Timestamp: %s", e.Timestamp())
		b.logEntry("     -> DroppedAttributesCount: %d", e.DroppedAttributesCount())

		if e.Attributes().Len() == 0 {
			continue
		}
		b.logEntry("     -> Attributes:")
		e.Attributes().Range(func(k string, v pcommon.Value) bool {
			b.logEntry("         -> %s: %s(%s)", k, v.Type().String(), attributeValueToString(v))
			return true
		})
	}
}

func (b *dataBuffer) logLinks(description string, sl ptrace.SpanLinkSlice) {
	if sl.Len() == 0 {
		return
	}

	b.logEntry("%s:", description)

	for i := 0; i < sl.Len(); i++ {
		l := sl.At(i)
		b.logEntry("SpanLink #%d", i)
		b.logEntry("     -> Trace ID: %s", l.TraceID().HexString())
		b.logEntry("     -> ID: %s", l.SpanID().HexString())
		b.logEntry("     -> TraceState: %s", l.TraceState())
		b.logEntry("     -> DroppedAttributesCount: %d", l.DroppedAttributesCount())
		if l.Attributes().Len() == 0 {
			continue
		}
		b.logEntry("     -> Attributes:")
		l.Attributes().Range(func(k string, v pcommon.Value) bool {
			b.logEntry("         -> %s: %s(%s)", k, v.Type().String(), attributeValueToString(v))
			return true
		})
	}
}

func attributeValueToString(av pcommon.Value) string {
	switch av.Type() {
	case pcommon.ValueTypeString:
		return av.StringVal()
	case pcommon.ValueTypeBool:
		return strconv.FormatBool(av.BoolVal())
	case pcommon.ValueTypeDouble:
		return strconv.FormatFloat(av.DoubleVal(), 'f', -1, 64)
	case pcommon.ValueTypeInt:
		return strconv.FormatInt(av.IntVal(), 10)
	case pcommon.ValueTypeSlice:
		return attributeValueArrayToString(av.SliceVal())
	case pcommon.ValueTypeMap:
		return attributeMapToString(av.MapVal())
	default:
		return fmt.Sprintf("<Unknown OpenTelemetry attribute value type %q>", av.Type())
	}
}

func attributeValueArrayToString(av pcommon.Slice) string {
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < av.Len(); i++ {
		if i < av.Len()-1 {
			fmt.Fprintf(&b, "%s, ", attributeValueToString(av.At(i)))
		} else {
			b.WriteString(attributeValueToString(av.At(i)))
		}
	}

	b.WriteByte(']')
	return b.String()
}

func attributeMapToString(av pcommon.Map) string {
	var b strings.Builder
	b.WriteString("{\n")

	av.Sort().Range(func(k string, v pcommon.Value) bool {
		fmt.Fprintf(&b, "     -> %s: %s(%s)\n", k, v.Type(), v.AsString())
		return true
	})
	b.WriteByte('}')
	return b.String()
}
