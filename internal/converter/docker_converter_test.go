package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	conventions "go.opentelemetry.io/collector/semconv/v1.8.0"
)

func TestCanConvert(t *testing.T) {
	converter := DockerContainerMetricConverter{}

	attributes := pcommon.NewMap()
	attributes.InsertString(conventions.AttributeContainerRuntime, "docker")
	attributes.InsertString(conventions.AttributeContainerID, "abc")
	attributes.InsertString(conventions.AttributeContainerImageName, "ubuntu")
	attributes.InsertString(conventions.AttributeContainerImageTag, "latest")
	attributes.InsertString(conventions.AttributeContainerName, "my-container")

	metrics := pmetric.NewMetricSlice()
	metric := metrics.AppendEmpty()
	metric.SetName("container.network.io.usage.tx_packets")
	metric.SetDescription("")
	metric.SetUnit("1")
	metric.SetDataType(pmetric.MetricDataTypeSum)
	metric.Sum().SetIsMonotonic(true)
	metric.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityCumulative)
	md := metric.Sum().DataPoints().AppendEmpty()
	md.SetIntVal(0)
	md.Attributes().InsertString("interface", "eth0")

	plugins := converter.ConvertMetrics(attributes, metrics)

	assert.Equal(t, 1, len(plugins))
	assert.Equal(t, "com.instana.plugin.docker", plugins[0].Name)
}
