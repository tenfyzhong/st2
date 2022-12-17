package st2

import (
	"strings"
)

type Node struct {
	Field       string
	Type        Type
	Children    []*Node
	fingerprint string
}

func (node Node) FieldCamel() string {
	return Camel(node.Field)
}

func (node *Node) Fingerprint() string {
	if node.fingerprint == "" {
		node.fingerprint = node.fingerprintHelper()
	}
	return node.fingerprint
}

func (node *Node) fingerprintHelper() string {
	switch node.Type {
	case ArrayVal:
		child := "null"
		if len(node.Children) > 0 {
			child = node.Children[0].Field + ":" + node.Children[0].Fingerprint()
		}
		return "[" + child + "]"
	case StructVal:
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

type NodeList []*Node

func (l NodeList) Len() int {
	return len(l)
}

func (l NodeList) Less(i, j int) bool {
	return l[i].Field < l[j].Field
}

func (l NodeList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
