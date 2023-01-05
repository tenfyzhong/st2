package st2

import (
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
	return "`" + strings.Join(m.GoTag, " ") + "`"
}

type Struct struct {
	Type    Type
	Members []*Member
	Comment Comment
}
