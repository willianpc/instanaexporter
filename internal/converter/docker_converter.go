package converter

import (
	"time"

	"github.com/ibm-observability/instanaexporter/internal/converter/model"
	instanaacceptor "github.com/instana/go-sensor/acceptor"
	"github.com/instana/go-sensor/docker"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.8.0"
)

var DockerMetricMap = map[string]string{
	"container.blockio.io_service_bytes_recursive.read":  "",
	"container.blockio.io_service_bytes_recursive.write": "",
}

var _ Converter = (*DockerContainerMetricConverter)(nil)

type DockerContainerMetricConverter struct{}

func (c *DockerContainerMetricConverter) AcceptsMetrics(attributes pcommon.Map, metricSlice pmetric.MetricSlice) bool {
	if !containsMetricWithPrefix(metricSlice, "container.") {
		return false
	}

	//_, ex := attributes.Get(conventions.AttributeContainerRuntime)
	//if !ex { // || containerRuntime.AsString() != "docker"
	//	return false
	//}

	_, ex := attributes.Get(conventions.AttributeContainerID)
	if !ex {
		return false
	}

	_, ex = attributes.Get(conventions.AttributeContainerImageName)
	if !ex {
		return false
	}

	_, ex = attributes.Get(conventions.AttributeContainerName)
	if !ex {
		return false
	}

	return true
}

func (c *DockerContainerMetricConverter) ConvertMetrics(attributes pcommon.Map, metricSlice pmetric.MetricSlice) []instanaacceptor.PluginPayload {
	containerId, ex := attributes.Get(conventions.AttributeContainerID)
	if !ex {
		return make([]instanaacceptor.PluginPayload, 0)
	}

	containerImage, ex := attributes.Get(conventions.AttributeContainerImageName)
	if !ex {
		return make([]instanaacceptor.PluginPayload, 0)
	}

	containerName, ex := attributes.Get(conventions.AttributeContainerName)
	if !ex {
		return make([]instanaacceptor.PluginPayload, 0)
	}

	dockerData := instanaacceptor.DockerData{}
	dockerData.ID = containerId.AsString()

	// TODO: Calculate them deltas
	dockerData.BlockIO = instanaacceptor.NewDockerBlockIOStatsDelta(
		docker.ContainerBlockIOStats{},
		docker.ContainerBlockIOStats{},
	)

	// TODO: Calculate them deltas
	dockerData.CPU = instanaacceptor.NewDockerCPUStatsDelta(
		docker.ContainerCPUStats{},
		docker.ContainerCPUStats{},
	)

	// TODO: Add to attributes in dockerstatssreceiver
	dockerData.Command = "/bin/bash"

	// TODO: Add to attributes in dockerstatssreceiver
	created, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	dockerData.CreatedAt = created

	// TODO: Add to attributes in dockerstatssreceiver
	dockerData.DockerAPIVersion = "0.0.0-invalid"

	// TODO: Add to attributes in dockerstatssreceiver
	dockerData.DockerVersion = "0.0.0-invalid"

	dockerData.Image = containerImage.AsString()

	// TODO: Calculate them deltas
	dockerData.Memory = instanaacceptor.NewDockerMemoryStatsUpdate(
		docker.ContainerMemoryStats{},
		docker.ContainerMemoryStats{},
	)

	dockerData.Names = []string{containerName.AsString()}

	// TODO: Calculate them deltas
	dockerData.Network = instanaacceptor.NewDockerNetworkAggregatedStatsDelta(
		map[string]docker.ContainerNetworkStats{},
		map[string]docker.ContainerNetworkStats{},
	)

	// TODO: Add to attributes in dockerstatssreceiver
	dockerData.NetworkMode = "host"

	// TODO: Add to attributes in dockerstatssreceiver
	dockerData.PortBindings = "80:80"

	// TODO: Add to attributes in dockerstatssreceiver
	started, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	dockerData.StartedAt = started

	// TODO: Add to attributes in dockerstatssreceiver
	dockerData.StorageDriver = "my-driver"

	return []instanaacceptor.PluginPayload{
		instanaacceptor.NewDockerPluginPayload(
			containerId.AsString(),
			dockerData,
		),
	}
}

func (c *DockerContainerMetricConverter) AcceptsSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) bool {

	return false
}

func (c *DockerContainerMetricConverter) ConvertSpans(attributes pcommon.Map, spanSlice ptrace.SpanSlice) model.Bundle {

	return model.NewBundle()
}

func (c *DockerContainerMetricConverter) Name() string {
	return "ConvertAllConverter"
}
