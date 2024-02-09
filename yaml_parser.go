package st2

import (
	"gopkg.in/yaml.v3"
)

type YamlUnmarshalTagFormat struct{}

func (j YamlUnmarshalTagFormat) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

func (j YamlUnmarshalTagFormat) TagFormat() string {
	return `yaml:"%s"`
}

// NewYamlParser create [StructuredParser] to parse yaml source
func NewYamlParser(ctx Context) *StructuredParser {
	return NewStructuredParser(ctx, &YamlUnmarshalTagFormat{})
}
