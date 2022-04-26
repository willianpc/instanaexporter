package converter

import (
	"os"
	"strconv"

	"github.com/ibm-observability/ibminstanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/model/pdata"
)

var _ Converter = (*CollectorMetricsConverter)(nil)

type CollectorMetricsConverter struct{}

func (c *CollectorMetricsConverter) AcceptsMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) bool {

	return containsMetricWithPrefix(metricSlice, "otelcol_")
}

func (c *CollectorMetricsConverter) ConvertMetrics(attributes pdata.AttributeMap, metricSlice pdata.MetricSlice) []instanaacceptor.PluginPayload {
	pid := os.Getpid()

	collectorProcessPlugin := instanaacceptor.NewProcessPluginPayload(strconv.Itoa(pid), instanaacceptor.ProcessData{
		PID:     pid,
		Exec:    "otelcol-idot",
		HostPID: pid,
	})

	goPlugin := instanaacceptor.NewGoProcessPluginPayload(instanaacceptor.GoProcessData{
		PID: pid,
		Snapshot: &instanaacceptor.RuntimeInfo{
			Name: "OpenTelemetry Collector",
		},
	})

	customMetricsData := model.NewOpenTelemetryCustomMetricsData()
	customMetricsData.Pid = strconv.Itoa(pid)

	for i := 0; i < metricSlice.Len(); i++ {
		metric := metricSlice.At(i)

		customMetricsData.AppendMetric(metric)
	}

	otelMetricsCollectorPlugin := model.NewOpenTelemetryMetricsPluginPayload(strconv.Itoa(pid), customMetricsData)

	return []instanaacceptor.PluginPayload{
		collectorProcessPlugin,
		goPlugin,
		otelMetricsCollectorPlugin,
	}
}

func (c *CollectorMetricsConverter) AcceptsSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) bool {

	return false
}

func (c *CollectorMetricsConverter) ConvertSpans(attributes pdata.AttributeMap, spanSlice pdata.SpanSlice) model.Bundle {

	return model.NewBundle()
}

func (c *CollectorMetricsConverter) Name() string {
	return "CollectorMetricsConverter"
}
