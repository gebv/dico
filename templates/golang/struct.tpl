{{define "struct" }}
{{- with .comment }}// {{.name}} {{.comment}}{{end}}
type {{.name}} struct {
    {{ range $key, $field := .fields}}
    {{with $field.comment}}// {{$field.comment}}{{end}}
    {{$field.name}} {{$field.type}} {{template "structtags" $field.tag}}  
    {{ end }}
}
{{- $structName := .name}}
{{range $key, $field := .fields}}
{{- template "setter" (map "structname" $structName "field" $field) -}}
{{- template "getter" (map "structname" $structName "field" $field) -}}  
{{ end }}
{{ end }}

{{define "structtags" }}{{with .}}`{{.}}`{{end}}{{end}}