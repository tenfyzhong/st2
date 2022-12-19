package tmpl

const Go = `
{{- define "MEMBER" }}
	{{.FieldCamel}} {{.Go}}` + " `json:\"{{.Field}}\"`" + `{{- end -}}

{{- define "STRUCT" -}}
type {{ .Type.StructName }} struct {
{{- range $member := .Members }}    
{{- template "MEMBER" $member }}
{{- end }}
}
{{- end }}

{{- define "ENUM" -}}
type {{ .Type.StructName }} int
const (
{{- range $member := .Members }} 
	{{ $member.FieldCamel }} {{ $member.Go }} = {{ $member.Index }} {{- end}}
)
{{- end }}

{{- range $st := . -}}
{{- if eq $st.Type.GoStructType "enum" }}
{{- template "ENUM" $st }}
{{- else }}
{{- template "STRUCT" $st }}
{{- end }}

{{ end }}`
