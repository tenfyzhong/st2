package st2

import "strings"

func Camel(s string) string {
	items := strings.Split(s, "_")
	for i, item := range items {
		if len(item) == 0 {
			continue
		}
		items[i] = strings.ToUpper(item[0:1]) + strings.ToLower(item[1:])
	}
	return strings.Join(items, "")
}
