package st2

import "io"

// Parse interface has a method `Parse` which parse source code to a list
// of [Struct]
type Parse interface {
	// Parse parse source code to Struct
	Parse(r io.Reader) ([]*Struct, error)
}
