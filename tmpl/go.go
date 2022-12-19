package tmpl

const Go = `
{{- define "MEMBER" }}
	{{- range $comment := .Comment.BeginningComments }}
	{{ $comment }}
	{{- end}}
	{{.FieldCamel}} {{.Go}}` + " `json:\"{{.Field}}\"`" + ` {{ .Comment.InlineComment }} {{- end -}}

{{- define "STRUCT" -}}
{{- range $comment := .Comment.BeginningComments -}}
{{- $comment }}
{{ end -}}
type {{ .Type.StructName }} struct { {{ .Comment.InlineComment }}
{{- range $member := .Members }}    
{{- template "MEMBER" $member }}
{{- end }}
}
{{- end }}

{{- define "ENUM" -}}
{{- range $comment := .Comment.BeginningComments -}}
{{- $comment }}
{{ end -}}
type {{ .Type.StructName }} int {{ .Comment.InlineComment }}

const (
{{- range $member := .Members }} 
	{{ $member.FieldCamel }} {{ $member.Go }} = {{ $member.Index }} {{ $member.Comment.InlineComment }} {{- end}}
)
{{- end }}

{{- range $st := . -}}
{{- if eq $st.Type.GoStructType "enum" }}
{{- template "ENUM" $st }}
{{- else }}
{{- template "STRUCT" $st }}
{{- end }}

{{ end }}`
