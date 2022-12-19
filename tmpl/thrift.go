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

{{- define "ENUM" -}}
enum {{ .Type.StructName }} {
{{- range $member := .Members }} 
	{{ $member.Field }} = {{ $member.Index }}; {{- end}}
}
{{- end }}

{{- range $st := . -}}
{{- if eq $st.Type.ThriftStructType "enum" }}
{{- template "ENUM" $st }}
{{- else }}
{{- template "STRUCT" $st }}
{{- end }}

{{ end }}`
