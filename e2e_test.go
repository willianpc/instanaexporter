package instanaexporter

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/ibm-observability/instanaexporter/config"
	"github.com/ibm-observability/instanaexporter/internal/converter"
	"github.com/ibm-observability/instanaexporter/internal/converter/model"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.8.0"
)

type SpanOptions struct {
	TraceId        [16]byte
	SpanId         [8]byte
	ParentId       [8]byte
	Error          string
	StartTimestamp time.Duration
	EndTimestamp   time.Duration
}

func setupSpan(span *ptrace.Span, opts SpanOptions) {
	var empty16 [16]byte
	var empty8 [8]byte

	now := time.Now().UnixMilli()

	traceId := opts.TraceId
	spanId := opts.SpanId
	parentId := opts.ParentId
	startTime := opts.StartTimestamp
	endTime := opts.EndTimestamp

	if bytes.Equal(traceId[:], empty16[:]) {
		traceId = generateTraceId()
	}

	if bytes.Equal(spanId[:], empty8[:]) {
		spanId = generateSpanId()
	}

	if startTime == time.Second*0 {
		startTime = time.Duration(now)
	}

	if endTime == time.Second*0 {
		endTime = startTime + 1000
	}

	if opts.Error != "" {
		span.Status().SetCode(ptrace.StatusCodeError)
		span.Status().SetMessage(opts.Error)
	}

	if !bytes.Equal(parentId[:], empty8[:]) {
		span.SetParentSpanID(pcommon.NewSpanID(parentId))
	}

	span.SetStartTimestamp(pcommon.Timestamp(startTime * 1e6))
	span.SetEndTimestamp(pcommon.Timestamp(endTime * 1e6))

	span.SetSpanID(pcommon.NewSpanID(spanId))
	span.SetKind(ptrace.SpanKindClient)
	span.SetName("my_operation")
	span.SetTraceState(ptrace.TraceStateEmpty)
	span.SetTraceID(pcommon.NewTraceID(traceId))

	// adding attributes (tags in the instana side)
	span.Attributes().Insert("some_key", pcommon.NewValueBool(true))
}

func generateAttrs() pcommon.Map {
	rawmap := map[string]interface{}{
		"some_boolean_key": true,
		"custom_attribute": "ok",
		// test non empty pid
		conventions.AttributeProcessPID: "1234",
		// test non empty service name
		conventions.AttributeServiceName: "myservice",
		// test non empty instana host id
		config.AttributeInstanaHostID: "myhost1",
	}

	attrs := pcommon.NewMapFromRaw(rawmap)
	attrs.InsertBool("itistrue", true)

	return attrs
}

func validateInstanaSpanBasics(sp model.Span, t *testing.T) {
	if sp.SpanID == "" {
		t.Error("expected span id not to be empty")
	}

	if sp.TraceID == "" {
		t.Error("expected trace id not to be empty")
	}

	if sp.Name != "otel" {
		t.Errorf("expected span name to be 'otel' but received '%v'", sp.Name)
	}

	if sp.Timestamp <= 0 {
		t.Errorf("expected timestamp to be provided but received %v", sp.Timestamp)
	}

	if sp.Duration <= 0 {
		t.Errorf("expected duration to be provided but received %v", sp.Duration)
	}
}

func validateBundle(jsonData []byte, t *testing.T, fn func(model.Span, *testing.T)) {
	var bundle model.Bundle

	err := json.Unmarshal(jsonData, &bundle)

	if err != nil {
		t.Fatal(err)
	}

	if len(bundle.Spans) == 0 {
		t.Log("bundle contains no spans")
		return
	}

	for _, span := range bundle.Spans {
		fn(span, t)
	}
}

func TestSpanBasics(t *testing.T) {
	spanSlice := ptrace.NewSpanSlice()

	sp1 := spanSlice.AppendEmpty()

	setupSpan(&sp1, SpanOptions{})

	attrs := generateAttrs()
	conv := converter.SpanConverter{}
	bundle := conv.ConvertSpans(attrs, spanSlice)
	data, _ := json.MarshalIndent(bundle, "", "  ")

	validateBundle(data, t, func(sp model.Span, t *testing.T) {
		validateInstanaSpanBasics(sp, t)
		validateSpanError(sp, false, t)
	})
}

func TestSpanCorrelation(t *testing.T) {
	spanSlice := ptrace.NewSpanSlice()

	sp1 := spanSlice.AppendEmpty()
	setupSpan(&sp1, SpanOptions{})

	sp2 := spanSlice.AppendEmpty()
	setupSpan(&sp2, SpanOptions{
		ParentId: sp1.SpanID().Bytes(),
	})

	sp3 := spanSlice.AppendEmpty()
	setupSpan(&sp3, SpanOptions{
		ParentId: sp2.SpanID().Bytes(),
	})

	sp4 := spanSlice.AppendEmpty()
	setupSpan(&sp4, SpanOptions{
		ParentId: sp1.SpanID().Bytes(),
	})

	attrs := generateAttrs()
	conv := converter.SpanConverter{}
	bundle := conv.ConvertSpans(attrs, spanSlice)
	data, _ := json.MarshalIndent(bundle, "", "  ")

	spanIdList := make(map[string]bool)

	validateBundle(data, t, func(sp model.Span, t *testing.T) {
		validateInstanaSpanBasics(sp, t)
		validateSpanError(sp, false, t)

		spanIdList[sp.SpanID] = true

		if sp.ParentID != "" && !spanIdList[sp.ParentID] {
			t.Errorf("span %v expected to have parent id %v", sp.SpanID, sp.ParentID)
		}
	})
}
func TestSpanWithError(t *testing.T) {
	spanSlice := ptrace.NewSpanSlice()

	sp1 := spanSlice.AppendEmpty()
	setupSpan(&sp1, SpanOptions{
		Error: "some error",
	})

	attrs := generateAttrs()
	conv := converter.SpanConverter{}
	bundle := conv.ConvertSpans(attrs, spanSlice)
	data, _ := json.MarshalIndent(bundle, "", "  ")

	validateBundle(data, t, func(sp model.Span, t *testing.T) {
		validateInstanaSpanBasics(sp, t)
		validateSpanError(sp, true, t)
	})
}

func validateSpanError(sp model.Span, shouldHaveError bool, t *testing.T) {
	if shouldHaveError {
		if sp.Ec <= 0 {
			t.Error("expected span to have errors (ec = 1)")
		}

		if sp.Data.Tags[model.INSTANA_DATA_ERROR] == "" {
			t.Error("expected data.error to exist")
		}

		if sp.Data.Tags[model.INSTANA_DATA_ERROR_DETAIL] == "" {
			t.Error("expected data.error_detail to exist")
		}

		return
	}

	if sp.Ec > 0 {
		t.Error("expected span not to have errors (ec = 0)")
	}

	if sp.Data.Tags[model.INSTANA_DATA_ERROR] != "" {
		t.Error("expected data.error to be empty")
	}

	if sp.Data.Tags[model.INSTANA_DATA_ERROR_DETAIL] != "" {
		t.Error("expected data.error_detail to be empty")
	}
}

func generateTraceId() (data [16]byte) {
	rand.Read(data[:])

	return data
}

func generateSpanId() (data [8]byte) {
	rand.Read(data[:])

	return data
}
