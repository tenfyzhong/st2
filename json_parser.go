package st2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
)

type JsonParser struct {
	fingerMap map[string]*Struct
	nameMap   map[string]bool

	structs []*Struct
}

func NewJsonParser() *JsonParser {
	return &JsonParser{
		fingerMap: make(map[string]*Struct),
		nameMap:   make(map[string]bool),
	}
}

func (p *JsonParser) Parse(reader io.Reader) ([]*Struct, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("read data failed")
	}

	if len(data) == 0 {
		return nil, nil
	}

	var v interface{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	root := p.parseNode("root", v)
	p.parseStructs(root)

	structs := p.structs
	p.fingerMap = make(map[string]*Struct)
	p.nameMap = make(map[string]bool)
	p.structs = p.structs[0:0]

	return structs, nil
}

func (p *JsonParser) genUniqName(seed string) string {
	if !p.nameMap[seed] {
		return seed
	}

	for i := 1; i < 1000; i++ {
		name := fmt.Sprintf("%s%02d", seed, i)
		if !p.nameMap[name] {
			return name
		}
	}

	return p.genUniqName(seed + "a")
}

func (p *JsonParser) parseStructs(root *Node) *Member {
	if root == nil {
		return nil
	}

	switch root.Type {
	case NullVal, BoolVal, Float64Val, StringVal:
		return &Member{
			Field: root.Field,
			Type:  root.Type,
		}
	case ArrayVal:
		if len(root.Children) == 0 {
			// ignore the current memeber if the array is empty
			// the type of element is unknown
			return nil
		}
		root.Children[0].Field = root.Field
		child := p.parseStructs(root.Children[0])
		member := &Member{
			Field: root.Field,
			Type: &ArrayType{
				ChildType: NullVal,
			},
		}
		if child != nil {
			member.Type = &ArrayType{
				ChildType: child.Type,
			}
		}
		return member
	case StructVal:
		finger := root.Fingerprint()
		if st, ok := p.fingerMap[finger]; ok {
			return &Member{
				Field: root.Field,
				Type:  st.Type,
			}
		}

		name := p.genUniqName(root.FieldCamel())
		p.nameMap[name] = true

		members := make([]*Member, 0, len(root.Children))

		for i, child := range root.Children {
			childMember := p.parseStructs(child)
			if childMember != nil {
				childMember.Index = i + 1
				members = append(members, childMember)
			}
		}

		st := &Struct{
			Members: members,
			Type: &StructType{
				Name: name,
			},
		}

		p.structs = append(p.structs, st)
		p.fingerMap[finger] = st

		return &Member{
			Field: root.Field,
			Type:  st.Type,
		}
	}

	return nil
}

func (p *JsonParser) parseNode(tag string, v interface{}) *Node {
	node := &Node{
		Field: tag,
	}
	switch c := v.(type) {
	case bool:
		node.Type = BoolVal
	case float64:
		node.Type = Float64Val
	case string:
		node.Type = StringVal
	case map[string]interface{}:
		node.Type = StructVal
		node.Children = []*Node{}
		for k, v := range c {
			child := p.parseNode(k, v)
			node.Children = append(node.Children, child)
		}
		sort.Sort(NodeList(node.Children))
	case []interface{}:
		node.Type = ArrayVal
		if len(c) > 0 {
			child := p.parseNode("", c[0])
			node.Children = append(node.Children, child)
		}
	}
	return node
}
