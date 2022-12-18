package st2

type Parse interface {
	Parse(data []byte) ([]*Struct, error)
}
