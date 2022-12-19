package st2

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/cloudwego/thriftgo/parser"
)

type ThriftParser struct {
}

func NewThriftParser() *ThriftParser {
	return &ThriftParser{}
}

func (p ThriftParser) Parse(reader io.Reader) ([]*Struct, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("read data failed")
	}

	thrift, err := parser.ParseString("", string(data))
	if err != nil {
		return nil, errors.New("Parse thrift failed")
	}

	if thrift == nil {
		return nil, errors.New("Parse thrift failed")
	}

	res := make([]*Struct, 0)

	for _, e := range thrift.Enums {
		res = append(res, p.enum2struct(e))
	}

	for _, s := range thrift.Structs {
		res = append(res, p.structLike2struct(s, "struct"))
	}

	for _, u := range thrift.Unions {
		res = append(res, p.structLike2struct(u, "union"))
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

func (p ThriftParser) structLike2struct(sl *parser.StructLike, stType string) *Struct {
	s := &Struct{
		Type: &StructType{
			Name: sl.Name,
			Type: stType,
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
	case "bool":
		return BoolVal
	case "byte":
		return Int8Val
	case "i16":
		return Int16Val
	case "i32":
		return Int32Val
	case "i64":
		return Int64Val
	case "double":
		return Float64Val
	case "string":
		return StringVal
	case "binary":
		return BinaryVal
	case "map":
		return &MapType{
			Key:   p.type2Type(t.KeyType),
			Value: p.type2Type(t.ValueType),
		}
	case "list":
		return &ArrayType{
			ChildType: p.type2Type(t.ValueType),
		}
	case "set":
		return &SetType{
			Key: p.type2Type(t.ValueType),
		}
	default:
		return &RawType{
			Name: t.Name,
		}
	}
	return nil
}
