package st2

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCsvParser_Parse(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) CsvParser
		inspect func(r CsvParser, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      []*Struct
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "empty",
			init: func(t *testing.T) CsvParser {
				return *NewCsvParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte("")),
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.EqualError(t, err, "EOF")
			},
		},
		{
			name: "one item",
			init: func(t *testing.T) CsvParser {
				return *NewCsvParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte("a\n")),
				}
			},
			want1: []*Struct{
				{
					Type: &StructLikeType{
						Name:   "Root",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "a",
							Type:  StringVal,
							Index: 1,
							GoTag: []string{`csv:"a"`},
						},
					},
				},
			},
		},
		{
			name: "succ",
			init: func(t *testing.T) CsvParser {
				return *NewCsvParser(Context{
					Root: "hello",
				})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte(`a,b,c,d,c d , a b ,hello world, t       t   
1,2,3,4
`)),
				}
			},
			want1: []*Struct{
				{
					Type: &StructLikeType{
						Name:   "Hello",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "a",
							Type:  StringVal,
							Index: 1,
							GoTag: []string{`csv:"a"`},
						},
						{
							Field: "b",
							Type:  StringVal,
							Index: 2,
							GoTag: []string{`csv:"b"`},
						},
						{
							Field: "c",
							Type:  StringVal,
							Index: 3,
							GoTag: []string{`csv:"c"`},
						},
						{
							Field: "d",
							Type:  StringVal,
							Index: 4,
							GoTag: []string{`csv:"d"`},
						},
						{
							Field: "c_d",
							Type:  StringVal,
							Index: 5,
							GoTag: []string{`csv:"c_d"`},
						},
						{
							Field: "a_b",
							Type:  StringVal,
							Index: 6,
							GoTag: []string{`csv:"a_b"`},
						},
						{
							Field: "hello_world",
							Type:  StringVal,
							Index: 7,
							GoTag: []string{`csv:"hello_world"`},
						},
						{
							Field: "t_t",
							Type:  StringVal,
							Index: 8,
							GoTag: []string{`csv:"t_t"`},
						},
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
				t.Errorf("CsvParser.Parse got1 = %v, want1: %v", string(got1json), string(want1json))
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CsvParser.Parse error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
