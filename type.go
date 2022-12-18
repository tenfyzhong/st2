package st2

var (
	NullVal    Type = &NullType{}
	BoolVal    Type = &BoolType{}
	Float64Val Type = &Float64Type{}
	StringVal  Type = &StringType{}
	ArrayVal   Type = &ArrayType{}
	StructVal  Type = &StructType{}
)

type Type interface {
	Json() string
	Go() string
	Protobuf() string
	Thrift() string
}

type NullType struct{}

func (v NullType) Json() string     { return "null" }
func (v NullType) Go() string       { return "nil" }
func (v NullType) Protobuf() string { return "" }
func (v NullType) Thrift() string   { return "" }

type BoolType struct{}

func (v BoolType) Json() string     { return "bool" }
func (v BoolType) Go() string       { return "bool" }
func (v BoolType) Protobuf() string { return "bool" }
func (v BoolType) Thrift() string   { return "bool" }

type Float64Type struct{}

func (v Float64Type) Json() string     { return "number" }
func (v Float64Type) Go() string       { return "float64" }
func (v Float64Type) Protobuf() string { return "double" }
func (v Float64Type) Thrift() string   { return "double" }

type StringType struct{}

func (v StringType) Json() string     { return "string" }
func (v StringType) Go() string       { return "string" }
func (v StringType) Protobuf() string { return "string" }
func (v StringType) Thrift() string   { return "string" }

type ArrayType struct {
	ChildType Type
}

func (v ArrayType) Json() string     { return "[]" + v.ChildType.Json() }
func (v ArrayType) Go() string       { return "[]" + v.ChildType.Go() }
func (v ArrayType) Protobuf() string { return "repeated " + v.ChildType.Protobuf() }
func (v ArrayType) Thrift() string   { return "list<" + v.ChildType.Thrift() + ">" }

type StructType struct {
	Name string
}

func (v StructType) Json() string       { return v.Name }
func (v StructType) Go() string         { return "*" + v.Name }
func (v StructType) Protobuf() string   { return v.Name }
func (v StructType) Thrift() string     { return v.Name }
func (v StructType) StructName() string { return v.Name }
