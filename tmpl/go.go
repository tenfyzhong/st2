package tmpl

const Go = `
{{- define "MEMBER" }}
	{{- range $comment := .Comment.BeginningComments }}
	{{ $comment }}
	{{- end}}
	{{.FieldCamel}} {{.Go}}
	{{- if .GoTagString }} {{ .GoTagString }}{{- end }}
	{{- if .Comment.InlineComment }} {{ .Comment.InlineComment }} {{- end -}}
{{- end }}

{{- define "STRUCT" -}}
{{- range $comment := .Comment.BeginningComments -}}
{{- $comment }}
{{ end -}}
type {{ .Type.StructName }} struct { {{- if .Comment.InlineComment }} {{ .Comment.InlineComment }} {{- end }}
{{- range $member := .Members }}    
{{- template "MEMBER" $member }}
{{- end }}
}
{{- end }}

{{- define "ENUM" -}}
{{- range $comment := .Comment.BeginningComments -}}
{{- $comment }}
{{ end -}}
type {{ .Type.StructName }} int {{- if .Comment.InlineComment }} {{ .Comment.InlineComment }} {{- end }}

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
