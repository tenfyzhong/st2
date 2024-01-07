package st2

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	type args struct {
		ctx    Context
		buffer *bytes.Buffer
		reader io.Reader
		writer io.Writer
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		wantData   []byte
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "reader is nil",
			args: func(t *testing.T) args {
				return args{
					reader: nil,
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.EqualError(t, err, "reader is nil")
			},
		},
		{
			name: "writer is nil",
			args: func(t *testing.T) args {
				return args{
					reader: bytes.NewReader([]byte{}),
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.EqualError(t, err, "writer is nil")
			},
		},
		{
			name: "parser not found",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "a",
					},
					reader: bytes.NewReader([]byte{}),
					buffer: bytes.NewBuffer(nil),
				}
				a.writer = a.buffer
				return a
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.EqualError(t, err, "Can not found parser")
			},
		},
		{
			name: "tmpl not found",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "json",
						Dst: "bb",
					},
					reader: bytes.NewReader([]byte{}),
					buffer: bytes.NewBuffer(nil),
				}
				a.writer = a.buffer
				return a
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.EqualError(t, err, "Can not found template")
			},
		},
		{
			name: "Parser parse failed",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "json",
						Dst: "go",
					},
					reader: bytes.NewReader([]byte(`a`)),
					buffer: bytes.NewBuffer(nil),
				}
				a.writer = a.buffer
				return a
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.EqualError(t, err, `Read: unexpected value type: 0, error found in #0 byte of ...|a|..., bigger context ...|a|...`)
			},
		},
		{
			name: "json to go",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "json",
						Dst: "go",
					},
					buffer: bytes.NewBuffer(nil),
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
    "c": ["123"],
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
    "gg": [],
	"h": null
}]
`)),
				}
				a.writer = a.buffer
				return a
			},
			wantData: []byte(`type A struct {
	B int64 ` + "`json:\"b,omitempty\"`" + `
	C string ` + "`json:\"c,omitempty\"`" + `
}

type D struct {
	B int64 ` + "`json:\"b,omitempty\"`" + `
	C int64 ` + "`json:\"c,omitempty\"`" + `
}

type E struct {
	Aa bool ` + "`json:\"aa,omitempty\"`" + `
	Bb bool ` + "`json:\"bb,omitempty\"`" + `
}

type A01 struct {
	Hello bool ` + "`json:\"hello,omitempty\"`" + `
}

type F struct {
	A *A01 ` + "`json:\"a,omitempty\"`" + `
}

type Root struct {
	A *A ` + "`json:\"a,omitempty\"`" + `
	B *A ` + "`json:\"b,omitempty\"`" + `
	C []string ` + "`json:\"c,omitempty\"`" + `
	D []*D ` + "`json:\"d,omitempty\"`" + `
	E *E ` + "`json:\"e,omitempty\"`" + `
	F *F ` + "`json:\"f,omitempty\"`" + `
	H any ` + "`json:\"h,omitempty\"`" + `
}

`),
			wantErr: false,
		},
		{
			name: "json to proto",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "json",
						Dst: "proto",
					},
					buffer: bytes.NewBuffer(nil),
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
    "c": ["123"],
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
    "gg": [],
	"h": null
}]
`)),
				}
				a.writer = a.buffer
				return a
			},
			wantData: []byte(`message A {
    int64 b = 1; 
    string c = 2; 
}

message D {
    int64 b = 1; 
    int64 c = 2; 
}

message E {
    bool aa = 1; 
    bool bb = 2; 
}

message A01 {
    bool hello = 1; 
}

message F {
    A01 a = 1; 
}

message Root {
    A a = 1; 
    A b = 2; 
    repeated string c = 3; 
    repeated D d = 4; 
    E e = 5; 
    F f = 6; 
    google.protobuf.Any h = 8; 
}

