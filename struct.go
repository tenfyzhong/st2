package st2

type Member struct {
	Field string
	Type
	Index int
}

func (m Member) FieldCamel() string {
	return Camel(m.Field)
}

type Struct struct {
	Members []*Member
	Type    Type
}
