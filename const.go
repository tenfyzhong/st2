package st2

import "github.com/tenfyzhong/st2/tmpl"

const (
	LangGo     = "go"
	LangJson   = "json"
	LangProto  = "proto"
	LangThrift = "thrift"
	LangCsv    = "csv"
	LangYaml   = "yaml"
	LangYml    = "yml"

	RootDefault = "Root"
)

const (
	StrInt        = "int"
	StrInt8       = "int8"
	StrInt16      = "int16"
	StrInt32      = "int32"
	StrInt64      = "int64"
	StrUint       = "uint"
	StrUint8      = "uint8"
	StrUint16     = "uint16"
	StrUint32     = "uint32"
	StrUint64     = "uint64"
	StrSint32     = "sint32"
	StrSint64     = "sint64"
	StrI16        = "i16"
	StrI32        = "i32"
	StrI64        = "i64"
	StrFixed32    = "fixed32"
	StrFixed64    = "fixed64"
	StrFloat32    = "float32"
	StrFloat64    = "float64"
	StrDouble     = "double"
	StrFloat      = "float"
	StrBool       = "bool"
	StrString     = "string"
	StrComplex64  = "complex64"
	StrComplex128 = "complex128"
	StrByte       = "byte"
	StrBytes      = "bytes"
	StrRune       = "rune"
	StrUintptr    = "uintptr"
	StrAny        = "any"
	StrPbAny      = "google.protobuf.Any"
	StrBinary     = "binary"
	StrMap        = "map"
	StrList       = "list"
	StrSet        = "set"
	StrNil        = "nil"
	StrNull       = "null"
	StrNumber     = "number"
	StrRepeated   = "repeated"
)

var (
	SourceLangs = []Lang{
		{
			Lang: LangJson,
		},
		{
			Lang:    LangYaml,
			Aliases: []string{LangYml},
		},
		{
			Lang: LangProto,
		},
		{
			Lang: LangThrift,
		},
		{
			Lang: LangGo,
		},
		{
			Lang: LangCsv,
		},
	}

	DestinationLangs = []Lang{
		{
			Lang: LangGo,
		},
		{
			Lang: LangProto,
		},
		{
			Lang: LangThrift,
		},
	}
	LangTmplMap = map[string]string{
		LangGo:     tmpl.Go,
		LangProto:  tmpl.Proto,
		LangThrift: tmpl.Thrift,
	}
)
