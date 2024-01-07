package st2

import (
	"fmt"
	"strconv"
	"strings"
)

type StructLikeSource int

const (
	SLSUnknown StructLikeSource = 0
	SLSStruct  StructLikeSource = 1
	SLSUnion   StructLikeSource = 2
)

var (
	AnyVal        Type = &AnyType{}
	BoolVal       Type = &BoolType{}
	Float32Val    Type = &Float32Type{}
	Float64Val    Type = &Float64Type{}
	StringVal     Type = &StringType{}
	ArrayVal      Type = &ArrayType{}
	Int8Val       Type = &Int8Type{}
	Int16Val      Type = &Int16Type{}
	Int32Val      Type = &Int32Type{}
	Int64Val      Type = &Int64Type{}
	Uint8Val      Type = &Uint8Type{}
	Uint16Val     Type = &Uint16Type{}
	Uint32Val     Type = &Uint32Type{}
	Uint64Val     Type = &Uint64Type{}
	BinaryVal     Type = &BinaryType{}
	MapVal        Type = &MapType{}
	SetVal        Type = &SetType{}
	EnumVal       Type = &EnumType{}
	StructLikeVal Type = &StructLikeType{}
)

type Type interface {
	Json() string
	Go() string
	Proto() string
	Thrift() string
	IsBasicType() bool
}

type AnyType struct{}

func (v AnyType) Json() string      { return StrNull }
func (v AnyType) Go() string        { return StrAny }
func (v AnyType) Proto() string     { return StrPbAny }
func (v AnyType) Thrift() string    { return StrBinary }
func (v AnyType) Value() string     { return StrNil }
func (v AnyType) IsBasicType() bool { return false }

type BoolType struct {
	V bool
}

func (v BoolType) Json() string      { return StrBool }
func (v BoolType) Go() string        { return StrBool }
func (v BoolType) Proto() string     { return StrBool }
func (v BoolType) Thrift() string    { return StrBool }
func (v BoolType) Value() string     { return strconv.FormatBool(v.V) }
func (v BoolType) IsBasicType() bool { return true }

type Float32Type struct {
	V float32
}

func (v Float32Type) Json() string      { return StrNumber }
func (v Float32Type) Go() string        { return StrFloat32 }
func (v Float32Type) Proto() string     { return StrFloat }
func (v Float32Type) Thrift() string    { return StrDouble }
func (v Float32Type) Value() string     { return strconv.FormatFloat(float64(v.V), 'f', -1, 32) }
func (v Float32Type) IsBasicType() bool { return true }

type Float64Type struct {
	V float64
}

func (v Float64Type) Json() string      { return StrNumber }
func (v Float64Type) Go() string        { return StrFloat64 }
func (v Float64Type) Proto() string     { return StrDouble }
func (v Float64Type) Thrift() string    { return StrDouble }
func (v Float64Type) Value() string     { return strconv.FormatFloat(v.V, 'f', -1, 64) }
func (v Float64Type) IsBasicType() bool { return true }

type StringType struct {
	V string
}

func (v StringType) Json() string      { return StrString }
func (v StringType) Go() string        { return StrString }
func (v StringType) Proto() string     { return StrString }
func (v StringType) Thrift() string    { return StrString }
func (v StringType) Value() string     { return v.V }
func (v StringType) IsBasicType() bool { return true }

type ArrayType struct {
	ChildType Type
}

func (v ArrayType) Json() string      { return "[]" + v.ChildType.Json() }
func (v ArrayType) Go() string        { return "[]" + v.ChildType.Go() }
func (v ArrayType) Proto() string     { return StrRepeated + " " + v.ChildType.Proto() }
func (v ArrayType) Thrift() string    { return StrList + "<" + v.ChildType.Thrift() + ">" }
func (v ArrayType) IsBasicType() bool { return false }

type Int8Type struct {
	V int8
}

func (v Int8Type) Json() string      { return StrNumber }
func (v Int8Type) Go() string        { return StrInt8 }
func (v Int8Type) Proto() string     { return StrInt32 }
func (v Int8Type) Thrift() string    { return StrByte }
func (v Int8Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int8Type) IsBasicType() bool { return true }

type Int16Type struct {
	V int16
}

func (v Int16Type) Json() string      { return StrNumber }
func (v Int16Type) Go() string        { return StrInt16 }
func (v Int16Type) Proto() string     { return StrInt32 }
func (v Int16Type) Thrift() string    { return StrI16 }
func (v Int16Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int16Type) IsBasicType() bool { return true }

type Int32Type struct {
	V int32
}

func (v Int32Type) Json() string      { return StrNumber }
func (v Int32Type) Go() string        { return StrInt32 }
func (v Int32Type) Proto() string     { return StrInt32 }
func (v Int32Type) Thrift() string    { return StrI32 }
func (v Int32Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int32Type) IsBasicType() bool { return true }

type Int64Type struct {
	V int64
}

