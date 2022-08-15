package converter

import (
	"strings"

	"github.com/ibm-observability/instanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.8.0"
)

var _ Converter = (*ProcessMetricConverter)(nil)

type ProcessMetricConverter struct{}

func (c *ProcessMetricConverter) AcceptsMetrics(attributes pcommon.Map, metricSlice pmetric.MetricSlice) bool {
	_, ex := attributes.Get(conventions.AttributeProcessPID)

	return ex
}

func (c *ProcessMetricConverter) ConvertMetrics(attributes pcommon.Map, metricSlice pmetric.MetricSlice) []instanaacceptor.PluginPayload {
	plugins := make([]instanaacceptor.PluginPayload, 0)

	processPidAttr, ex := attributes.Get(conventions.AttributeProcessPID)
	if !ex {
		return plugins
	}

	processData := createProcessData(attributes, int(processPidAttr.IntVal()))

	plugins = append(plugins, instanaacceptor.NewProcessPluginPayload(processPidAttr.AsString(), processData))

	return plugins
}

func (c *ProcessMetricConverter) AcceptsSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) bool {
	_, ex := attributes.Get(conventions.AttributeProcessPID)

	return ex
}

func (c *ProcessMetricConverter) ConvertSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) model.Bundle {

	bundle := model.NewBundle()

	processPidAttr, ex := attributes.Get(conventions.AttributeProcessPID)
	if !ex {
		return bundle
	}

	processData := createProcessData(attributes, int(processPidAttr.IntVal()))
	bundle.Metrics.Plugins = append(bundle.Metrics.Plugins, instanaacceptor.NewProcessPluginPayload(processPidAttr.AsString(), processData))

	return bundle
}

func (c *ProcessMetricConverter) Name() string {
	return "ProcessMetricConverter"
}

func createProcessData(attributes pcommon.Map, processPid int) instanaacceptor.ProcessData {
	processData := instanaacceptor.ProcessData{
		PID: int(processPid),
	}

	processExecAttr, ex := attributes.Get(conventions.AttributeProcessExecutablePath)
	if ex {
		processData.Exec = processExecAttr.AsString()
	}

	processCommmandArgsAttr, ex := attributes.Get(conventions.AttributeProcessCommandArgs)
	if ex {
		processData.Args = strings.Split(processCommmandArgsAttr.AsString(), ", ")
	}

	// container info
	processContainerIdAttr, ex := attributes.Get(conventions.AttributeContainerID)
	if ex {
		processData.ContainerID = processContainerIdAttr.AsString()
	}

	return processData
}
