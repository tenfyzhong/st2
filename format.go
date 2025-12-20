package st2

import (
	"fmt"
	"go/format"
)

// Format is an interface to format source code
type Format interface {
	Format(data []byte) ([]byte, error)
}

// EmptyFormat is a struct implement the [Format] interface with empty action
type EmptyFormater struct{}

func (f EmptyFormater) Format(data []byte) ([]byte, error) {
	return data, nil
}

// GoFormater is a struct implement the [Format] interface with format golang
// source data
type GoFormater struct{}

func (f GoFormater) Format(data []byte) ([]byte, error) {
	formatted, err := format.Source(data)
	if err != nil {
		return data, fmt.Errorf("failed to format Go code: %w", err)
	}
	return formatted, nil
}
