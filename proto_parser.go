package st2

import (
	"io"
	"strconv"

	"github.com/yoheimuta/go-protoparser/v4"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

type ProtoParser struct {
	ctx Context
}

func NewProtoParser(ctx Context) *ProtoParser {
	return &ProtoParser{
		ctx: ctx,
	}
}

func (p ProtoParser) Parse(reader io.Reader) ([]*Struct, error) {
	got, err := protoparser.Parse(reader)
	if err != nil {
		return nil, err
	}

	res := make([]*Struct, 0, len(got.ProtoBody))
	for _, pb := range got.ProtoBody {
		v := NewProtoVisitor(p.ctx)
		pb.Accept(v)
		if v.Struct != nil {
			res = append(res, v.Struct)
		}
	}
	return res, nil
}

type ProtoVisitor struct {
	Struct *Struct
	ctx    Context
}

func NewProtoVisitor(ctx Context) *ProtoVisitor {
	return &ProtoVisitor{
		ctx: ctx,
	}
}

func (v *ProtoVisitor) VisitComment(c *parser.Comment) {
}

func (v *ProtoVisitor) VisitEmptyStatement(s *parser.EmptyStatement) bool {
	return true
}

func (v *ProtoVisitor) VisitEnum(e *parser.Enum) bool {
	v.Struct = &Struct{
		Type: &EnumType{
			Name: e.EnumName,
		},
		Comment: v.comment2Comment(e.Comments, e.InlineCommentBehindLeftCurly),
	}
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
	return true
}

func (v *ProtoVisitor) VisitExtensions(e *parser.Extensions) bool {
	return true
}

func (v *ProtoVisitor) VisitField(f *parser.Field) bool {
	fieldNumber, _ := strconv.ParseInt(f.FieldNumber, 10, 64)
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
	return true
}

func (v *ProtoVisitor) VisitImport(i *parser.Import) bool {
	// fmt.Printf("%+v\n", i)
	return true
}

func (v *ProtoVisitor) VisitMapField(f *parser.MapField) bool {
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
	v.Struct = &Struct{
		Type: &StructLikeType{
			Name:   m.MessageName,
			Source: SLSStruct,
		},
		Comment: v.comment2Comment(m.Comments, m.InlineCommentBehindLeftCurly),
	}
	return true
}

func (v *ProtoVisitor) VisitOneof(o *parser.Oneof) bool {
	return true
}

func (v *ProtoVisitor) VisitOneofField(f *parser.OneofField) bool {
	// fmt.Printf("%+v\n", f)
	return true
}

func (v *ProtoVisitor) VisitOption(o *parser.Option) bool {
	// fmt.Printf("%+v\n", o)
	return true
}

func (v *ProtoVisitor) VisitPackage(p *parser.Package) bool {
	// fmt.Printf("%+v\n", p)
	return true
}

func (v *ProtoVisitor) VisitReserved(r *parser.Reserved) bool {
	// fmt.Printf("%+v\n", r)
	return true
}

func (v *ProtoVisitor) VisitRPC(rpc *parser.RPC) bool {
	// fmt.Printf("%+v\n", rpc)
	return true
}

func (v *ProtoVisitor) VisitService(s *parser.Service) bool {
	// fmt.Printf("%+v\n", s)
	return true
}

func (v *ProtoVisitor) VisitSyntax(s *parser.Syntax) bool {
	// fmt.Printf("%+v\n", s)
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
		return Int64Val
	case "fixed32":
		return Int32Val
	case "fixed64":
		return Int64Val
	case "bool":
		return BoolVal
	case "string":
		return StringVal
	case "bytes":
		return BinaryVal
	}
	return &StructLikeType{
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
