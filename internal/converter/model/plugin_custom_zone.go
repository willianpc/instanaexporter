package model

import (
	instanaacceptor "github.com/instana/go-sensor/acceptor"
)

type GenericZoneData struct {
	GroupId string `json:"availability-zone"`
}

func NewGenericZone(zoneName string) instanaacceptor.PluginPayload {
	const pluginName = "com.instana.plugin.generic.hardware"

	genZone := GenericZoneData{GroupId: zoneName}

	return instanaacceptor.PluginPayload{
		Name:     pluginName,
		EntityID: "localhost",
		Data:     genZone,
	}
}
