package st2

import (
	"bytes"
	"errors"
	"io"
	"text/template"
)

// Convert is a wrap function parse from reader and write the output to writer
func Convert(ctx Context, reader io.Reader, writer io.Writer) error {
	if reader == nil {
		return errors.New("reader is nil")
	}

	if writer == nil {
		return errors.New("writer is nil")
	}

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

	b := new(bytes.Buffer)
	err = t.Execute(b, structs)
	if err != nil {
		return err
	}

	data := b.Bytes()
	formater := CreateFormater(ctx)
	data = formater.Format(data)

	_, err = writer.Write(data)
	return err
}
