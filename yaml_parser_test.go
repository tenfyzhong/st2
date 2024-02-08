package st2

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestYamlParser_Parse(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *StructuredParser
		inspect func(r *StructuredParser, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      []*Struct
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "empty",
			init: func(t *testing.T) *StructuredParser {
				return NewYamlParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte("")),
				}
			},
			want1:   nil,
			wantErr: false,
		},
		// {
		// 	name: "illegal yaml",
		// 	init: func(t *testing.T) *StructuredParser {
		// 		return NewYamlParser(Context{})
		// 	},
		// 	args: func(t *testing.T) args {
		// 		return args{
		// 			reader: bytes.NewReader([]byte("a")),
		// 		}
		// 	},
		// 	wantErr: true,
		// 	inspectErr: func(err error, t *testing.T) {
		// 		assert.EqualError(t, err, "Read: unexpected value type: 0, error found in #0 byte of ...|a|..., bigger context ...|a|...")
		// 	},
		// },
		{
			name: "simple struct",
			init: func(t *testing.T) *StructuredParser {
				return NewYamlParser(Context{
					Root: "helloWorld",
				})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte(`a: 1`)),
				}
			},
			want1: []*Struct{
				{
					Members: []*Member{
						{
							Field: "a",
							Type:  Int64Val,
							Index: 1,
							GoTag: []string{`yaml:"a"`},
						},
					},
					Type: &StructLikeType{
						Name:   "HelloWorld",
						Source: SLSStruct,
					},
				},
			},
			wantErr:    false,
			inspectErr: func(err error, t *testing.T) {},
		},
		{
			name: "simple struct with null",
			init: func(t *testing.T) *StructuredParser {
				return NewYamlParser(Context{
					Root: "helloWorld",
				})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte(`{"a":1, "b":null, "c": 1.1}`)),
				}
			},
			want1: []*Struct{
				{
					Members: []*Member{
						{
							Field: "a",
							Type:  Int64Val,
							Index: 1,
							GoTag: []string{`yaml:"a"`},
						},
						{
							Field: "b",
							Type:  AnyVal,
							Index: 2,
							GoTag: []string{`yaml:"b"`},
						},
						{
							Field: "c",
							Type:  Float64Val,
							Index: 3,
							GoTag: []string{`yaml:"c"`},
						},
					},
					Type: &StructLikeType{
						Name:   "HelloWorld",
						Source: SLSStruct,
					},
				},
			},
			wantErr:    false,
			inspectErr: func(err error, t *testing.T) {},
		},
		{
			name: "complex struct",
			init: func(t *testing.T) *StructuredParser {
				return NewYamlParser(Context{})
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
					Type: &StructLikeType{
						Name:   "A",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "b",
							Type:  Int64Val,
							Index: 1,
							GoTag: []string{`yaml:"b"`},
						},
						{
							Field: "c",
							Type:  StringVal,
							Index: 2,
							GoTag: []string{`yaml:"c"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "D",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "b",
							Type:  Int64Val,
							Index: 1,
							GoTag: []string{`yaml:"b"`},
						},
						{
							Field: "c",
							Type:  Int64Val,
							Index: 2,
							GoTag: []string{`yaml:"c"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "E",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "aa",
							Type:  BoolVal,
							Index: 1,
							GoTag: []string{`yaml:"aa"`},
						},
						{
							Field: "bb",
							Type:  BoolVal,
							Index: 2,
							GoTag: []string{`yaml:"bb"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "A01",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "hello",
							Type:  BoolVal,
							Index: 1,
							GoTag: []string{`yaml:"hello"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "F",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "a",
							Type: &StructLikeType{
								Name: "A01",
							},
							Index: 1,
							GoTag: []string{`yaml:"a"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "Ggg",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "ggg",
							Type:  Int64Val,
							Index: 1,
							GoTag: []string{`yaml:"ggg"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "Root",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "a",
							Type: &StructLikeType{
								Name: "A",
							},
							Index: 1,
							GoTag: []string{`yaml:"a"`},
						},
						{
							Field: "b",
							Type: &StructLikeType{
								Name: "A",
							},
							Index: 2,
							GoTag: []string{`yaml:"b"`},
						},
						{
							Field: "cccc",
							Type: &ArrayType{
								ChildType: StringVal,
							},
							Index: 3,
							GoTag: []string{`yaml:"cccc"`},
						},
						{
							Field: "d",
							Type: &ArrayType{
								ChildType: &StructLikeType{
									Name: "D",
								},
							},
							Index: 4,
							GoTag: []string{`yaml:"d"`},
						},
						{
							Field: "e",
							Type: &StructLikeType{
								Name: "E",
							},
							Index: 5,
							GoTag: []string{`yaml:"e"`},
						},
						{
							Field: "f",
							Type: &StructLikeType{
								Name: "F",
							},
							Index: 6,
							GoTag: []string{`yaml:"f"`},
						},
						{
							Field: "ggg",
							Type: &ArrayType{
								ChildType: &ArrayType{
									ChildType: &StructLikeType{
										Name: "Ggg",
									},
								},
							},
							Index: 7,
							GoTag: []string{`yaml:"ggg"`},
						},
					},
				},
			},
			wantErr:    false,
			inspectErr: func(err error, t *testing.T) {},
		},
		{
			name: "complex array struct",
			init: func(t *testing.T) *StructuredParser {
				return NewYamlParser(Context{})
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
	"hh": [],
	"jj": [[]]
}]`)),
				}
			},
			want1: []*Struct{
				{
					Type: &StructLikeType{
						Name:   "A",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "b",
							Type:  Int64Val,
							Index: 1,
							GoTag: []string{`yaml:"b"`},
						},
						{
							Field: "c",
							Type:  StringVal,
							Index: 2,
							GoTag: []string{`yaml:"c"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "D",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "b",
							Type:  Int64Val,
							Index: 1,
							GoTag: []string{`yaml:"b"`},
						},
						{
							Field: "c",
							Type:  Int64Val,
							Index: 2,
							GoTag: []string{`yaml:"c"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "E",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "aa",
							Type:  BoolVal,
							Index: 1,
							GoTag: []string{`yaml:"aa"`},
						},
						{
							Field: "bb",
							Type:  BoolVal,
							Index: 2,
							GoTag: []string{`yaml:"bb"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "A01",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "hello",
							Type:  BoolVal,
							Index: 1,
							GoTag: []string{`yaml:"hello"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "F",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "a",
							Type: &StructLikeType{
								Name: "A01",
							},
							Index: 1,
							GoTag: []string{`yaml:"a"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "Ggg",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "ggg",
							Type:  Int64Val,
							Index: 1,
							GoTag: []string{`yaml:"ggg"`},
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "Root",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "a",
							Type: &StructLikeType{
								Name: "A",
							},
							Index: 1,
							GoTag: []string{`yaml:"a"`},
						},
						{
							Field: "b",
							Type: &StructLikeType{
								Name: "A",
							},
							Index: 2,
							GoTag: []string{`yaml:"b"`},
						},
						{
							Field: "cccc",
							Type: &ArrayType{
								ChildType: StringVal,
							},
							Index: 3,
							GoTag: []string{`yaml:"cccc"`},
						},
						{
							Field: "d",
							Type: &ArrayType{
								ChildType: &StructLikeType{
									Name: "D",
								},
							},
							Index: 4,
							GoTag: []string{`yaml:"d"`},
						},
						{
							Field: "e",
							Type: &StructLikeType{
								Name: "E",
							},
							Index: 5,
							GoTag: []string{`yaml:"e"`},
						},
						{
							Field: "f",
							Type: &StructLikeType{
								Name: "F",
							},
							Index: 6,
							GoTag: []string{`yaml:"f"`},
						},
						{
							Field: "ggg",
							Type: &ArrayType{
								ChildType: &ArrayType{
									ChildType: &StructLikeType{
										Name: "Ggg",
									},
								},
							},
							Index: 7,
							GoTag: []string{`yaml:"ggg"`},
						},
						{
							Field: "hh",
							Type: &ArrayType{
								ChildType: AnyVal,
							},
							Index: 8,
							GoTag: []string{`yaml:"hh"`},
						},
						{
							Field: "jj",
							Type: &ArrayType{
								ChildType: &ArrayType{
									ChildType: AnyVal,
								},
							},
							Index: 9,
							GoTag: []string{`yaml:"jj"`},
						},
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
				got1yaml, _ := yaml.Marshal(got1)
				want1yaml, _ := yaml.Marshal(tt.want1)
				t.Errorf("YamlParser.Parse got1 = %v, want1: %v", string(got1yaml), string(want1yaml))
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("YamlParser.Parse error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
