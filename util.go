package st2

import (
	"strings"

	"github.com/iancoleman/strcase"
)

var acronyms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

func Camel(s string) string {
	items := strings.Split(s, "_")
	for i, item := range items {
		if len(item) == 0 {
			continue
		}

		upper := strings.ToUpper(item)
		if acronyms[upper] {
			items[i] = upper
		} else {
			items[i] = strings.ToUpper(item[0:1]) + strings.ToLower(item[1:])
		}
	}
	return strings.Join(items, "")
}

func Snake(s string) string {
	return strcase.ToSnake(s)
}

func init() {
	for str := range acronyms {
		strcase.ConfigureAcronym(str, strings.ToLower(str))
	}
}
