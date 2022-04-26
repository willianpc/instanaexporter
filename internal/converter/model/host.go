package model

import (
	instanaacceptor "github.com/instana/go-sensor/acceptor"
)

type NetworkInterface struct {
}

type FileSystem struct {
	Options   string `json:"options,omitempty"`
	Systype   string `json:"systype,omitempty"`
	Mount     string `json:"mount,omitempty"`
	ICapacity int64  `json:"icapacity,omitempty"`
	Capacity  int64  `json:"capacity,omitempty"`
}

type CpuSummary struct {
	Sys   float64 `json:"sys,omitempty"`
	Idle  float64 `json:"idle,omitempty"`
	User  float64 `json:"user,omitempty"`
	Steal float64 `json:"steal,omitempty"`
}

type HostData struct {
	Tags
	BootId string `json:"bootId,omitempty"`

	Cpu          CpuSummary   `json:"cpu,omitempty"`
	CpuSummaries []CpuSummary `json:"cpus,omitempty"`

	HostName string `json:"hostname,omitempty"`

	CpuCount int    `json:"cpu.count,omitempty"`
	CpuModel string `json:"cpu.model,omitempty"`

	FileSystems map[string]FileSystem       `json:"filesystems,omitempty"`
	Interfaces  map[string]NetworkInterface `json:"interfaces,omitempty"`

	MemoryTotal int64 `json:"memory.total,omitempty"`

	OpenFilesMax int64 `json:"openFiles.max,omitempty"`

	OsArch    string `json:"os.arch,omitempty"`
	OsName    string `json:"os.name,omitempty"`
	OsVersion string `json:"os.version,omitempty"`

	Start int64 `json:"start,omitempty"`

	SystemSerialNumber string `json:"systemSerialnumber,omitempty"`
}

func NewHostData() HostData {
	return HostData{}
}

// NewDockerPluginPayload returns payload for the Docker plugin of Instana acceptor
func NewHostPluginPayload(entityID string, data HostData) instanaacceptor.PluginPayload {
	const pluginName = "com.instana.plugin.host"

	return instanaacceptor.PluginPayload{
		Name:     pluginName,
		EntityID: entityID,
		Data:     data,
	}
}
