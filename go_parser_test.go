package st2

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoParser_Parse(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) GoParser
		inspect func(r GoParser, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      []*Struct
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "parse error",
			init: func(t *testing.T) GoParser {
				return *NewGoParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte("a")),
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.EqualError(t, err, "1:1: expected 'package', found a")
			},
		},
		{
			name: "succ",
			init: func(t *testing.T) GoParser {
				return *NewGoParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte(`
package main

// EEEE
type Eeee int // EEEE
type Efff Eeee
type IntList []int
type StringList []string
type IntIntMap map[int]int
type IntStringMap map[int]string
type IntEeeeMap map[int]Eeee
type IntCccMap map[int]Ccc
type IntArray [1]int

const (
	EEEA Eeee = 0 // a
	EEEB Eeee = 1 // a
	EEEC Eeee = 3 // a

	II1 int = 1
	II2 int = 2

	IN1 = 1
	IN2 = 2
)

// haha
type Aaa struct { // aaa
	// a
	A  []int32 // a
	B  int64   
	C  *string 
	MM map[int]string
}

type BbbBB struct {
	A int32
	B int64
	C string
}

type Ccc struct {
	A   int32
	B   int64
	C   string
	Aaa *Aaa
}

type ErrorStatus struct {
	Message string
	Details []*protobuf.Any
}

type SampleMessage struct {
}
`)),
				}
			},
			wantErr: false,
			want1: []*Struct{
				{
					Type: &EnumType{
						Name: "Eeee",
					},
					Members: []*Member{
						{
							Field: "EEEA",
							Type: &EnumType{
								Name: "Eeee",
							},
							Index: 0,
							Comment: Comment{
								InlineComment: "// a",
							},
						},
						{
							Field: "EEEB",
							Type: &EnumType{
								Name: "Eeee",
							},
							Index: 1,
							Comment: Comment{
								InlineComment: "// a",
							},
						},
						{
							Field: "EEEC",
							Type: &EnumType{
								Name: "Eeee",
							},
							Index: 3,
							Comment: Comment{
								InlineComment: "// a",
							},
						},
					},
				},
				{
					Type: &StructType{
						Name: "Aaa",
						Type: "struct",
					},
					Comment: Comment{
						BeginningComments: []string{"// haha"},
					},
					Members: []*Member{
						{
							Field: "A",
							Type: &ArrayType{
								ChildType: Int32Val,
							},
							Index: 1,
							Comment: Comment{
								BeginningComments: []string{"// a"},
								InlineComment:     "// a",
							},
						},
						{
							Field: "B",
							Type:  Int64Val,
							Index: 2,
						},
						{
							Field:    "C",
							Type:     StringVal,
							Index:    3,
							Optional: true,
						},
						{
							Field: "MM",
							Type: &MapType{
								Key:   Int64Val,
								Value: StringVal,
							},
							Index: 4,
						},
					},
				},
				{
					Type: &StructType{
						Name: "BbbBB",
						Type: "struct",
					},
					Members: []*Member{
						{
							Field: "A",
							Type:  Int32Val,
							Index: 1,
						},
						{
							Field: "B",
							Type:  Int64Val,
							Index: 2,
						},
						{
							Field: "C",
							Type:  StringVal,
							Index: 3,
						},
					},
				},
				{
					Type: &StructType{
						Name: "Ccc",
						Type: "struct",
					},
					Members: []*Member{
						{
							Field: "A",
							Type:  Int32Val,
							Index: 1,
						},
						{
							Field: "B",
							Type:  Int64Val,
							Index: 2,
						},
						{
							Field: "C",
							Type:  StringVal,
							Index: 3,
						},
						{
							Field: "Aaa",
							Type: &RawType{
								Name: "Aaa",
							},
							Index: 4,
						},
					},
				},
				{
					Type: &StructType{
						Name: "ErrorStatus",
						Type: "struct",
					},
					Members: []*Member{
						{
							Field: "Message",
							Type:  StringVal,
							Index: 1,
						},
						{
							Field: "Details",
							Type: &ArrayType{
								ChildType: &RawType{
									Name: "protobuf.Any",
								},
							},
							Index: 2,
						},
					},
				},
				{
					Type: &StructType{
						Name: "SampleMessage",
						Type: "struct",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			got1, err := receiver.Parse(tArgs.reader)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				got1json, _ := json.MarshalIndent(got1, "", "  ")
				want1json, _ := json.MarshalIndent(tt.want1, "", "  ")
				t.Errorf("GoParser.Parse got1 = %v, want1: %v", string(got1json), string(want1json))
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("GoParser.Parse error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