func (v Int64Type) Json() string      { return StrNumber }
func (v Int64Type) Go() string        { return StrInt64 }
func (v Int64Type) Proto() string     { return StrInt64 }
func (v Int64Type) Thrift() string    { return StrI64 }
func (v Int64Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int64Type) IsBasicType() bool { return true }

type Uint8Type struct {
	V int8
}

func (v Uint8Type) Json() string      { return StrNumber }
func (v Uint8Type) Go() string        { return StrUint8 }
func (v Uint8Type) Proto() string     { return StrUint32 }
func (v Uint8Type) Thrift() string    { return StrByte }
func (v Uint8Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Uint8Type) IsBasicType() bool { return true }

type Uint16Type struct {
	V int16
}

func (v Uint16Type) Json() string      { return StrNumber }
func (v Uint16Type) Go() string        { return StrUint16 }
func (v Uint16Type) Proto() string     { return StrUint32 }
func (v Uint16Type) Thrift() string    { return StrI16 }
func (v Uint16Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Uint16Type) IsBasicType() bool { return true }

type Uint32Type struct {
	V int32
}

func (v Uint32Type) Json() string      { return StrNumber }
func (v Uint32Type) Go() string        { return StrUint32 }
func (v Uint32Type) Proto() string     { return StrUint32 }
func (v Uint32Type) Thrift() string    { return StrI32 }
func (v Uint32Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Uint32Type) IsBasicType() bool { return true }

type Uint64Type struct {
	V int64
}

func (v Uint64Type) Json() string      { return StrNumber }
func (v Uint64Type) Go() string        { return StrUint64 }
func (v Uint64Type) Proto() string     { return StrUint64 }
func (v Uint64Type) Thrift() string    { return StrI64 }
func (v Uint64Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Uint64Type) IsBasicType() bool { return true }

type BinaryType struct{}

func (v BinaryType) Json() string      { return StrString }
func (v BinaryType) Go() string        { return "[]" + StrByte }
func (v BinaryType) Proto() string     { return StrBytes }
func (v BinaryType) Thrift() string    { return StrBinary }
func (v BinaryType) IsBasicType() bool { return false }

type MapType struct {
	Key   Type
	Value Type
}

func (v MapType) Json() string { return "{}" }
func (v MapType) Go() string   { return fmt.Sprintf("map[%s]%s", v.Key.Go(), v.Value.Go()) }
func (v MapType) Proto() string {
	return fmt.Sprintf("map<%s, %s>", v.Key.Proto(), v.Value.Proto())
}
func (v MapType) Thrift() string    { return fmt.Sprintf("map<%s, %s>", v.Key.Thrift(), v.Value.Thrift()) }
func (v MapType) IsBasicType() bool { return false }

type SetType struct {
	Key Type
}

func (v SetType) Json() string      { return "{}" }
func (v SetType) Go() string        { return fmt.Sprintf("%s[%s]%s", StrMap, v.Key.Go(), StrBool) }
func (v SetType) Proto() string     { return fmt.Sprintf("%s<%s, %s>", StrMap, v.Key.Proto(), StrBool) }
func (v SetType) Thrift() string    { return fmt.Sprintf("%s<%s>", StrSet, v.Key.Thrift()) }
func (v SetType) IsBasicType() bool { return false }

type EnumType struct {
	Name string
}

func (v EnumType) Json() string             { return v.Name }
func (v EnumType) Go() string               { return v.Name }
func (v EnumType) Proto() string            { return v.Name }
func (v EnumType) Thrift() string           { return v.Name }
func (v EnumType) IsBasicType() bool        { return false }
func (v EnumType) StructName() string       { return v.Name }
func (v EnumType) GoStructType() string     { return "enum" }
func (v EnumType) ProtoStructType() string  { return "enum" }
func (v EnumType) ThriftStructType() string { return "enum" }

type StructLikeType struct {
	Name   string
	Source StructLikeSource
}

func (v StructLikeType) Json() string            { return v.Name }
func (v StructLikeType) Go() string              { return "*" + goWithPackageName(v.Name) }
func (v StructLikeType) Proto() string           { return v.Name }
func (v StructLikeType) Thrift() string          { return v.Name }
func (v StructLikeType) IsBasicType() bool       { return false }
func (v StructLikeType) StructName() string      { return v.Name }
func (v StructLikeType) GoStructType() string    { return "struct" }
func (v StructLikeType) ProtoStructType() string { return "message" }
func (v StructLikeType) ThriftStructType() string {
	switch v.Source {
	case SLSStruct:
		return "struct"
	case SLSUnion:
		return "union"
	}
	return "struct"
}

func goWithPackageName(name string) string {
	// If the name of the filed is in other package,
	// use the last part of the package name as go's package
	names := strings.Split(name, ".")
	if len(names) > 2 {
		names = names[len(names)-2:]
	}
	return strings.Join(names, ".")
}
