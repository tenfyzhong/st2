package st2

import (
	"bufio"
	"io"
	"strings"
)

type CsvParser struct {
	ctx Context
}

func NewCsvParser(ctx Context) *CsvParser {
	return &CsvParser{
		ctx: ctx,
	}
}

func (p CsvParser) Parse(reader io.Reader) ([]*Struct, error) {
	bufioReader := bufio.NewReader(reader)
	str, err := bufioReader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	items := strings.Split(str, ",")
	rootName := p.ctx.Root
	if rootName == "" {
		rootName = RootDefault
	}
	rootName = Camel(rootName)

	st := &Struct{
		Type: &StructType{
			Name: rootName,
			Type: "struct",
		},
	}
	for i, item := range items {
		member := &Member{
			Field: p.formatItem(item),
			Type:  StringVal,
			Index: i + 1,
		}
		st.Members = append(st.Members, member)
	}
	return []*Struct{st}, nil
}

func (p CsvParser) formatItem(str string) string {
	str = strings.TrimSpace(str)
	items := strings.Split(str, " ")
	last := 0
	for i := 0; i < len(items); i++ {
		if items[i] != "" {
			items[last] = items[i]
			last++
		}
	}
	items = items[0:last]
	return strings.Join(items, "_")
}
