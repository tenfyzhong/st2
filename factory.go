package st2

import "github.com/tenfyzhong/st2/tmpl"

// CreateParser Create a [Parse] belongs to the context
func CreateParser(ctx Context) Parse {
	switch ctx.Src {
	case LangGo:
		return NewGoParser(ctx)
	case LangJson:
		return NewJsonParser(ctx)
	case LangYaml:
		return NewYamlParser(ctx)
	case LangProto:
		return NewProtoParser(ctx)
	case LangThrift:
		return NewThriftParser(ctx)
	case LangCsv:
		return NewCsvParser(ctx)
	case LangXML:
		return NewXMLParser(ctx)
	case LangToml:
		return NewTomlParser(ctx)
	}
	return nil
}

// CreateTmpl Get the template data belongs to the context
func CreateTmpl(ctx Context) string {
	switch ctx.Dst {
	case LangGo:
		return tmpl.Go
	case LangJson:
	case LangYaml:
	case LangProto:
		return tmpl.Proto
	case LangThrift:
		return tmpl.Thrift
	}
	return ""
}
