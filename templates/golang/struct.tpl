{{define "struct" }}
// {{.name}} {{.comment}}
type {{.name}} struct {
    {{ range $key, $field := .fields }}
    {{with $field.comment}}// {{$field.comment}}{{end}}
    {{$field.name}} {{$field.type}} {{template "structtags" $field.tag}}  
    {{ end }}
}
{{end}}

{{define "structtags" }}{{with .}}`{{.}}`{{end}}{{end}}