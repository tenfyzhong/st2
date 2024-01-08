package st2

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strconv"
	"strings"
)

// GoParser is a Parser to parse golang source.
type GoParser struct {
	ctx Context
}

// NewGoParser create [GoParser]
func NewGoParser(ctx Context) *GoParser {
	return &GoParser{
		ctx: ctx,
	}
}

// Parse method parse golang source
func (p GoParser) Parse(reader io.Reader) ([]*Struct, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", reader, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	ast.FileExports(f)

	res := make([]*Struct, 0)
	for _, node := range f.Decls {
		st := p.processNode(fset, node)
		if st != nil {
			res = append(res, st...)
		}
	}
	return res, nil
}

func (p GoParser) processNode(fset *token.FileSet, node ast.Decl) []*Struct {
	switch n := node.(type) {
	case *ast.GenDecl:
		if n.Tok == token.TYPE {
			st := p.processType(n)
			if st == nil {
				return nil
			}
			return []*Struct{st}
		} else if n.Tok == token.CONST {
			return p.processConst(n)
		}
	}
	return nil
}

func (p GoParser) processType(decl *ast.GenDecl) *Struct {
	if len(decl.Specs) == 0 {
		return nil
	}
	spec, ok := decl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}
	name := spec.Name.Name
	st, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	res := &Struct{
		Type: &StructLikeType{
			Name:   name,
			Source: SLSStruct,
		},
		Comment: p.parseComment(decl.Doc, nil),
	}

	if st.Fields != nil {
		for i, field := range st.Fields.List {
			if len(field.Names) == 0 {
				continue
			}
			fieldName := field.Names[0].Name

			t := p.type2Type(field.Type)
			member := &Member{
				Field:    snake(fieldName),
				Type:     t,
				Index:    i + 1,
				Optional: p.isOptional(field.Type),
				Comment:  p.parseComment(field.Doc, field.Comment),
			}
			if tag := p.tag2GoTag(field.Tag); tag != "" {
				// if there any tag, use the first tag field name as the Member.Field value
				member.Field = tag
			}
			res.Members = append(res.Members, member)
		}
	}
	return res
}

func (p GoParser) tag2GoTag(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}

	// "`json:\"hello_world,omitempty\" proto:\"hello_world\"`"
	item := strings.Trim(tag.Value, "`") // "json:\"hello_world,omitempty\" proto:\"hello_world\""
	item = strings.Split(item, " ")[0]   // "json:\"hello_world,omitempty\""
	items := strings.Split(item, ":")    // ["json:", "\"hello_world,omitempty\""]
	if len(items) < 2 {
		return ""
	}
	item = items[1]                    // "\"hello_world,omitempty\""
	item = strings.TrimSpace(item)     // "\"hello_world,omitempty\""
	item = strings.Trim(item, "\"")    // "hello_world,omitempty"
	item = strings.Split(item, ",")[0] // "hello_world"
	return item
}

func (p GoParser) isOptional(field ast.Expr) bool {
	if f, ok := field.(*ast.StarExpr); ok {
		t := p.type2Type(f.X)
		if t == nil {
			return false
		}
		if p.isBasicType(t.Go()) {
			return true
		}
	}
	return false
}

// There are may be many enum in a const block
func (p GoParser) processConst(decl *ast.GenDecl) []*Struct {
	m := make(map[string]*Struct)
	for _, spec := range decl.Specs {
		v, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		if len(v.Names) == 0 {
			continue
		}
		name := v.Names[0].Name

		typeIdent, ok := v.Type.(*ast.Ident)
		if !ok {
			continue
		}
		if p.isBasicType(typeIdent.Name) {
			continue
		}

		typeName := typeIdent.Name
		st := m[typeName]
		if st == nil {
			st = &Struct{
				Type: &EnumType{
					Name: typeName,
				},
			}
		}

		if len(v.Values) == 0 {
			continue
		}

		value, ok := v.Values[0].(*ast.BasicLit)
		if !ok {
			continue
		}
		valueInt, err := strconv.ParseInt(value.Value, 10, 64)
		if err != nil {
			continue
		}

		st.Members = append(st.Members, &Member{
			Field:   name,
			Type:    st.Type,
			Index:   int(valueInt),
			Comment: p.parseComment(v.Doc, v.Comment),
		})

		m[typeName] = st
	}

	res := make([]*Struct, 0, len(m))
	for _, st := range m {
		res = append(res, st)
	}
	return res
}

func (p GoParser) parseComment(doc *ast.CommentGroup, comment *ast.CommentGroup) Comment {
	c := Comment{}
	if doc != nil {
		for _, comment := range doc.List {
			c.BeginningComments = append(c.BeginningComments, comment.Text)
		}
	}
	if comment != nil && len(comment.List) > 0 {
		c.InlineComment = comment.List[0].Text
	}
	return c
}

func (p GoParser) isBasicType(str string) bool {
	basic := []string{
		StrInt,
		StrInt8,
		StrInt16,
		StrInt32,
		StrInt64,
		StrUint8,
		StrUint16,
		StrUint32,
		StrUint64,
		StrFloat32,
		StrFloat64,
		StrBool,
		StrString,
		StrComplex64,
		StrComplex128,
		StrByte,
		StrRune,
		StrUintptr}
	for _, name := range basic {
		if name == str {
			return true
		}
	}
	return false
}

func (p GoParser) type2Type(t ast.Expr) Type {
	switch t := t.(type) {
	case *ast.Ident:
		return p.nameType(t.Name)
	case *ast.ArrayType:
		return &ArrayType{
			ChildType: p.type2Type(t.Elt),
		}
	case *ast.MapType:
		return &MapType{
			Key:   p.type2Type(t.Key),
			Value: p.type2Type(t.Value),
		}
	case *ast.StarExpr:
		return p.type2Type(t.X)
	case *ast.SelectorExpr:
		selector := p.type2Type(t.X)
		sub := p.type2Type(t.Sel)
		return &StructLikeType{
			Name: strings.TrimLeft(selector.Go(), "*") + "." + strings.TrimLeft(sub.Go(), "*"),
		}
	}
	return nil
}

func (p GoParser) nameType(name string) Type {
	switch name {
	case StrInt:
		return Int64Val
	case StrInt8:
		return Int8Val
	case StrInt16:
		return Int16Val
	case StrInt32:
		return Int32Val
	case StrInt64:
		return Int64Val
	case StrUint:
		return Uint64Val
	case StrUint8:
		return Uint8Val
	case StrUint16:
		return Uint16Val
	case StrUint32:
		return Uint32Val
	case StrUint64:
		return Uint64Val
	case StrFloat32:
		return Float32Val
	case StrFloat64:
		return Float64Val
	case StrBool:
		return BoolVal
	case StrString:
		return StringVal
	case StrComplex64:
		// TODO
	case StrComplex128:
		// TODO
	case StrByte:
		return Uint8Val
	case StrRune:
		// TODO
	case StrUintptr:
		// TODO
	case StrAny:
		return AnyVal
	}
	return &StructLikeType{
		Name: name,
	}
}
