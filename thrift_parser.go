package st2

import (
	"errors"
	"io"

	"github.com/cloudwego/thriftgo/parser"
)

// ThriftParser is a Parser to parse thrift source
type ThriftParser struct {
	ctx Context
}

// NewThriftParser create [ThriftParser]
func NewThriftParser(ctx Context) *ThriftParser {
	return &ThriftParser{
		ctx: ctx,
	}
}

// Parse method parse thrift source
func (p ThriftParser) Parse(reader io.Reader) ([]*Struct, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("read data failed")
	}

	if len(data) == 0 {
		return nil, nil
	}

	thrift, err := parser.ParseString("", string(data))
	if err != nil {
		return nil, err
	}

	if thrift == nil {
		return nil, errors.New("Parse thrift failed")
	}

	res := make([]*Struct, 0)

	for _, e := range thrift.Enums {
		res = append(res, p.enum2struct(e))
	}

	for _, s := range thrift.Structs {
		res = append(res, p.structLike2struct(s, SLSStruct))
	}

	for _, u := range thrift.Unions {
		res = append(res, p.structLike2struct(u, SLSUnion))
	}

	return res, nil
}

func (p ThriftParser) enum2struct(e *parser.Enum) *Struct {
	s := &Struct{
		Type: &EnumType{
			Name: e.Name,
		},
	}

	for _, value := range e.Values {
		member := &Member{
			Field: value.Name,
			Type:  s.Type,
			Index: int(value.Value),
		}
		s.Members = append(s.Members, member)
	}
	return s
}

func (p ThriftParser) structLike2struct(sl *parser.StructLike, source StructLikeSource) *Struct {
	s := &Struct{
		Type: &StructLikeType{
			Name:   sl.Name,
			Source: source,
		},
	}

	for _, field := range sl.Fields {
		t := p.type2Type(field.Type)
		if t == nil {
			continue
		}

		member := &Member{
			Field:    field.Name,
			Type:     t,
			Index:    int(field.ID),
			Optional: field.Requiredness == parser.FieldType_Optional,
		}
		s.Members = append(s.Members, member)
	}

	return s
}

func (p ThriftParser) type2Type(t *parser.Type) Type {
	if t == nil {
		return nil
	}

	switch t.Name {
	case StrBool:
		return BoolVal
	case StrByte:
		return Int8Val
	case StrI16:
		return Int16Val
	case StrI32:
		return Int32Val
	case StrI64:
		return Int64Val
	case StrDouble:
		return Float64Val
	case StrString:
		return StringVal
	case StrBinary:
		return BinaryVal
	case StrMap:
		return &MapType{
			Key:   p.type2Type(t.KeyType),
			Value: p.type2Type(t.ValueType),
		}
	case StrList:
		return &ArrayType{
			ChildType: p.type2Type(t.ValueType),
		}
	case StrSet:
		return &SetType{
			Key: p.type2Type(t.ValueType),
		}
	default:
		return &StructLikeType{
			Name: t.Name,
		}
	}
	return nil
}
