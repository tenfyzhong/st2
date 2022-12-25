package st2

import (
	"errors"
	"io"
	"text/template"
)

func Convert(ctx Context, reader io.Reader, writer io.Writer) error {
	parse := CreateParser(ctx)
	if parse == nil {
		return errors.New("Can not found parser")
	}

	tmpl := CreateTmpl(ctx)
	if tmpl == "" {
		return errors.New("Can not found template")
	}

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
