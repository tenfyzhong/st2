package st2

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

type UnmarshalFunc func(data []byte, v any) error

type UnmarshalTagFormat interface {
	Unmarshal(data []byte, v any) error
	TagFormat() string
}

// StructuredParser is a Parser to parse structured data source
type StructuredParser struct {
	ctx Context

	fingerMap map[string]*Struct
	nameMap   map[string]bool
	structs   []*Struct

	unmarshalTagFormat UnmarshalTagFormat
}

// NewStructuredParser create [StructuredParser]
func NewStructuredParser(ctx Context, unmarshalTagFormat UnmarshalTagFormat) *StructuredParser {
	return &StructuredParser{
		fingerMap:          make(map[string]*Struct),
		nameMap:            make(map[string]bool),
		unmarshalTagFormat: unmarshalTagFormat,
		ctx:                ctx,
	}
}

// Parse method parse structured data source
func (p *StructuredParser) Parse(reader io.Reader) ([]*Struct, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	var v any
	err = p.unmarshalTagFormat.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	rootName := p.ctx.Root
	if rootName == "" {
		rootName = RootDefault
	}

	root := p.parseNode(rootName, v)
	p.parseStructs(root)

	structs := p.structs
	p.fingerMap = make(map[string]*Struct)
	p.nameMap = make(map[string]bool)
	p.structs = p.structs[0:0]

	return structs, nil
}

func (p *StructuredParser) genUniqName(seed string) string {
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

func (p *StructuredParser) parseStructs(root *rawNode) *Member {
	if root == nil {
		return nil
	}

	member := &Member{
		Field: root.Field,
		GoTag: []string{fmt.Sprintf(p.unmarshalTagFormat.TagFormat(), root.Field)},
	}

	switch root.Type {
	case AnyVal,
		BoolVal,
		Float64Val,
		Float32Val,
		StringVal,
		Int8Val,
		Int16Val,
		Int32Val,
		Int64Val,
		Uint8Val,
		Uint16Val,
		Uint32Val,
		Uint64Val:
		member.Type = root.Type
	case ArrayVal:
		if len(root.Children) == 0 {
			// ignore the current member if the array is empty
			// the type of element is unknown
			return nil
		}
		root.Children[0].Field = root.Field
		child := p.parseStructs(root.Children[0])
		if child == nil {
			member.Type = &ArrayType{
				ChildType: AnyVal,
			}
		} else {
			member.Type = &ArrayType{
				ChildType: child.Type,
			}
		}
	case StructLikeVal:
		finger := root.Fingerprint()
		if st, ok := p.fingerMap[finger]; ok {
			member.Field = root.Field
			t, ok := st.Type.(*StructLikeType)
			if !ok {
				return nil
			}
			member.Type = &StructLikeType{
				Name: t.Name,
			}
			return member
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
			Type: &StructLikeType{
				Name:   name,
				Source: SLSStruct,
			},
		}

		p.structs = append(p.structs, st)
		p.fingerMap[finger] = st

		member.Type = &StructLikeType{
			Name: name,
		}
	default:
		return nil
	}

	return member
}

func (p *StructuredParser) parseNode(tag string, v any) *rawNode {
	node := &rawNode{
		Field: tag,
	}
	if v == nil {
		node.Type = AnyVal
	}

	switch c := v.(type) {
	case bool:
		node.Type = BoolVal
	case json.Number:
		_, err := c.Int64()
		if err == nil {
			node.Type = Int64Val
		} else {
			node.Type = Float64Val
		}
	case float32:
		node.Type = Float32Val
	case float64:
		node.Type = Float64Val
	case int8:
		node.Type = Int8Val
	case int16:
		node.Type = Int16Val
	case int32:
		node.Type = Int32Val
	case int64:
		node.Type = Int64Val
	case int:
		node.Type = Int64Val
	case uint8:
		node.Type = Uint8Val
	case uint16:
		node.Type = Uint16Val
	case uint32:
		node.Type = Uint32Val
	case uint64:
		node.Type = Uint64Val
	case uint:
		node.Type = Uint64Val
	case string:
		node.Type = StringVal
	case map[string]any:
		node.Type = StructLikeVal
		node.Children = []*rawNode{}
		for k, v := range c {
			child := p.parseNode(k, v)
			node.Children = append(node.Children, child)
		}
		sort.Sort(NodeList(node.Children))
	case []any:
		node.Type = ArrayVal
		if len(c) > 0 {
			child := p.parseNode("", c[0])
			node.Children = append(node.Children, child)
		} else {
			child := &rawNode{
				Type: AnyVal,
			}
			node.Children = append(node.Children, child)
		}
	}
	return node
}
