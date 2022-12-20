package st2

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

type Struct struct {
	Members []*Member
	Type    Type
	Comment Comment
}
