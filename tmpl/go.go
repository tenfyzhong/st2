package tmpl

const Go = `
{{- define "MEMBER" }}
	{{.FieldCamel}} {{.Go}}` + " `json:\"{{.Field}}\"`" + `{{- end -}}

{{- define "STRUCT" }}
type {{ .Type.GoStruct }} struct {
{{- range $member := .Members }}    
{{- template "MEMBER" $member }}
{{- end}}
}
{{- end }}

{{- range $st := . -}}
{{- template "STRUCT" $st }}
{{ end }}`
