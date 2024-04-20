package st2

import (
	"bytes"

	xj "github.com/basgys/goxml2json"
)

type XMLUnmarshalTagFormat struct {
	ContentTagPrefix   string
	AttributeTagPrefix string
}

func (x XMLUnmarshalTagFormat) Unmarshal(data []byte, v any) error {
	reader := bytes.NewReader(data)
	root := &xj.Node{}
	decoder := xj.NewDecoder(reader)
	err := decoder.DecodeWithCustomPrefixes(root, x.ContentTagPrefix, x.AttributeTagPrefix)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	encoder := xj.NewEncoder(buf)
	err = encoder.Encode(root)
	if err != nil {
		return err
	}

	json := buf.Bytes()

	return jsonapi.Unmarshal(json, v)
}

func (x XMLUnmarshalTagFormat) TagFormat() string {
	return `xml:"%s"`
}

func NewXMLParser(ctx Context) *StructuredParser {
	return NewStructuredParser(ctx, &XMLUnmarshalTagFormat{
		ContentTagPrefix:   ctx.XMLContext.ContentTagPrefix,
		AttributeTagPrefix: ctx.XMLContext.AttributeTagPrefix,
	})
}
