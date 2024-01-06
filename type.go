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
	NullVal       Type = &NullType{}
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

type NullType struct{}

func (v NullType) Json() string      { return "null" }
func (v NullType) Go() string        { return "interface{}" }
func (v NullType) Proto() string     { return "" }
func (v NullType) Thrift() string    { return "" }
func (v NullType) Value() string     { return "nil" }
func (v NullType) IsBasicType() bool { return false }

type BoolType struct {
	V bool
}

func (v BoolType) Json() string      { return "bool" }
func (v BoolType) Go() string        { return "bool" }
func (v BoolType) Proto() string     { return "bool" }
func (v BoolType) Thrift() string    { return "bool" }
func (v BoolType) Value() string     { return strconv.FormatBool(v.V) }
func (v BoolType) IsBasicType() bool { return true }

type Float32Type struct {
	V float32
}

func (v Float32Type) Json() string      { return "number" }
func (v Float32Type) Go() string        { return "float32" }
func (v Float32Type) Proto() string     { return "float" }
func (v Float32Type) Thrift() string    { return "double" }
func (v Float32Type) Value() string     { return strconv.FormatFloat(float64(v.V), 'f', -1, 32) }
func (v Float32Type) IsBasicType() bool { return true }

type Float64Type struct {
	V float64
}

func (v Float64Type) Json() string      { return "number" }
func (v Float64Type) Go() string        { return "float64" }
func (v Float64Type) Proto() string     { return "double" }
func (v Float64Type) Thrift() string    { return "double" }
func (v Float64Type) Value() string     { return strconv.FormatFloat(v.V, 'f', -1, 64) }
func (v Float64Type) IsBasicType() bool { return true }

type StringType struct {
	V string
}

func (v StringType) Json() string      { return "string" }
func (v StringType) Go() string        { return "string" }
func (v StringType) Proto() string     { return "string" }
func (v StringType) Thrift() string    { return "string" }
func (v StringType) Value() string     { return v.V }
func (v StringType) IsBasicType() bool { return true }

type ArrayType struct {
	ChildType Type
}

func (v ArrayType) Json() string      { return "[]" + v.ChildType.Json() }
func (v ArrayType) Go() string        { return "[]" + v.ChildType.Go() }
func (v ArrayType) Proto() string     { return "repeated " + v.ChildType.Proto() }
func (v ArrayType) Thrift() string    { return "list<" + v.ChildType.Thrift() + ">" }
func (v ArrayType) IsBasicType() bool { return false }

type Int8Type struct {
	V int8
}

func (v Int8Type) Json() string      { return "number" }
func (v Int8Type) Go() string        { return "int8" }
func (v Int8Type) Proto() string     { return "int32" }
func (v Int8Type) Thrift() string    { return "byte" }
func (v Int8Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int8Type) IsBasicType() bool { return true }

type Int16Type struct {
	V int16
}

func (v Int16Type) Json() string      { return "number" }
func (v Int16Type) Go() string        { return "int16" }
func (v Int16Type) Proto() string     { return "int32" }
func (v Int16Type) Thrift() string    { return "i16" }
func (v Int16Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int16Type) IsBasicType() bool { return true }

type Int32Type struct {
	V int32
}

func (v Int32Type) Json() string      { return "number" }
func (v Int32Type) Go() string        { return "int32" }
func (v Int32Type) Proto() string     { return "int32" }
func (v Int32Type) Thrift() string    { return "i32" }
func (v Int32Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int32Type) IsBasicType() bool { return true }

type Int64Type struct {
	V int64
}

func (v Int64Type) Json() string      { return "number" }
func (v Int64Type) Go() string        { return "int64" }
func (v Int64Type) Proto() string     { return "int64" }
func (v Int64Type) Thrift() string    { return "i64" }
func (v Int64Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int64Type) IsBasicType() bool { return true }

type Uint8Type struct {
	V int8
}

func (v Uint8Type) Json() string      { return "number" }
func (v Uint8Type) Go() string        { return "uint8" }
func (v Uint8Type) Proto() string     { return "uint32" }
func (v Uint8Type) Thrift() string    { return "byte" }
func (v Uint8Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Uint8Type) IsBasicType() bool { return true }

type Uint16Type struct {
	V int16
}

func (v Uint16Type) Json() string      { return "number" }
func (v Uint16Type) Go() string        { return "uint16" }
func (v Uint16Type) Proto() string     { return "uint32" }
func (v Uint16Type) Thrift() string    { return "i16" }
func (v Uint16Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Uint16Type) IsBasicType() bool { return true }

type Uint32Type struct {
	V int32
}

func (v Uint32Type) Json() string      { return "number" }
func (v Uint32Type) Go() string        { return "uint32" }
func (v Uint32Type) Proto() string     { return "uint32" }
func (v Uint32Type) Thrift() string    { return "i32" }
func (v Uint32Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Uint32Type) IsBasicType() bool { return true }

type Uint64Type struct {
	V int64
}

func (v Uint64Type) Json() string      { return "number" }
func (v Uint64Type) Go() string        { return "uint64" }
func (v Uint64Type) Proto() string     { return "uint64" }
func (v Uint64Type) Thrift() string    { return "i64" }
func (v Uint64Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Uint64Type) IsBasicType() bool { return true }

type BinaryType struct{}

func (v BinaryType) Json() string      { return "string" }
func (v BinaryType) Go() string        { return "[]byte" }
func (v BinaryType) Proto() string     { return "bytes" }
func (v BinaryType) Thrift() string    { return "binary" }
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
func (v SetType) Go() string        { return fmt.Sprintf("map[%s]bool", v.Key.Go()) }
func (v SetType) Proto() string     { return fmt.Sprintf("map<%s, bool>", v.Key.Proto()) }
func (v SetType) Thrift() string    { return fmt.Sprintf("set<%s>", v.Key.Thrift()) }
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
