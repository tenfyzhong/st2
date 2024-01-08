package st2

import (
	"strings"
)

// Comment record a [Member] or [Struct] comments
type Comment struct {
	InlineComment string

	BeginningComments []string
}

// Member is fields of [Struct]
type Member struct {
	Field string
	Type
	Index    int
	Optional bool
	Comment  Comment
	GoTag    []string
}

// FieldCamel get a camel type field name
func (m Member) FieldCamel() string {
	return camel(m.Field)
}

// Go get the golang file type string
func (m Member) Go() string {
	name := m.Type.Go()
	if m.Optional && m.Type.IsBasicType() {
		return "*" + name
	}
	return name
}

// GoTagString get the go field tag string
func (m Member) GoTagString() string {
	if len(m.GoTag) == 0 {
		return ""
	}
	return "`" + strings.Join(m.GoTag, " ") + "`"
}

// Struct is a parsed result which contains the source struct data
type Struct struct {
	Type    Type
	Members []*Member
	Comment Comment
}
