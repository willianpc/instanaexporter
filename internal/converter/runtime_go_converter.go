package converter

import (
	"github.com/ibm-observability/instanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.8.0"
)

var _ Converter = (*RuntimeGoConverter)(nil)

type RuntimeGoConverter struct{}

func (c *RuntimeGoConverter) AcceptsMetrics(attributes pcommon.Map, metricSlice pmetric.MetricSlice) bool {
	runtimeAttr, ex := attributes.Get(conventions.AttributeTelemetrySDKLanguage)

	return ex && runtimeAttr.AsString() == conventions.AttributeTelemetrySDKLanguageGo
}

func (c *RuntimeGoConverter) ConvertMetrics(attributes pcommon.Map, metricSlice pmetric.MetricSlice) []instanaacceptor.PluginPayload {
	processPidAttr, ex := attributes.Get(conventions.AttributeProcessPID)
	if !ex {
		return make([]instanaacceptor.PluginPayload, 0)
	}

	return []instanaacceptor.PluginPayload{
		instanaacceptor.NewGoProcessPluginPayload(
			createGoSnapshot(attributes, int(processPidAttr.IntVal())),
		),
	}
}

func (c *RuntimeGoConverter) AcceptsSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) bool {

	runtimeAttr, ex := attributes.Get(conventions.AttributeTelemetrySDKLanguage)

	return ex && runtimeAttr.AsString() == conventions.AttributeTelemetrySDKLanguageGo
}

func (c *RuntimeGoConverter) ConvertSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) model.Bundle {
	bundle := model.NewBundle()
	processPidAttr, ex := attributes.Get(conventions.AttributeProcessPID)
	if !ex {
		return bundle
	}

	bundle.Metrics.Plugins = append(bundle.Metrics.Plugins, instanaacceptor.NewGoProcessPluginPayload(
		createGoSnapshot(attributes, int(processPidAttr.IntVal())),
	))

	return bundle
}

func (c *RuntimeGoConverter) Name() string {
	return "RuntimeGoConverter"
}

func createGoSnapshot(attributes pcommon.Map, processPid int) instanaacceptor.GoProcessData {
	runtimeInfo := instanaacceptor.RuntimeInfo{}

	processNameAttr, ex := attributes.Get(conventions.AttributeProcessExecutableName)
	if ex {
		runtimeInfo.Name = processNameAttr.AsString()
	}

	runtimeNameAttr, ex := attributes.Get(conventions.AttributeProcessRuntimeName)
	if ex {
		runtimeInfo.Compiler = runtimeNameAttr.AsString()
	}

	runtimeVersionAttr, ex := attributes.Get(conventions.AttributeProcessRuntimeVersion)
	if ex {
		runtimeInfo.Version = runtimeVersionAttr.AsString()
	}

	return instanaacceptor.GoProcessData{
		PID:      processPid,
		Snapshot: &runtimeInfo,
	}
}
