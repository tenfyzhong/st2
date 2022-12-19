package st2

type Member struct {
	Field string
	Type
	Index    int
	Optional bool
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
}
