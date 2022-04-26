package converter

import (
	"github.com/ibm-observability/ibminstanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/model/pdata"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.8.0"
)

var _ Converter = (*RuntimeGoConverter)(nil)

type RuntimeGoConverter struct{}

func (c *RuntimeGoConverter) AcceptsMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) bool {
	runtimeAttr, ex := attributes.Get(conventions.AttributeTelemetrySDKLanguage)

	return ex && runtimeAttr.AsString() == conventions.AttributeTelemetrySDKLanguageGo
}

func (c *RuntimeGoConverter) ConvertMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) []instanaacceptor.PluginPayload {
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

func (c *RuntimeGoConverter) AcceptsSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) bool {

	runtimeAttr, ex := attributes.Get(conventions.AttributeTelemetrySDKLanguage)

	return ex && runtimeAttr.AsString() == conventions.AttributeTelemetrySDKLanguageGo
}

func (c *RuntimeGoConverter) ConvertSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) model.Bundle {
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

func createGoSnapshot(attributes pdata.AttributeMap, processPid int) instanaacceptor.GoProcessData {
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
