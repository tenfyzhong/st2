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
	B  int64   `json:"b,omitempty"`
	C  *string `json:"C,omitempty" proto:"C"`
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
