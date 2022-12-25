package st2

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonParser_Parse(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *JsonParser
		inspect func(r *JsonParser, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      []*Struct
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "empty",
			init: func(t *testing.T) *JsonParser {
				return NewJsonParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte("")),
				}
			},
			want1:   nil,
			wantErr: false,
		},
		{
			name: "illegal json",
			init: func(t *testing.T) *JsonParser {
				return NewJsonParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte("a")),
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.EqualError(t, err, "invalid character 'a' looking for beginning of value")
			},
		},
		{
			name: "simple struct",
			init: func(t *testing.T) *JsonParser {
				return NewJsonParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte(`{"a":1}`)),
				}
			},
			want1: []*Struct{
				{
					Members: []*Member{
						{
							Field: "a",
							Type:  Float64Val,
							Index: 1,
						},
					},
					Type: &StructType{
						Name: "Root",
					},
				},
			},
			wantErr:    false,
			inspectErr: func(err error, t *testing.T) {},
		},
		{
			name: "complex struct",
			init: func(t *testing.T) *JsonParser {
				return NewJsonParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte(`
{
	"a": {
		"b": 1,
		"c": "hello"
	},
	"b": {
		"b": 2,
		"c": "world"
	},
	"cccc": ["123"],
	"d": [{
		"b": 3,
		"c": 4
	}],
	"e": {
		"aa": true,
		"bb": false
	},
	"f": {
		"a": {
			"hello": true
		}
	},
	"ggg": [[{
		"ggg": 1
	}]]
}`)),
				}
			},
			want1: []*Struct{
				{
					Members: []*Member{
						{
							Field: "b",
							Type:  Float64Val,
							Index: 1,
						},
						{
							Field: "c",
							Type:  StringVal,
							Index: 2,
						},
					},
					Type: &StructType{
						Name: "A",
					},
				},
				{
					Members: []*Member{
						{
							Field: "b",
							Type:  Float64Val,
							Index: 1,
						},
						{
							Field: "c",
							Type:  Float64Val,
							Index: 2,
						},
					},
					Type: &StructType{
						Name: "D",
					},
				},
				{
					Members: []*Member{
						{
							Field: "aa",
							Type:  BoolVal,
							Index: 1,
						},
						{
							Field: "bb",
							Type:  BoolVal,
							Index: 2,
						},
					},
					Type: &StructType{
						Name: "E",
					},
				},
				{
					Members: []*Member{
						{
							Field: "hello",
							Type:  BoolVal,
							Index: 1,
						},
					},
					Type: &StructType{
						Name: "A01",
					},
				},
				{
					Members: []*Member{
						{
							Field: "a",
							Type: &StructType{
								Name: "A01",
							},
							Index: 1,
						},
					},
					Type: &StructType{
						Name: "F",
					},
				},
				{
					Members: []*Member{
						{
							Field: "ggg",
							Type:  Float64Val,
							Index: 1,
						},
					},
					Type: &StructType{
						Name: "Ggg",
					},
				},
				{
					Members: []*Member{
						{
							Field: "a",
							Type: &StructType{
								Name: "A",
							},
							Index: 1,
						},
						{
							Field: "b",
							Type: &StructType{
								Name: "A",
							},
							Index: 2,
						},
						{
							Field: "cccc",
							Type: &ArrayType{
								ChildType: StringVal,
							},
							Index: 3,
						},
						{
							Field: "d",
							Type: &ArrayType{
								ChildType: &StructType{
									Name: "D",
								},
							},
							Index: 4,
						},
						{
							Field: "e",
							Type: &StructType{
								Name: "E",
							},
							Index: 5,
						},
						{
							Field: "f",
							Type: &StructType{
								Name: "F",
							},
							Index: 6,
						},
						{
							Field: "ggg",
							Type: &ArrayType{
								ChildType: &ArrayType{
									ChildType: &StructType{
										Name: "Ggg",
									},
								},
							},
							Index: 7,
						},
					},
					Type: &StructType{
						Name: "Root",
					},
				},
			},
			wantErr:    false,
			inspectErr: func(err error, t *testing.T) {},
		},
		{
			name: "complex array struct",
			init: func(t *testing.T) *JsonParser {
				return NewJsonParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte(`
[{
	"a": {
		"b": 1,
		"c": "hello"
	},
	"b": {
		"b": 2,
		"c": "world"
	},
	"cccc": ["123"],
	"d": [{
		"b": 3,
		"c": 4
	}],
	"e": {
		"aa": true,
		"bb": false
	},
	"f": {
		"a": {
			"hello": true
		}
	},
	"ggg": [[{
		"ggg": 1
	}]],
	"hh": []
}]`)),
				}
			},
			want1: []*Struct{
				{
					Members: []*Member{
						{
							Field: "b",
							Type:  Float64Val,
							Index: 1,
						},
						{
							Field: "c",
							Type:  StringVal,
							Index: 2,
						},
					},
					Type: &StructType{
						Name: "A",
					},
				},
				{
					Members: []*Member{
						{
							Field: "b",
							Type:  Float64Val,
							Index: 1,
						},
						{
							Field: "c",
							Type:  Float64Val,
							Index: 2,
						},
					},
					Type: &StructType{
						Name: "D",
					},
				},
				{
					Members: []*Member{
						{
							Field: "aa",
							Type:  BoolVal,
							Index: 1,
						},
						{
							Field: "bb",
							Type:  BoolVal,
							Index: 2,
						},
					},
					Type: &StructType{
						Name: "E",
					},
				},
				{
					Members: []*Member{
						{
							Field: "hello",
							Type:  BoolVal,
							Index: 1,
						},
					},
					Type: &StructType{
						Name: "A01",
					},
				},
				{
					Members: []*Member{
						{
							Field: "a",
							Type: &StructType{
								Name: "A01",
							},
							Index: 1,
						},
					},
					Type: &StructType{
						Name: "F",
					},
				},
				{
					Members: []*Member{
						{
							Field: "ggg",
							Type:  Float64Val,
							Index: 1,
						},
					},
					Type: &StructType{
						Name: "Ggg",
					},
				},
				{
					Members: []*Member{
						{
							Field: "a",
							Type: &StructType{
								Name: "A",
							},
							Index: 1,
						},
						{
							Field: "b",
							Type: &StructType{
								Name: "A",
							},
							Index: 2,
						},
						{
							Field: "cccc",
							Type: &ArrayType{
								ChildType: StringVal,
							},
							Index: 3,
						},
						{
							Field: "d",
							Type: &ArrayType{
								ChildType: &StructType{
									Name: "D",
								},
							},
							Index: 4,
						},
						{
							Field: "e",
							Type: &StructType{
								Name: "E",
							},
							Index: 5,
						},
						{
							Field: "f",
							Type: &StructType{
								Name: "F",
							},
							Index: 6,
						},
						{
							Field: "ggg",
							Type: &ArrayType{
								ChildType: &ArrayType{
									ChildType: &StructType{
										Name: "Ggg",
									},
								},
							},
							Index: 7,
						},
					},
					Type: &StructType{
						Name: "Root",
					},
				},
			},
			wantErr:    false,
			inspectErr: func(err error, t *testing.T) {},
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
				t.Errorf("JsonParser.Parse got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("JsonParser.Parse error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
