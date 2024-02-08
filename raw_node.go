package st2

import (
	"strings"
)

type rawNode struct {
	Field       string
	Type        Type
	Children    []*rawNode
	fingerprint string
}

func (node rawNode) FieldCamel() string {
	return camel(node.Field)
}

func (node *rawNode) Fingerprint() string {
	if node.fingerprint == "" {
		node.fingerprint = node.fingerprintHelper()
	}
	return node.fingerprint
}

func (node *rawNode) fingerprintHelper() string {
	switch node.Type {
	case ArrayVal:
		child := "null"
		if len(node.Children) > 0 {
			child = node.Children[0].Field + ":" + node.Children[0].Fingerprint()
		}
		return "[" + child + "]"
	case StructLikeVal:
		children := make([]string, 0)
		for _, child := range node.Children {
			finger := child.Fingerprint()
			children = append(children, child.Field+":"+finger)
		}
		return "{" + strings.Join(children, ";") + "}"
	default:
		return node.Type.Json()
	}
}

type NodeList []*rawNode

func (l NodeList) Len() int {
	return len(l)
}

func (l NodeList) Less(i, j int) bool {
	return l[i].Field < l[j].Field
}

func (l NodeList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
