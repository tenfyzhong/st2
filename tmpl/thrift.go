package tmpl

const Thrift = `
{{- define "MEMBER" }}
    {{.Index}}: {{.Thrift}} {{.Field}},
{{- end -}}

{{- define "STRUCT" -}}
struct {{ .Type.StructName }} {
{{- range $member := .Members }}
{{- template "MEMBER" $member }}
{{- end }}
}
{{- end }}

{{- range $st := . -}}
{{- template "STRUCT" $st }}

{{ end -}}
`
