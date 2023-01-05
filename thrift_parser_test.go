package st2

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThriftParser_Parse(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) ThriftParser
		inspect func(r ThriftParser, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      []*Struct
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "empty",
			init: func(t *testing.T) ThriftParser {
				return *NewThriftParser(Context{})
			},
			args: func(t *testing.T) args {
				data := []byte("")
				reader := bytes.NewReader(data)
				return args{
					reader: reader,
				}
			},
			want1:   nil,
			wantErr: false,
		},
		{
			name: "illegal",
			init: func(t *testing.T) ThriftParser {
				return *NewThriftParser(Context{})
			},
			args: func(t *testing.T) args {
				data := []byte("a")
				reader := bytes.NewReader(data)
				return args{
					reader: reader,
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.EqualError(t, err, "\nparse error near Unknown (line 1 symbol 1 - line 1 symbol 1):\n\"\"\n")
			},
		},
		{
			name: "succ",
			init: func(t *testing.T) ThriftParser {
				return *NewThriftParser(Context{})
			},
			args: func(t *testing.T) args {
				data := []byte(`
enum EEE {
    A = 1;
    B = 2;
}

struct SS {
    1: optional bool a,
    2: byte b,
    3: i16 c,
    4: i32 d,
    5: i64 e,
    6: double f,
    7: string g,
    8: binary h,
    9: map<i32, i32> i,
    10: optional list<i32> j,
    11: set<i32> k,
}

struct AAA {
    1: string hello,
}

struct BBB {
    1: i16 b1,
    2: i32 b2,
    3: EEE e,
    4: map<AAA, BBB> mapab,
    5: set<AAA> seta,
    6: list<BBB> listb,
}

union UUU {
    1: AAA a;
    2: BBB b;
}
`)
				reader := bytes.NewReader(data)
				return args{
					reader: reader,
				}
			},
			want1: []*Struct{
				{
					Type: &EnumType{
						Name: "EEE",
					},
					Members: []*Member{
						{
							Field: "A",
							Type: &EnumType{
								Name: "EEE",
							},
							Index:    1,
							Optional: false,
						},
						{
							Field: "B",
							Type: &EnumType{
								Name: "EEE",
							},
							Index:    2,
							Optional: false,
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "SS",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field:    "a",
							Type:     BoolVal,
							Index:    1,
							Optional: true,
						},
						{
							Field:    "b",
							Type:     Int8Val,
							Index:    2,
							Optional: false,
						},
						{
							Field:    "c",
							Type:     Int16Val,
							Index:    3,
							Optional: false,
						},
						{
							Field:    "d",
							Type:     Int32Val,
							Index:    4,
							Optional: false,
						},
						{
							Field:    "e",
							Type:     Int64Val,
							Index:    5,
							Optional: false,
						},
						{
							Field:    "f",
							Type:     Float64Val,
							Index:    6,
							Optional: false,
						},
						{
							Field:    "g",
							Type:     StringVal,
							Index:    7,
							Optional: false,
						},
						{
							Field:    "h",
							Type:     BinaryVal,
							Index:    8,
							Optional: false,
						},
						{
							Field: "i",
							Type: &MapType{
								Key:   Int32Val,
								Value: Int32Val,
							},
							Index:    9,
							Optional: false,
						},
						{
							Field: "j",
							Type: &ArrayType{
								ChildType: Int32Val,
							},
							Index:    10,
							Optional: true,
						},
						{
							Field: "k",
							Type: &SetType{
								Key: Int32Val,
							},
							Index:    11,
							Optional: false,
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "AAA",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "hello",
							Type:  StringVal,
							Index: 1,
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "BBB",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "b1",
							Type:  Int16Val,
							Index: 1,
						},
						{
							Field: "b2",
							Type:  Int32Val,
							Index: 2,
						},
						{
							Field: "e",
							Type: &StructLikeType{
								Name: "EEE",
							},
							Index: 3,
						},
						{
							Field: "mapab",
							Type: &MapType{
								Key: &StructLikeType{
									Name: "AAA",
								},
								Value: &StructLikeType{
									Name: "BBB",
								},
							},
							Index: 4,
						},
						{
							Field: "seta",
							Type: &SetType{
								Key: &StructLikeType{
									Name: "AAA",
								},
							},
							Index: 5,
						},
						{
							Field: "listb",
							Type: &ArrayType{
								ChildType: &StructLikeType{
									Name: "BBB",
								},
							},
							Index: 6,
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "UUU",
						Source: SLSUnion,
					},
					Members: []*Member{
						{
							Field: "a",
							Type: &StructLikeType{
								Name: "AAA",
							},
							Index: 1,
						},
						{
							Field: "b",
							Type: &StructLikeType{
								Name: "BBB",
							},
							Index: 2,
						},
					},
				},
			},
			wantErr: false,
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
				got1Json, _ := json.MarshalIndent(got1, "", "  ")
				want1Json, _ := json.MarshalIndent(tt.want1, "", "  ")
				t.Errorf("ThriftParser.Parse got1 = %v, want1: %v", string(got1Json), string(want1Json))
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("ThriftParser.Parse error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
