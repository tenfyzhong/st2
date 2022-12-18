package st2

import (
	"io"
	"io/ioutil"
	"text/template"
)

func Convert(reader io.Reader, parse Parse, tmpl string, writer io.Writer) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	structs, err := parse.Parse(data)
	if err != nil {
		return err
	}

	t, err := template.New("st2").Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(writer, structs)
}
