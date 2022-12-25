package tmpl

const Proto = `
{{- define "MEMBER" }}
    {{- range $comment := .Comment.BeginningComments }}
    {{ $comment }}
    {{- end}}
    {{.Proto}} {{.Field}} = {{.Index}}; {{ .Comment.InlineComment }} {{- end -}}

{{- define "STRUCT" -}}
{{- range $comment := .Comment.BeginningComments -}}
{{- $comment }}
{{ end -}}
message {{ .Type.StructName }} { {{- if .Comment.InlineComment }} {{ .Comment.InlineComment }} {{- end }}
{{- range $member := .Members }}
{{- template "MEMBER" $member }}
{{- end }}
}
{{- end }}

{{- define "ENUM" -}}
{{- range $comment := .Comment.BeginningComments -}}
{{ $comment }}
{{- end }}
enum {{ .Type.StructName }} { {{- if .Comment.InlineComment }} {{ .Comment.InlineComment }} {{- end }}
{{- range $member := .Members }} 
    {{ $member.Field }} = {{ $member.Index }}; {{ $member.Comment.InlineComment}} {{- end}}
}
{{- end }}

{{- range $st := . }}
{{- if eq $st.Type.ProtoStructType "enum" }}
{{- template "ENUM" $st -}}
{{- else -}}
{{- template "STRUCT" $st }}
{{- end }}

{{ end }}`
