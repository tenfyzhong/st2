package st2

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// CsvParser is a Parser to parse `.csv` source.
type CsvParser struct {
	ctx Context
}

// NewCsvParser create [CsvParser]
func NewCsvParser(ctx Context) *CsvParser {
	return &CsvParser{
		ctx: ctx,
	}
}

// Parse method parse csv source
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
	rootName = camel(rootName)

	st := &Struct{
		Type: &StructLikeType{
			Name:   rootName,
			Source: SLSStruct,
		},
	}
	for i, item := range items {
		member := &Member{
			Field: p.formatItem(item),
			Type:  StringVal,
			Index: i + 1,
		}
		member.GoTag = append(member.GoTag, fmt.Sprintf(`csv:"%s"`, member.Field))
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
