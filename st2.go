package st2

import (
	"io"
	"text/template"
)

func Convert(reader io.Reader, parse Parse, tmpl string, writer io.Writer) error {
	structs, err := parse.Parse(reader)
	if err != nil {
		return err
	}

	t, err := template.New("st2").Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(writer, structs)
}
