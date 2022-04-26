package model

import (
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"go.opentelemetry.io/collector/model/pdata"
)

type metricsInner struct {
	Gauges         map[string]float64 `json:"gauges,omitempty"`
	HistogramsMean map[string]float64 `json:"histograms_mean,omitempty"`
	Sums           map[string]float64 `json:"sums,omitempty"`
}

func newMetricsInner() metricsInner {
	return metricsInner{
		Gauges:         make(map[string]float64),
		HistogramsMean: make(map[string]float64),
		Sums:           make(map[string]float64),
	}
}

type OpenTelemetryCustomMetricsData struct {
	Metrics metricsInner `json:"metrics,omitempty"`
	Pid     string       `json:"pid,omitempty"`
}

func NewOpenTelemetryCustomMetricsData() OpenTelemetryCustomMetricsData {
	return OpenTelemetryCustomMetricsData{
		Metrics: newMetricsInner(),
	}
}

func (omData *OpenTelemetryCustomMetricsData) AppendMetric(metric pdata.Metric) {
	metricName := metric.Name()

	switch metric.DataType() {
	case pdata.MetricDataTypeGauge:
		for j := 0; j < metric.Gauge().DataPoints().Len(); j++ {
			dp := metric.Gauge().DataPoints().At(j)

			if dp.ValueType() == pdata.MetricValueTypeDouble {
				omData.Metrics.Gauges[metricNameToCompact(metricName, dp.Attributes())] = dp.DoubleVal()
			}

			if dp.ValueType() == pdata.MetricValueTypeInt {
				omData.Metrics.Gauges[metricNameToCompact(metricName, dp.Attributes())] = float64(dp.IntVal())
			}
		}
	case pdata.MetricDataTypeSum:
		for j := 0; j < metric.Sum().DataPoints().Len(); j++ {
			dp := metric.Sum().DataPoints().At(j)

			if dp.ValueType() == pdata.MetricValueTypeDouble {
				omData.Metrics.Sums[metricNameToCompact(metricName, dp.Attributes())] = dp.DoubleVal()
			}

			if dp.ValueType() == pdata.MetricValueTypeInt {
				omData.Metrics.Sums[metricNameToCompact(metricName, dp.Attributes())] = float64(dp.IntVal())
			}
		}
	case pdata.MetricDataTypeHistogram:
		for j := 0; j < metric.Histogram().DataPoints().Len(); j++ {
			dp := metric.Histogram().DataPoints().At(j)

			omData.Metrics.HistogramsMean[metricNameToCompact(metricName, dp.Attributes())] = dp.Sum()
		}
	}
}

func NewOpenTelemetryMetricsPluginPayload(entityID string, data OpenTelemetryCustomMetricsData) instanaacceptor.PluginPayload {
	const pluginName = "com.instana.plugin.otel.metrics"

	return instanaacceptor.PluginPayload{
		Name:     pluginName,
		EntityID: entityID,
		Data:     data,
	}
}
