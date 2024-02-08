package st2

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	jsonapi = jsoniter.Config{UseNumber: true}.Froze()
)

type JsonUnmarshalTagFormat struct{}

func (j JsonUnmarshalTagFormat) Unmarshal(data []byte, v interface{}) error {
	return jsonapi.Unmarshal(data, v)
}

func (j JsonUnmarshalTagFormat) TagFormat() string {
	return `json:"%s,omitempty"`
}

// NewJsonParser create [StructuredParser] to parse json source
func NewJsonParser(ctx Context) *StructuredParser {
	return NewStructuredParser(ctx, &JsonUnmarshalTagFormat{})
}
