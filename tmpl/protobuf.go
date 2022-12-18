package tmpl

const Protobuf = `
{{- define "MEMBER" }}
    {{.Protobuf}} {{.Field}} = {{.Index}};
{{- end -}}

{{- define "STRUCT" -}}
message {{ .Type.StructName }} {
{{- range $member := .Members }}
{{- template "MEMBER" $member }}
{{- end }}
}

{{- end }}
{{- range $st := . -}}
{{- template "STRUCT" $st }}

{{ end -}}
`
