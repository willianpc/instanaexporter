package model

import (
	"strconv"

	instanaacceptor "github.com/instana/go-sensor/acceptor"
)

type PythonProcessData struct {
	PID           int    `json:"pid"`
	Name          string `json:"snapshot.name"`
	PythonVersion string `json:"snapshot.version"`
	PythonArch    string `json:"snapshot.a"`
	PythonFlavor  string `json:"snapshot.f"`
}

func NewPythonRuntimePlugin(data PythonProcessData) instanaacceptor.PluginPayload {
	const pluginName = "com.instana.plugin.python"

	return instanaacceptor.PluginPayload{
		Name:     pluginName,
		EntityID: strconv.Itoa(data.PID),
		Data:     data,
	}
}
