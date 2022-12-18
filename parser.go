package st2

import "io"

type Parse interface {
	Parse(r io.Reader) ([]*Struct, error)
}
