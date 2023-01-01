package st2

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type GoParser struct {
	ctx Context
}

func NewGoParser(ctx Context) *GoParser {
	return &GoParser{
		ctx: ctx,
	}
}

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
		Type: &StructType{
			Name: name,
			Type: "struct",
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
				Field:    fieldName,
				Type:     t,
				Index:    i + 1,
				Optional: p.isOptional(field.Type),
				Comment:  p.parseComment(field.Doc, field.Comment),
			}
			if field.Tag != nil {
				member.GoTag = p.tag2GoTag(field.Tag)
			}
			res.Members = append(res.Members, member)
		}
	}
	return res
}

func (p GoParser) tag2GoTag(tag *ast.BasicLit) []string {
	// "`json:\"C,omitempty\" proto:\"C\"`"
	re := regexp.MustCompile(`(\w*):`)
	items := re.FindAllStringSubmatch(tag.Value, -1)
	res := make([]string, 0)
	for _, item := range items {
		if len(item) > 1 {
			res = append(res, item[1])
		}
	}
	return res
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
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"float32",
		"float64",
		"bool",
		"string",
		"complex64",
		"complex128",
		"byte",
		"rune",
		"uintptr"}
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
		return &RawType{
			Name: strings.TrimLeft(selector.Go(), "*") + "." + strings.TrimLeft(sub.Go(), "*"),
		}
	}
	return nil
}

func (p GoParser) nameType(name string) Type {
	switch name {
	case "int":
		return Int64Val
	case "int8":
		return Int8Val
	case "int16":
		return Int16Val
	case "int32":
		return Int32Val
	case "int64":
		return Int64Val
	case "uint8":
		return Uint8Val
	case "uint16":
		return Uint16Val
	case "uint32":
		return Uint32Val
	case "uint64":
		return Uint64Val
	case "float32":
		return Float32Val
	case "float64":
		return Float64Val
	case "bool":
		return BoolVal
	case "string":
		return StringVal
	case "complex64":
		// TODO
	case "complex128":
		// TODO
	case "byte":
		return Uint8Val
	case "rune":
		// TODO
	case "uintptr":
		// TODO
	}
	return &RawType{
		Name: name,
	}
}
