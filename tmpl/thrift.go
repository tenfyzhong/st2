package tmpl

const Thrift = `
{{- define "MEMBER" }}
    {{- range $comment := .Comment.BeginningComments }}
    {{ $comment }}
    {{- end}}
    {{.Index}}: {{.Thrift}} {{.Field}}, {{ .Comment.InlineComment }} {{- end -}}

{{- define "STRUCT" -}}
{{- range $comment := .Comment.BeginningComments -}}
{{- $comment }}
{{ end -}}
{{- .Type.ThriftStructType }} {{ .Type.StructName }} { {{- if .Comment.InlineComment }} {{ .Comment.InlineComment }} {{- end }}
{{- range $member := .Members }}
{{- template "MEMBER" $member }}
{{- end }}
}
{{- end }}

{{- define "ENUM" -}}
{{- range $comment := .Comment.BeginningComments -}}
{{- $comment }}
{{ end -}}
enum {{ .Type.StructName }} { {{- if .Comment.InlineComment }} {{ .Comment.InlineComment }} {{- end }}
{{- range $member := .Members }} 
    {{ $member.Field }} = {{ $member.Index }}; {{ $member.Comment.InlineComment}} {{- end}}
}
{{- end }}

{{- range $st := . }}
{{- if eq $st.Type.ThriftStructType "enum" }}
{{- template "ENUM" $st -}}
{{- else -}}
{{- template "STRUCT" $st }}
{{- end }}

{{ end }}`
