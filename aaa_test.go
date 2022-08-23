package instanaexporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/ibm-observability/instanaexporter/config"
	instanaConfig "github.com/ibm-observability/instanaexporter/config"
	"github.com/ibm-observability/instanaexporter/internal/converter"
	"github.com/ibm-observability/instanaexporter/internal/converter/model"

	// "go.opentelemetry.io/collector/model/pdata"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.8.0"
)

func setupSpan(span *ptrace.Span) {
	traceId := generateTraceId()
	now := time.Now()

	nowMinus := now.Sub(time.UnixMilli(1000))

	span.SetStartTimestamp(pcommon.Timestamp(nowMinus))
	span.SetEndTimestamp(pcommon.Timestamp(now.UnixMilli()))
	span.SetParentSpanID(pcommon.NewSpanID([8]byte{0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf}))
	span.SetSpanID(pcommon.NewSpanID([8]byte{0xd, 0xa, 0xd, 0xa, 0xf, 0xe, 0xe, 0xb}))
	span.SetKind(ptrace.SpanKindClient)
	span.SetName("span_name")
	span.SetTraceState(ptrace.TraceStateEmpty)
	// span.SetTraceID(pcommon.NewTraceID([16]byte{0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf, 0xd, 0xa, 0xd, 0xa, 0xf, 0xe, 0xe, 0xb}))
	span.SetTraceID(pcommon.NewTraceID(traceId))

	// adding attributes (tags in the instana side)
	span.Attributes().Insert("some_key", pcommon.NewValueBool(true))
}

func generateAttrs() pcommon.Map {
	rawmap := map[string]interface{}{
		"aaa":              true,
		"custom_attribute": "ok",
		// test non empty pid
		conventions.AttributeProcessPID: "bla",
		// test non empty service name
		conventions.AttributeServiceName: "ble",
		// test non empty instana host id
		config.AttributeInstanaHostID: "blu",
	}

	attrs := pcommon.NewMapFromRaw(rawmap)
	attrs.InsertBool("itistrue", true)

	return attrs
}

func validateInstanaSpan(sp model.Span, t *testing.T) {
	if sp.SpanID == "" {
		t.Error("span id must not be empty")
	}

	if sp.TraceID == "" {
		t.Error("trace id must not be empty")
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

func validateBundle(jsonData []byte, t *testing.T) {
	var bundle model.Bundle

	err := json.Unmarshal(jsonData, &bundle)

	if err != nil {
		t.Fatal(err)
	}

	if len(bundle.Spans) == 1 {
		t.Log("bundle contains no spans")
		return
	}

	for _, span := range bundle.Spans {
		validateInstanaSpan(span, t)
	}
}

func TestSpanConvert(t *testing.T) {
	spanSlice := ptrace.NewSpanSlice()

	sp1 := spanSlice.AppendEmpty()
	setupSpan(&sp1)

	sp2 := spanSlice.AppendEmpty()
	setupSpan(&sp2)

	attrs := generateAttrs()

	conv := converter.SpanConverter{}

	bundle := conv.ConvertSpans(attrs, spanSlice)

	data, _ := json.MarshalIndent(bundle, "", "  ")

	validateBundle(data, t)

	// endpoint := "https://pink.instana.rocks/serverless/bundle"

	// headers := map[string]string{
	// 	instanaConfig.HeaderKey:  "agent_key",
	// 	instanaConfig.HeaderHost: "", //host id
	// 	instanaConfig.HeaderTime: "0",
	// }

	// httpRequest(endpoint, data, headers)
}

func httpRequest(url string, data []byte, header map[string]string) {
	// req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(data[:100]))
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	userAgent := fmt.Sprintf("%s/%s (%s/%s)", "Some description", "1.0", runtime.GOOS, runtime.GOARCH)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	for name, value := range header {
		req.Header.Set(name, value)
	}

	req.Header.Set(instanaConfig.HeaderKey, "2Zykc2m_RiKJnVE-TNNdrA")
	req.Header.Set(instanaConfig.HeaderHost, "")
	req.Header.Set(instanaConfig.HeaderTime, "0")

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Request is successful.
		fmt.Println("success", resp.StatusCode, resp.Status)
	} else {
		fmt.Println("not a success", resp)
	}
}

func generateTraceId() (data [16]byte) {
	rand.Read(data[:])

	return data
}
