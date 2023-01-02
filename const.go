package st2

import "github.com/tenfyzhong/st2/tmpl"

const (
	LangGo     = "go"
	LangJson   = "json"
	LangProto  = "proto"
	LangThrift = "thrift"

	RootDefault = "Root"
)

var (
	SourceLangs      = []string{LangJson, LangProto, LangThrift, LangGo}
	DestinationLangs = []string{LangGo, LangProto, LangThrift}
	LangTmplMap      = map[string]string{
		LangGo:     tmpl.Go,
		LangProto:  tmpl.Proto,
		LangThrift: tmpl.Thrift,
	}
)
