package converter

import (
	"github.com/ibm-observability/instanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/model/pdata"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.8.0"
)

var _ Converter = (*RuntimeJavaConverter)(nil)

type RuntimeJavaConverter struct{}

func (c *RuntimeJavaConverter) AcceptsMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) bool {
	runtimeAttr, ex := attributes.Get(conventions.AttributeTelemetrySDKLanguage)

	return ex && runtimeAttr.AsString() == conventions.AttributeTelemetrySDKLanguageJava
}

func (c *RuntimeJavaConverter) ConvertMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) []instanaacceptor.PluginPayload {
	processPidAttr, ex := attributes.Get(conventions.AttributeProcessPID)
	if !ex {
		return make([]instanaacceptor.PluginPayload, 0)
	}

	return []instanaacceptor.PluginPayload{
		model.NewJvmRuntimePlugin(
			createJvmSnapshot(attributes, int(processPidAttr.IntVal())),
		),
	}
}

func (c *RuntimeJavaConverter) AcceptsSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) bool {

	runtimeAttr, ex := attributes.Get(conventions.AttributeTelemetrySDKLanguage)

	return ex && runtimeAttr.AsString() == conventions.AttributeTelemetrySDKLanguageJava
}

func (c *RuntimeJavaConverter) ConvertSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) model.Bundle {
	bundle := model.NewBundle()
	processPidAttr, ex := attributes.Get(conventions.AttributeProcessPID)
	if !ex {
		return bundle
	}

	bundle.Metrics.Plugins = append(bundle.Metrics.Plugins, model.NewJvmRuntimePlugin(
		createJvmSnapshot(attributes, int(processPidAttr.IntVal())),
	))

	return bundle
}

func (c *RuntimeJavaConverter) Name() string {
	return "RuntimeJavaConverter"
}

func createJvmSnapshot(attributes pdata.AttributeMap, processPid int) model.JVMProcessData {
	processData := model.JVMProcessData{
		PID: processPid,
	}

	processNameAttr, ex := attributes.Get(conventions.AttributeProcessExecutableName)
	if ex {
		processData.Name = processNameAttr.AsString()
	}

	runtimeNameAttr, ex := attributes.Get(conventions.AttributeProcessRuntimeName)
	if ex {
		processData.JvmVendor = runtimeNameAttr.AsString()
	}

	runtimeVersionAttr, ex := attributes.Get(conventions.AttributeProcessRuntimeVersion)
	if ex {
		processData.JvmVersion = runtimeVersionAttr.AsString()
	}

	return processData
}
