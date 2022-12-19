package st2

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/yoheimuta/go-protoparser/v4"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

type ProtoParser struct {
}

func NewProtoParser() *ProtoParser {
	return &ProtoParser{}
}

func (p ProtoParser) Parse(reader io.Reader) ([]*Struct, error) {
	got, err := protoparser.Parse(reader)
	if err != nil {
		return nil, errors.New("Parse proto failed")
	}

	res := make([]*Struct, 0, len(got.ProtoBody))
	for _, pb := range got.ProtoBody {
		v := NewProtoVisitor()
		pb.Accept(v)
		res = append(res, v.Struct)
	}
	return res, nil
}

type ProtoVisitor struct {
	Struct *Struct
}

func NewProtoVisitor() *ProtoVisitor {
	return &ProtoVisitor{
		Struct: &Struct{},
	}
}

func (v *ProtoVisitor) VisitComment(c *parser.Comment) {
	// fmt.Printf("%+v\n", c)
}

func (v *ProtoVisitor) VisitEmptyStatement(s *parser.EmptyStatement) bool {
	// fmt.Printf("%+v\n", s)
	return true
}

func (v *ProtoVisitor) VisitEnum(e *parser.Enum) bool {
	v.Struct.Type = &EnumType{
		Name: e.EnumName,
	}
	v.Struct.Comment = v.comment2Comment(e.Comments, e.InlineCommentBehindLeftCurly)
	return true
}

func (v *ProtoVisitor) VisitEnumField(f *parser.EnumField) bool {
	number, _ := strconv.ParseInt(f.Number, 10, 32)
	v.Struct.Members = append(v.Struct.Members, &Member{
		Field:   f.Ident,
		Type:    v.Struct.Type,
		Index:   int(number),
		Comment: v.comment2Comment(f.Comments, f.InlineComment),
	})
	return true
}

func (v *ProtoVisitor) VisitExtend(e *parser.Extend) bool {
	// fmt.Printf("%+v\n", e)
	return true
}

func (v *ProtoVisitor) VisitExtensions(e *parser.Extensions) bool {
	// fmt.Printf("%+v\n", e)
	return true
}

func (v *ProtoVisitor) VisitField(f *parser.Field) bool {
	// fmt.Printf("%+v\n", f)
	fieldNumber, _ := strconv.ParseInt(f.FieldNumber, 10, 64)
	if f.IsRepeated {

	}
	v.Struct.Members = append(v.Struct.Members, &Member{
		Field: f.FieldName,
		Type: func() Type {
			if f.IsRepeated {
				return &ArrayType{
					ChildType: v.type2Type(f.Type),
				}
			}
			return v.type2Type(f.Type)
		}(),
		Index:    int(fieldNumber),
		Optional: f.IsOptional,
		Comment:  v.comment2Comment(f.Comments, f.InlineComment),
	})
	return true
}

func (v *ProtoVisitor) VisitGroupField(g *parser.GroupField) bool {
	// fmt.Printf("%+v\n", g)
	return true
}

func (v *ProtoVisitor) VisitImport(i *parser.Import) bool {
	fmt.Printf("%+v\n", i)
	return true
}

func (v *ProtoVisitor) VisitMapField(f *parser.MapField) bool {
	// fmt.Printf("%+v\n", f)
	fieldNumber, _ := strconv.ParseInt(f.FieldNumber, 10, 64)
	v.Struct.Members = append(v.Struct.Members, &Member{
		Field: f.MapName,
		Type: &MapType{
			Key:   v.type2Type(f.KeyType),
			Value: v.type2Type(f.Type),
		},
		Index:   int(fieldNumber),
		Comment: v.comment2Comment(f.Comments, f.InlineComment),
	})
	return true
}

func (v *ProtoVisitor) VisitMessage(m *parser.Message) bool {
	// fmt.Printf("%+v\n", m)
	v.Struct.Type = &StructType{
		Name: m.MessageName,
		Type: "struct",
	}
	v.Struct.Comment = v.comment2Comment(m.Comments, m.InlineCommentBehindLeftCurly)
	return true
}

func (v *ProtoVisitor) VisitOneof(o *parser.Oneof) bool {
	fmt.Printf("%+v\n", o)
	return true
}

func (v *ProtoVisitor) VisitOneofField(f *parser.OneofField) bool {
	fmt.Printf("%+v\n", f)
	return true
}

func (v *ProtoVisitor) VisitOption(o *parser.Option) bool {
	fmt.Printf("%+v\n", o)
	return true
}

func (v *ProtoVisitor) VisitPackage(p *parser.Package) bool {
	fmt.Printf("%+v\n", p)
	return true
}

func (v *ProtoVisitor) VisitReserved(r *parser.Reserved) bool {
	fmt.Printf("%+v\n", r)
	return true
}

func (v *ProtoVisitor) VisitRPC(rpc *parser.RPC) bool {
	fmt.Printf("%+v\n", rpc)
	return true
}

func (v *ProtoVisitor) VisitService(s *parser.Service) bool {
	fmt.Printf("%+v\n", s)
	return true
}

func (v *ProtoVisitor) VisitSyntax(s *parser.Syntax) bool {
	fmt.Printf("%+v\n", s)
	return true
}

func (v *ProtoVisitor) type2Type(str string) Type {
	switch str {
	case "double":
		return Float64Val
	case "float":
		return Float32Val
	case "int32":
		return Int32Val
	case "int64":
		return Int64Val
	case "uint32":
		return Uint32Val
	case "uint64":
		return Uint64Val
	case "sint32":
		return Int32Val
	case "sint64":
		return Int32Val
	case "fixed32":
		return Int32Val
	case "fixed64":
		return Int32Val
	case "bool":
		return BoolVal
	case "string":
		return StringVal
	case "bytes":
		return BinaryVal
	}
	return &RawType{
		Name: str,
	}
}

func (v *ProtoVisitor) comment2Comment(beginComments []*parser.Comment, inlineComment *parser.Comment) Comment {
	comment := Comment{}
	for _, c := range beginComments {
		comment.BeginningComments = append(comment.BeginningComments, c.Raw)
	}
	if inlineComment != nil {
		comment.InlineComment = inlineComment.Raw
	}
	return comment
}