`),
			wantErr: false,
		},
		{
			name: "json to thrift",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "json",
						Dst: "thrift",
					},
					buffer: bytes.NewBuffer(nil),
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
    "c": ["123"],
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
    "gg": [],
	"h": null
}]
`)),
				}
				a.writer = a.buffer
				return a
			},
			wantData: []byte(`struct A {
    1: i64 b, 
    2: string c, 
}

struct D {
    1: i64 b, 
    2: i64 c, 
}

struct E {
    1: bool aa, 
    2: bool bb, 
}

struct A01 {
    1: bool hello, 
}

struct F {
    1: A01 a, 
}

struct Root {
    1: A a, 
    2: A b, 
    3: list<string> c, 
    4: list<D> d, 
    5: E e, 
    6: F f, 
    8: binary h, 
}

`),
			wantErr: false,
		},
		{
			name: "proto to go",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "proto",
						Dst: "go",
					},
					buffer: bytes.NewBuffer(nil),
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
    int32 a = 1;
    int64 b = 2;
    string c = 3;
}

message Ccc {
    int32 a = 1;
    int64 b = 2;
    string c = 3;
    Aaa aaa = 4;
}

message ErrorStatus {
  string message = 1;
  repeated google.protobuf.Any details = 2;
}

message SampleMessage {
  oneof test_oneof {
    string name = 4;
    ErrorStatus sub_message = 9;
  }
}
`)),
				}
				a.writer = a.buffer
				return a
			},
			wantData: []byte(`// EEEE
type Eeee int // EEEE 

const ( 
	A Eeee = 0 // a
)

// haha
type Aaa struct { // aaa
	// a
	A int32 // a
	B int64
	C string
}

type BbbBB struct {
	A int32
	B int64
	C string
}

type Ccc struct {
	A int32
	B int64
	C string
	Aaa *Aaa
}

type ErrorStatus struct {
	Message string
	Details []*protobuf.Any
}

type SampleMessage struct {
}

`),
			wantErr: false,
		},
		{
			name: "proto to thrift",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "proto",
						Dst: "thrift",
					},
					buffer: bytes.NewBuffer(nil),
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
    int32 a = 1;
    int64 b = 2;
    string c = 3;
}

message Ccc {
    int32 a = 1;
    int64 b = 2;
    string c = 3;
    Aaa aaa = 4;
}

message ErrorStatus {
  string message = 1;
  repeated google.protobuf.Any details = 2;
}

message SampleMessage {
  oneof test_oneof {
    string name = 4;
    ErrorStatus sub_message = 9;
  }
}
`)),
				}
				a.writer = a.buffer
				return a
			},
			wantData: []byte(`// EEEE
enum Eeee { // EEEE  
    A = 0; // a
}

// haha
struct Aaa { // aaa
    // a
    1: i32 a, // a
    2: i64 b, 
    3: string c, 
}

struct BbbBB {
    1: i32 a, 
    2: i64 b, 
    3: string c, 
}

struct Ccc {
    1: i32 a, 
    2: i64 b, 
    3: string c, 
    4: Aaa aaa, 
}

struct ErrorStatus {
    1: string message, 
    2: list<google.protobuf.Any> details, 
}

struct SampleMessage {
}

`),
			wantErr: false,
		},
		{
			name: "thrift to go",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "thrift",
						Dst: "go",
					},
					buffer: bytes.NewBuffer(nil),
					reader: bytes.NewReader([]byte(`
enum EEE {
    A = 1;
    B = 2;
}

// hhhh
struct SS { // aa
    // ss
    1: optional bool a, // jjj
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
					`)),
				}
				a.writer = a.buffer
				return a
			},
			wantData: []byte(`type EEE int

const ( 
	A EEE = 1  
	B EEE = 2 
)

type SS struct {
	A *bool
	B int8
	C int16
	D int32
	E int64
	F float64
	G string
	H []byte
	I map[int32]int32
	J []int32
	K map[int32]bool
}

type AAA struct {
	Hello string
}

type BBB struct {
	B1 int16
	B2 int32
	E *EEE
	Mapab map[*AAA]*BBB
	Seta map[*AAA]bool
	Listb []*BBB
}

type UUU struct {
	A *AAA
	B *BBB
}

`),
			wantErr: false,
		},
		{
			name: "thrift to proto",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "thrift",
						Dst: "proto",
					},
					buffer: bytes.NewBuffer(nil),
					reader: bytes.NewReader([]byte(`
enum EEE {
    A = 1;
    B = 2;
}

// hhhh
struct SS { // aa
    // ss
    1: optional bool a, // jjj
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
					`)),
				}
				a.writer = a.buffer
				return a
			},
			wantData: []byte(`enum EEE { 
    A = 1;  
    B = 2; 
}

message SS {
    bool a = 1; 
    int32 b = 2; 
    int32 c = 3; 
    int32 d = 4; 
    int64 e = 5; 
    double f = 6; 
    string g = 7; 
    bytes h = 8; 
    map<int32, int32> i = 9; 
    repeated int32 j = 10; 
    map<int32, bool> k = 11; 
}

message AAA {
    string hello = 1; 
}

message BBB {
    int32 b1 = 1; 
    int32 b2 = 2; 
    EEE e = 3; 
    map<AAA, BBB> mapab = 4; 
    map<AAA, bool> seta = 5; 
    repeated BBB listb = 6; 
}

message UUU {
    AAA a = 1; 
    BBB b = 2; 
}

`),
			wantErr: false,
		},
		{
			name: "go to proto",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "go",
						Dst: "proto",
					},
					buffer: bytes.NewBuffer(nil),
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
	// comment EEEA Eeee block
	// comment EEEA Eeee block
	EEEA Eeee = 0 // comment EEEA Eeee inline
	EEEB Eeee = 1 // a
	EEEC Eeee = 3 // a

	II1 int = 1
	II2 int = 2

	IN1 = 1
	IN2 = 2
)

