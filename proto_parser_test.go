package st2

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtoParser_Parse(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) ProtoParser
		inspect func(r ProtoParser, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      []*Struct
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "err",
			init: func(t *testing.T) ProtoParser {
				return *NewProtoParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte("a")),
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.EqualError(t, err, "found \"\\\"a\\\"(Token=2, Pos=<input>:1:1)\" but expected [syntax] at /Users/bytedance/go/pkg/mod/github.com/yoheimuta/go-protoparser/v4@v4.7.0/parser/syntax.go:62")
			},
		},
		{
			name: "succ",
			init: func(t *testing.T) ProtoParser {
				return *NewProtoParser(Context{})
			},
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte(`
syntax = "proto3";

import "google/protobuf/any.proto";
option go_package="tenfyzhong/st2";

// EEEE
enum Eeee { // EEEE
    // A
    A = 0; // a
}

// haha
message Aaa { // aaa
    // a
    int32 a = 1; // a
    int64 b = 2;
    string c = 3;
}

message BbbBB {
    uint32 a = 1;
    uint64 b = 2;
    string c = 3;
}

message Ccc {
    sint32 a = 1;
    sint64 b = 2;
    string c = 3;
    Aaa aaa = 4;
}

message ErrorStatus {
  string message = 1;
  repeated google.protobuf.Any details = 2;
  map<int32, Ccc> m = 3;
}
`)),
				}
			},
			want1: []*Struct{
				{
					Type: &EnumType{
						Name: "Eeee",
					},
					Members: []*Member{
						{
							Field: "A",
							Type: &EnumType{
								Name: "Eeee",
							},
							Index: 0,
							Comment: Comment{
								BeginningComments: []string{"// A"},
								InlineComment:     "// a",
							},
						},
					},
					Comment: Comment{
						BeginningComments: []string{"// EEEE"},
						InlineComment:     "// EEEE",
					},
				},
				{
					Type: &StructLikeType{
						Name:   "Aaa",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "a",
							Type:  Int32Val,
							Index: 1,
							Comment: Comment{
								BeginningComments: []string{"// a"},
								InlineComment:     "// a",
							},
						},
						{
							Field: "b",
							Type:  Int64Val,
							Index: 2,
						},
						{
							Field: "c",
							Type:  StringVal,
							Index: 3,
						},
					},
					Comment: Comment{
						BeginningComments: []string{"// haha"},
						InlineComment:     "// aaa",
					},
				},
				{
					Type: &StructLikeType{
						Name:   "BbbBB",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "a",
							Type:  Uint32Val,
							Index: 1,
						},
						{
							Field: "b",
							Type:  Uint64Val,
							Index: 2,
						},
						{
							Field: "c",
							Type:  StringVal,
							Index: 3,
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "Ccc",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "a",
							Type:  Int32Val,
							Index: 1,
						},
						{
							Field: "b",
							Type:  Int64Val,
							Index: 2,
						},
						{
							Field: "c",
							Type:  StringVal,
							Index: 3,
						},
						{
							Field: "aaa",
							Type: &StructLikeType{
								Name: "Aaa",
							},
							Index: 4,
						},
					},
				},
				{
					Type: &StructLikeType{
						Name:   "ErrorStatus",
						Source: SLSStruct,
					},
					Members: []*Member{
						{
							Field: "message",
							Type:  StringVal,
							Index: 1,
						},
						{
							Field: "details",
							Type: &ArrayType{
								ChildType: &StructLikeType{
									Name: "google.protobuf.Any",
								},
							},
							Index: 2,
						},
						{
							Field: "m",
							Type: &MapType{
								Key: Int32Val,
								Value: &StructLikeType{
									Name: "Ccc",
								},
							},
							Index: 3,
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
				t.Errorf("ProtoParser.Parse got1 = %v, want1: %v", string(got1Json), string(want1Json))
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("ProtoParser.Parse error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
