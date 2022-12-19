package st2

import "fmt"
import "strconv"

var (
	NullVal    Type = &NullType{}
	BoolVal    Type = &BoolType{}
	Float64Val Type = &Float64Type{}
	StringVal  Type = &StringType{}
	ArrayVal   Type = &ArrayType{}
	StructVal  Type = &StructType{}
	Int8Val    Type = &Int8Type{}
	Int16Val   Type = &Int16Type{}
	Int32Val   Type = &Int32Type{}
	Int64Val   Type = &Int64Type{}
	BinaryVal  Type = &BinaryType{}
	MapVal     Type = &MapType{}
	SetVal     Type = &SetType{}
	EnumVal    Type = &EnumType{}
)

type Type interface {
	Json() string
	Go() string
	Protobuf() string
	Thrift() string
	IsBasicType() bool
}

type NullType struct{}

func (v NullType) Json() string      { return "null" }
func (v NullType) Go() string        { return "nil" }
func (v NullType) Protobuf() string  { return "" }
func (v NullType) Thrift() string    { return "" }
func (v NullType) Value() string     { return "nil" }
func (v NullType) IsBasicType() bool { return false }

type BoolType struct {
	V bool
}

func (v BoolType) Json() string      { return "bool" }
func (v BoolType) Go() string        { return "bool" }
func (v BoolType) Protobuf() string  { return "bool" }
func (v BoolType) Thrift() string    { return "bool" }
func (v BoolType) Value() string     { return strconv.FormatBool(v.V) }
func (v BoolType) IsBasicType() bool { return true }

type Float64Type struct {
	V float64
}

func (v Float64Type) Json() string      { return "number" }
func (v Float64Type) Go() string        { return "float64" }
func (v Float64Type) Protobuf() string  { return "double" }
func (v Float64Type) Thrift() string    { return "double" }
func (v Float64Type) Value() string     { return strconv.FormatFloat(v.V, 'f', -1, 64) }
func (v Float64Type) IsBasicType() bool { return true }

type StringType struct {
	V string
}

func (v StringType) Json() string      { return "string" }
func (v StringType) Go() string        { return "string" }
func (v StringType) Protobuf() string  { return "string" }
func (v StringType) Thrift() string    { return "string" }
func (v StringType) Value() string     { return v.V }
func (v StringType) IsBasicType() bool { return true }

type ArrayType struct {
	ChildType Type
}

func (v ArrayType) Json() string      { return "[]" + v.ChildType.Json() }
func (v ArrayType) Go() string        { return "[]" + v.ChildType.Go() }
func (v ArrayType) Protobuf() string  { return "repeated " + v.ChildType.Protobuf() }
func (v ArrayType) Thrift() string    { return "list<" + v.ChildType.Thrift() + ">" }
func (v ArrayType) IsBasicType() bool { return false }

type StructType struct {
	Name string
	Type string // struct or union, default struct
}

func (v StructType) Json() string               { return v.Name }
func (v StructType) Go() string                 { return "*" + v.Name }
func (v StructType) Protobuf() string           { return v.Name }
func (v StructType) Thrift() string             { return v.Name }
func (v StructType) IsBasicType() bool          { return false }
func (v StructType) StructName() string         { return v.Name }
func (v StructType) GoStructType() string       { return "struct" }
func (v StructType) ProtobufStructType() string { return "message" }
func (v StructType) ThriftStructType() string {
	if v.Type != "" {
		return v.Type
	}
	return "struct"
}

type Int8Type struct {
	V int8
}

func (v Int8Type) Json() string      { return "number" }
func (v Int8Type) Go() string        { return "int8" }
func (v Int8Type) Protobuf() string  { return "int32" }
func (v Int8Type) Thrift() string    { return "byte" }
func (v Int8Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int8Type) IsBasicType() bool { return true }

type Int16Type struct {
	V int16
}

func (v Int16Type) Json() string      { return "number" }
func (v Int16Type) Go() string        { return "int16" }
func (v Int16Type) Protobuf() string  { return "int32" }
func (v Int16Type) Thrift() string    { return "i16" }
func (v Int16Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int16Type) IsBasicType() bool { return true }

type Int32Type struct {
	V int32
}

func (v Int32Type) Json() string      { return "number" }
func (v Int32Type) Go() string        { return "int32" }
func (v Int32Type) Protobuf() string  { return "int32" }
func (v Int32Type) Thrift() string    { return "i32" }
func (v Int32Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int32Type) IsBasicType() bool { return true }

type Int64Type struct {
	V int64
}

func (v Int64Type) Json() string      { return "number" }
func (v Int64Type) Go() string        { return "int64" }
func (v Int64Type) Protobuf() string  { return "int64" }
func (v Int64Type) Thrift() string    { return "i64" }
func (v Int64Type) Value() string     { return strconv.FormatInt(int64(v.V), 10) }
func (v Int64Type) IsBasicType() bool { return true }

type BinaryType struct{}

func (v BinaryType) Json() string      { return "string" }
func (v BinaryType) Go() string        { return "[]byte" }
func (v BinaryType) Protobuf() string  { return "bytes" }
func (v BinaryType) Thrift() string    { return "binary" }
func (v BinaryType) IsBasicType() bool { return false }

type MapType struct {
	Key   Type
	Value Type
}

func (v MapType) Json() string { return "{}" }
func (v MapType) Go() string   { return fmt.Sprintf("map[%s]%s", v.Key.Go(), v.Value.Go()) }
func (v MapType) Protobuf() string {
	return fmt.Sprintf("map<%s, %s>", v.Key.Protobuf(), v.Value.Protobuf())
}
func (v MapType) Thrift() string    { return fmt.Sprintf("map<%s, %s>", v.Key.Thrift(), v.Value.Thrift()) }
func (v MapType) IsBasicType() bool { return false }

type SetType struct {
	Key Type
}

func (v SetType) Json() string      { return "{}" }
func (v SetType) Go() string        { return fmt.Sprintf("map[%s]bool", v.Key.Go()) }
func (v SetType) Protobuf() string  { return fmt.Sprintf("map<%s, bool>", v.Key.Protobuf()) }
func (v SetType) Thrift() string    { return fmt.Sprintf("set<%s>", v.Key.Thrift()) }
func (v SetType) IsBasicType() bool { return false }

type EnumType struct {
	Name string
}

func (v EnumType) Json() string               { return v.Name }
func (v EnumType) Go() string                 { return v.Name }
func (v EnumType) Protobuf() string           { return v.Name }
func (v EnumType) Thrift() string             { return v.Name }
func (v EnumType) IsBasicType() bool          { return false }
func (v EnumType) StructName() string         { return v.Name }
func (v EnumType) GoStructType() string       { return "enum" }
func (v EnumType) ProtobufStructType() string { return "enum" }
func (v EnumType) ThriftStructType() string   { return "enum" }

type RawType struct {
	Name string
}

func (v RawType) Json() string      { return v.Name }
func (v RawType) Go() string        { return "*" + v.Name }
func (v RawType) Protobuf() string  { return v.Name }
func (v RawType) Thrift() string    { return v.Name }
func (v RawType) IsBasicType() bool { return false }
