package st2

import (
	"reflect"
	"testing"
)

func TestJsonNode_Fingerprint(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *jsonNode
		inspect func(r *jsonNode, t *testing.T) //inspects receiver after test run

		want1 string
	}{
		{
			name: "null",
			init: func(t *testing.T) *jsonNode {
				return &jsonNode{
					Field: "root",
					Type:  AnyVal,
				}
			},
			want1: "null",
		},
		{
			name: "number",
			init: func(t *testing.T) *jsonNode {
				return &jsonNode{
					Field: "number1",
					Type:  Float64Val,
				}
			},
			want1: "number",
		},
		{
			name: "empty array",
			init: func(t *testing.T) *jsonNode {
				return &jsonNode{
					Field: "arr",
					Type:  ArrayVal,
				}
			},
			want1: "[null]",
		},
		{
			name: "string array",
			init: func(t *testing.T) *jsonNode {
				return &jsonNode{
					Field: "arr",
					Type:  ArrayVal,
					Children: []*jsonNode{
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
			init: func(t *testing.T) *jsonNode {
				return &jsonNode{
					Field: "arr",
					Type:  ArrayVal,
					Children: []*jsonNode{
						{
							Field: "second",
							Type:  StructLikeVal,
							Children: []*jsonNode{
								{
									Field: "third",
									Type:  Float64Val,
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
			init: func(t *testing.T) *jsonNode {
				return &jsonNode{
					Field: "st",
					Type:  StructLikeVal,
					Children: []*jsonNode{
						{
							Field: "a1",
							Type:  StructLikeVal,
							Children: []*jsonNode{
								{
									Field: "b1",
									Type:  Float64Val,
								},
							},
						},
						{
							Field: "a2",
							Type:  StructLikeVal,
							Children: []*jsonNode{
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
