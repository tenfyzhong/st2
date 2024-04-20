package st2

import (
	"github.com/pelletier/go-toml/v2"
)

type TomlUnmarshalTagFormat struct{}

func (t TomlUnmarshalTagFormat) Unmarshal(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}

func (t TomlUnmarshalTagFormat) TagFormat() string {
	return `toml:"%s"`
}

func NewTomlParser(ctx Context) *StructuredParser {
	return NewStructuredParser(ctx, &TomlUnmarshalTagFormat{})
}
