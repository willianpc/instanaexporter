package converter

import (
	"strings"

	"github.com/ibm-observability/ibminstanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/model/pdata"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.8.0"
)

var _ Converter = (*ProcessMetricConverter)(nil)

type ProcessMetricConverter struct{}

func (c *ProcessMetricConverter) AcceptsMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) bool {
	_, ex := attributes.Get(conventions.AttributeProcessPID)

	return ex
}

func (c *ProcessMetricConverter) ConvertMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) []instanaacceptor.PluginPayload {
	plugins := make([]instanaacceptor.PluginPayload, 0)

	processPidAttr, ex := attributes.Get(conventions.AttributeProcessPID)
	if !ex {
		return plugins
	}

	processData := createProcessData(attributes, int(processPidAttr.IntVal()))

	plugins = append(plugins, instanaacceptor.NewProcessPluginPayload(processPidAttr.AsString(), processData))

	return plugins
}

func (c *ProcessMetricConverter) AcceptsSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) bool {
	_, ex := attributes.Get(conventions.AttributeProcessPID)

	return ex
}

func (c *ProcessMetricConverter) ConvertSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) model.Bundle {

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

func createProcessData(attributes pdata.AttributeMap, processPid int) instanaacceptor.ProcessData {
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
