package model

import (
	"fmt"
	"time"

	"go.opentelemetry.io/collector/model/pdata"
)

const (
	OTEL_SPAN_TYPE = "otel"

	INSTANA_SPAN_KIND_SERVER   = "server"
	INSTANA_SPAN_KIND_CLIENT   = "client"
	INSTANA_SPAN_KIND_PRODUCER = "producer"
	INSTANA_SPAN_KIND_CONSUMER = "consumer"
	INSTANA_SPAN_KIND_INTERNAL = "internal"

	INSTANA_DATA_SERVICE      = "service"
	INSTANA_DATA_OPERATION    = "operation"
	INSTANA_DATA_TRACE_STATE  = "trace_state"
	INSTANA_DATA_ERROR        = "error"
	INSTANA_DATA_ERROR_DETAIL = "error_detail"
)

type BatchInfo struct {
	Size int `json:"s"`
}

type FromS struct {
	EntityID string `json:"e"`
	// Serverless agents fields
	Hostless      bool   `json:"hl,omitempty"`
	CloudProvider string `json:"cp,omitempty"`
	// Host agent fields
	HostID string `json:"h,omitempty"`
}

type TraceReference struct {
	TraceID  string `json:"t"`
	ParentID string `json:"p,omitempty"`
}

type OTelSpanData struct {
	Kind           string            `json:"kind"`
	HasTraceParent bool              `json:"tp,omitempty"`
	ServiceName    string            `json:"service"`
	Operation      string            `json:"operation"`
	TraceState     string            `json:"trace_state,omitempty"`
	Tags           map[string]string `json:"tags,omitempty"`
}

type Span struct {
	TraceReference

	SpanID          string          `json:"s"`
	LongTraceID     string          `json:"lt,omitempty"`
	Timestamp       uint64          `json:"ts"`
	Duration        uint64          `json:"d"`
	Name            string          `json:"n"`
	From            *FromS          `json:"f"`
	Batch           *BatchInfo      `json:"b,omitempty"`
	Ec              int             `json:"ec,omitempty"`
	Synthetic       bool            `json:"sy,omitempty"`
	CorrelationType string          `json:"crtp,omitempty"`
	CorrelationID   string          `json:"crid,omitempty"`
	ForeignTrace    bool            `json:"tp,omitempty"`
	Ancestor        *TraceReference `json:"ia,omitempty"`
	Data            OTelSpanData    `json:"data,omitempty"`
}

func ConvertPDataSpanToInstanaSpan(fromS FromS, otelSpan pdata.Span, serviceName string, attributes pdata.AttributeMap) (Span, error) {
	traceId := convertTraceId(otelSpan.TraceID())

	instanaSpan := Span{
		Name:           OTEL_SPAN_TYPE,
		TraceReference: TraceReference{},
		Timestamp:      uint64(otelSpan.StartTimestamp()) / uint64(time.Millisecond),
		Duration:       (uint64(otelSpan.EndTimestamp()) - uint64(otelSpan.StartTimestamp())) / uint64(time.Millisecond),
		Data: OTelSpanData{
			Tags: make(map[string]string),
		},
		From: &fromS,
	}

	if len(traceId) != 32 {
		return Span{}, fmt.Errorf("failed parsing span, length of TraceId should be 32, but got %d", len(traceId))
	}

	instanaSpan.TraceReference.TraceID = traceId[16:32]
	instanaSpan.LongTraceID = traceId

	if !otelSpan.ParentSpanID().IsEmpty() {
		instanaSpan.TraceReference.ParentID = convertSpanId(otelSpan.ParentSpanID())
	}

	instanaSpan.SpanID = convertSpanId(otelSpan.SpanID())

	kind, isEntry := oTelKindToInstanaKind(otelSpan.Kind())
	instanaSpan.Data.Kind = kind

	if !otelSpan.ParentSpanID().IsEmpty() && isEntry {
		instanaSpan.Data.HasTraceParent = true
	}

	instanaSpan.Data.ServiceName = serviceName

	instanaSpan.Data.Operation = otelSpan.Name()

	if otelSpan.TraceState() != pdata.TraceStateEmpty {
		instanaSpan.Data.TraceState = string(otelSpan.TraceState())
	}

	otelSpan.Attributes().Sort().Range(func(k string, v pdata.AttributeValue) bool {
		instanaSpan.Data.Tags[k] = v.AsString()

		return true
	})

	errornous := false
	if otelSpan.Status().Code() == pdata.StatusCodeError {
		errornous = true
		instanaSpan.Data.Tags[INSTANA_DATA_ERROR] = otelSpan.Status().Code().String()
		instanaSpan.Data.Tags[INSTANA_DATA_ERROR_DETAIL] = otelSpan.Status().Message()
	}

	if errornous {
		instanaSpan.Ec = 1
	}

	return instanaSpan, nil
}
