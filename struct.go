package st2

import (
	"fmt"
	"strings"
)

type Comment struct {
	BeginningComments []string
	InlineComment     string
}

type Member struct {
	Field string
	Type
	Index    int
	Optional bool
	Comment  Comment
	GoTag    []string
}

func (m Member) FieldCamel() string {
	return Camel(m.Field)
}

func (m Member) Go() string {
	name := m.Type.Go()
	if m.Optional && m.Type.IsBasicType() {
		return "*" + name
	}
	return name
}

func (m Member) GoTagString() string {
	if len(m.GoTag) == 0 {
		return ""
	}

	strs := make([]string, 0)
	for _, tag := range m.GoTag {
		str := fmt.Sprintf("%s:\"%s\"", tag, m.Field)
		strs = append(strs, str)
	}
	return "`" + strings.Join(strs, " ") + "`"
}

type Struct struct {
	Members []*Member
	Type    Type
	Comment Comment
}
