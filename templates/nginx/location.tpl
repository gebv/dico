{{define "location" }}
location {{$.path}} {
    {{with $.proxy_pass}}proxy_pass {{$.proxy_pass}};{{end}}
    {{with $.root}}root {{$.root}};{{end}}
}
{{end}}