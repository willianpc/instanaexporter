package model

type Tags struct {
	Tags []string `json:"tags,omitempty"`
}

func (t *Tags) AddTag(tag string) {
	t.Tags = append(t.Tags, tag)
}
