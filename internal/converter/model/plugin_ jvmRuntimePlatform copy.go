package model

import (
	"strconv"

	instanaacceptor "github.com/instana/go-sensor/acceptor"
)

type JVMProcessData struct {
	PID        int    `json:"pid"`
	Name       string `json:"name"`
	JvmName    string `json:"jvm.name"`
	JvmVendor  string `json:"jvm.vendor"`
	JvmVersion string `json:"jvm.version"`
}

func NewJvmRuntimePlugin(data JVMProcessData) instanaacceptor.PluginPayload {
	const pluginName = "com.instana.plugin.java"

	return instanaacceptor.PluginPayload{
		Name:     pluginName,
		EntityID: strconv.Itoa(data.PID),
		Data:     data,
	}
}
