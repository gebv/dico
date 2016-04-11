{{define "struct" }}
func New{{toUpper .name}}() *{{.name}} {
    model := new({{.name}})
    {{ range $key, $field := .fields}}
    {{if hasPrefix $field.type "map"}}
    model.{{$field.name}} = make({{$field.type}})
    {{end}}  
    {{ end }}
    return model
}

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