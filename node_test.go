package st2

import (
	"reflect"
	"testing"
)

func TestNode_Fingerprint(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *Node
		inspect func(r *Node, t *testing.T) //inspects receiver after test run

		want1 string
	}{
		{
			name: "null",
			init: func(t *testing.T) *Node {
				return &Node{
					Field: "root",
					Type:  NullVal,
				}
			},
			want1: "null",
		},
		{
			name: "number",
			init: func(t *testing.T) *Node {
				return &Node{
					Field: "number1",
					Type:  NumberVal,
				}
			},
			want1: "number",
		},
		{
			name: "empty array",
			init: func(t *testing.T) *Node {
				return &Node{
					Field: "arr",
					Type:  ArrayVal,
				}
			},
			want1: "[null]",
		},
		{
			name: "string array",
			init: func(t *testing.T) *Node {
				return &Node{
					Field: "arr",
					Type:  ArrayVal,
					Children: []*Node{
						{
							Field: "",
							Type:  StringVal,
						},
					},
				}
			},
			want1: "[:string]",
		},
		{
			name: "struct array",
			init: func(t *testing.T) *Node {
				return &Node{
					Field: "arr",
					Type:  ArrayVal,
					Children: []*Node{
						{
							Field: "second",
							Type:  StructVal,
							Children: []*Node{
								{
									Field: "third",
									Type:  NumberVal,
								},
							},
						},
					},
				}
			},
			want1: "[second:{third:number}]",
		},
		{
			name: "struct",
			init: func(t *testing.T) *Node {
				return &Node{
					Field: "st",
					Type:  StructVal,
					Children: []*Node{
						{
							Field: "a1",
							Type:  StructVal,
							Children: []*Node{
								{
									Field: "b1",
									Type:  NumberVal,
								},
							},
						},
						{
							Field: "a2",
							Type:  StructVal,
							Children: []*Node{
								{
									Field: "b2",
									Type:  StringVal,
								},
							},
						},
						{
							Field: "a3",
							Type:  StringVal,
						},
					},
				}
			},
			want1: "{a1:{b1:number};a2:{b2:string};a3:string}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			got1 := receiver.Fingerprint()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Node.Fingerprint got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

