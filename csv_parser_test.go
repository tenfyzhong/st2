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
					Type: &StructType{
						Name: "Root",
						Type: "struct",
					},
					Members: []*Member{
						{
							Field: "a",
							Type:  StringVal,
							Index: 1,
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
					Type: &StructType{
						Name: "Hello",
						Type: "struct",
					},
					Members: []*Member{
						{
							Field: "a",
							Type:  StringVal,
							Index: 1,
						},
						{
							Field: "b",
							Type:  StringVal,
							Index: 2,
						},
						{
							Field: "c",
							Type:  StringVal,
							Index: 3,
						},
						{
							Field: "d",
							Type:  StringVal,
							Index: 4,
						},
						{
							Field: "c_d",
							Type:  StringVal,
							Index: 5,
						},
						{
							Field: "a_b",
							Type:  StringVal,
							Index: 6,
						},
						{
							Field: "hello_world",
							Type:  StringVal,
							Index: 7,
						},
						{
							Field: "t_t",
							Type:  StringVal,
							Index: 8,
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
