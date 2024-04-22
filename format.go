package st2

import "go/format"

// Format is an interface to format source code
type Format interface {
	Format(data []byte) []byte
}

// EmptyFormat is a struct implement the [Format] interface with empty action
type EmptyFormater struct {
}

func (f EmptyFormater) Format(data []byte) []byte {
	return data
}

// GoFormater is a struct implement the [Format] interface with format golang
// source data
type GoFormater struct {
}

func (f GoFormater) Format(data []byte) []byte {
	// do not return err if format failed
	// fallback to the origin data
	data, _ = format.Source(data)
	return data
}
