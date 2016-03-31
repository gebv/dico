{{define "server" }}
server {
    {{with $.listen}}listen {{$.listen}};{{end}}
    {{ range $key, $location := $.locations}}
        {{template "location" $location}}
    {{end}}
}
{{ end}}