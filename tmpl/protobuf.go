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

{{- define "ENUM" -}}
enum {{ .Type.StructName }} {
{{- range $member := .Members }} 
	{{ $member.Field }} = {{ $member.Index }} {{- end}};
}
{{- end }}

{{- range $st := . -}}
{{- if eq $st.Type.ProtobufStructType "enum" }}
{{- template "ENUM" $st }}
{{- else }}
{{- template "STRUCT" $st }}
{{- end }}

{{ end }}`
