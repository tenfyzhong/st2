package st2

var (
	NullVal   Type = &NullType{}
	BoolVal   Type = &BoolType{}
	NumberVal Type = &NumberType{}
	StringVal Type = &StringType{}
	ArrayVal  Type = &ArrayType{}
	StructVal Type = &StructType{}
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

type NumberType struct{}

func (v NumberType) Json() string     { return "number" }
func (v NumberType) Go() string       { return "float64" }
func (v NumberType) Protobuf() string { return "double" }
func (v NumberType) Thrift() string   { return "double" }

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
