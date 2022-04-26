package model

import (
	"encoding/json"

	instanaacceptor "github.com/instana/go-sensor/acceptor"
)

type Bundle struct {
	Metrics PluginContainer `json:"metrics,omitempty"`
	Spans   []Span          `json:"spans,omitempty"`
}

func NewBundle() Bundle {
	return Bundle{
		Metrics: PluginContainer{
			Plugins: make([]instanaacceptor.PluginPayload, 0),
		},
		Spans: make([]Span, 0),
	}
}

func (b *Bundle) Marshal() ([]byte, error) {
	json, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	return json, nil
}

type PluginContainer struct {
	Plugins []instanaacceptor.PluginPayload `json:"plugins,omitempty"`
}
