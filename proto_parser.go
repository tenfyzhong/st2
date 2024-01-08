package st2

import (
	"io"
	"strconv"

	"github.com/yoheimuta/go-protoparser/v4"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

// ProtoParser is a Parser to parse protobuf source
type ProtoParser struct {
	ctx Context
}

// NewProtoParser create [ProtoParser]
func NewProtoParser(ctx Context) *ProtoParser {
	return &ProtoParser{
		ctx: ctx,
	}
}

// Parse method parse protobuf source
func (p ProtoParser) Parse(reader io.Reader) ([]*Struct, error) {
	got, err := protoparser.Parse(reader)
	if err != nil {
		return nil, err
	}

	res := make([]*Struct, 0, len(got.ProtoBody))
	for _, pb := range got.ProtoBody {
		v := newProtoVisitor(p.ctx)
		pb.Accept(v)
		if v.Struct != nil {
			res = append(res, v.Struct)
		}
	}
	return res, nil
}

type protoVisitor struct {
	Struct *Struct
	ctx    Context
}

func newProtoVisitor(ctx Context) *protoVisitor {
	return &protoVisitor{
		ctx: ctx,
	}
}

func (v *protoVisitor) VisitComment(c *parser.Comment) {
}

func (v *protoVisitor) VisitEmptyStatement(s *parser.EmptyStatement) bool {
	return true
}

func (v *protoVisitor) VisitEnum(e *parser.Enum) bool {
	v.Struct = &Struct{
		Type: &EnumType{
			Name: e.EnumName,
		},
		Comment: v.comment2Comment(e.Comments, e.InlineCommentBehindLeftCurly),
	}
	return true
}

func (v *protoVisitor) VisitEnumField(f *parser.EnumField) bool {
	number, _ := strconv.ParseInt(f.Number, 10, 32)
	v.Struct.Members = append(v.Struct.Members, &Member{
		Field:   f.Ident,
		Type:    v.Struct.Type,
		Index:   int(number),
		Comment: v.comment2Comment(f.Comments, f.InlineComment),
	})
	return true
}

func (v *protoVisitor) VisitExtend(e *parser.Extend) bool {
	return true
}

func (v *protoVisitor) VisitExtensions(e *parser.Extensions) bool {
	return true
}

func (v *protoVisitor) VisitField(f *parser.Field) bool {
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

func (v *protoVisitor) VisitGroupField(g *parser.GroupField) bool {
	return true
}

func (v *protoVisitor) VisitImport(i *parser.Import) bool {
	// fmt.Printf("%+v\n", i)
	return true
}

func (v *protoVisitor) VisitMapField(f *parser.MapField) bool {
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

func (v *protoVisitor) VisitMessage(m *parser.Message) bool {
	v.Struct = &Struct{
		Type: &StructLikeType{
			Name:   m.MessageName,
			Source: SLSStruct,
		},
		Comment: v.comment2Comment(m.Comments, m.InlineCommentBehindLeftCurly),
	}
	return true
}

func (v *protoVisitor) VisitOneof(o *parser.Oneof) bool {
	return true
}

func (v *protoVisitor) VisitOneofField(f *parser.OneofField) bool {
	// fmt.Printf("%+v\n", f)
	return true
}

func (v *protoVisitor) VisitOption(o *parser.Option) bool {
	// fmt.Printf("%+v\n", o)
	return true
}

func (v *protoVisitor) VisitPackage(p *parser.Package) bool {
	// fmt.Printf("%+v\n", p)
	return true
}

func (v *protoVisitor) VisitReserved(r *parser.Reserved) bool {
	// fmt.Printf("%+v\n", r)
	return true
}

func (v *protoVisitor) VisitRPC(rpc *parser.RPC) bool {
	// fmt.Printf("%+v\n", rpc)
	return true
}

func (v *protoVisitor) VisitService(s *parser.Service) bool {
	// fmt.Printf("%+v\n", s)
	return true
}

func (v *protoVisitor) VisitSyntax(s *parser.Syntax) bool {
	// fmt.Printf("%+v\n", s)
	return true
}

func (v *protoVisitor) type2Type(str string) Type {
	switch str {
	case StrDouble:
		return Float64Val
	case StrFloat:
		return Float32Val
	case StrInt32:
		return Int32Val
	case StrInt64:
		return Int64Val
	case StrUint32:
		return Uint32Val
	case StrUint64:
		return Uint64Val
	case StrSint32:
		return Int32Val
	case StrSint64:
		return Int64Val
	case StrFixed32:
		return Int32Val
	case StrFixed64:
		return Int64Val
	case StrBool:
		return BoolVal
	case StrString:
		return StringVal
	case StrBytes:
		return BinaryVal
	case StrPbAny:
		return AnyVal
	}
	return &StructLikeType{
		Name: str,
	}
}

func (v *protoVisitor) comment2Comment(beginComments []*parser.Comment, inlineComment *parser.Comment) Comment {
	comment := Comment{}
	for _, c := range beginComments {
		comment.BeginningComments = append(comment.BeginningComments, c.Raw)
	}
	if inlineComment != nil {
		comment.InlineComment = inlineComment.Raw
	}
	return comment
}