// comment haha
// comment hehe
type Aaa struct { // comment aaa inline
	// comment Aaa a
	A  []int32 // comment Aaa a inline
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
				a.writer = a.buffer
				return a
			},
			wantData: []byte(`enum Eeee { 
    EEEA = 0; // comment EEEA Eeee inline 
    EEEB = 1; // a 
    EEEC = 3; // a
}

// comment haha
// comment hehe
message Aaa {
    // comment Aaa a
    repeated int32 a = 1; // comment Aaa a inline
    int64 b = 2; 
    string c = 3; 
    map<int64, string> mm = 4; 
}

message BbbBB {
    int32 a = 1; 
    int64 b = 2; 
    string c = 3; 
}

message Ccc {
    int32 a = 1; 
    int64 b = 2; 
    string c = 3; 
    Aaa aaa = 4; 
}

message ErrorStatus {
    string message = 1; 
    repeated protobuf.Any details = 2; 
}

message SampleMessage {
}

`),
			wantErr: false,
		},
		{
			name: "go to thrift",
			args: func(t *testing.T) args {
				a := args{
					ctx: Context{
						Src: "go",
						Dst: "thrift",
					},
					buffer: bytes.NewBuffer(nil),
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
	// comment EEEA Eeee block
	// comment EEEA Eeee block
	EEEA Eeee = 0 // comment EEEA Eeee inline
	EEEB Eeee = 1 // a
	EEEC Eeee = 3 // a

	II1 int = 1
	II2 int = 2

	IN1 = 1
	IN2 = 2
)

// comment haha
// comment hehe
type Aaa struct { // comment aaa inline
	// comment Aaa a
	A  []int32 // comment Aaa a inline
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
				a.writer = a.buffer
				return a
			},
			wantData: []byte(`enum Eeee { 
    EEEA = 0; // comment EEEA Eeee inline 
    EEEB = 1; // a 
    EEEC = 3; // a
}

// comment haha
// comment hehe
struct Aaa {
    // comment Aaa a
    1: list<i32> a, // comment Aaa a inline
    2: i64 b, 
    3: string c, 
    4: map<i64, string> mm, 
}

struct BbbBB {
    1: i32 a, 
    2: i64 b, 
    3: string c, 
}

struct Ccc {
    1: i32 a, 
    2: i64 b, 
    3: string c, 
    4: Aaa aaa, 
}

struct ErrorStatus {
    1: string message, 
    2: list<protobuf.Any> details, 
}

struct SampleMessage {
}

`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			err := Convert(tArgs.ctx, tArgs.reader, tArgs.writer)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Convert error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}

			if !tt.wantErr {
				actual := tArgs.buffer.Bytes()
				assert.Equal(t, tt.wantData, actual, fmt.Sprintf("[%s] should equal to [%s]", string(tt.wantData), string(actual)))
			}
		})
	}
}
