{{define "helpfull" }}
// FromJson extract object from data (io.Reader OR []byte)
func FromJson(obj interface{}, data interface{}) error {
	switch data.(type) {
	case io.Reader:
		decoder := json.NewDecoder(data.(io.Reader))
		return decoder.Decode(obj)
	case []byte:
		return json.Unmarshal(data.([]byte), obj)
	}

	return ErrNotSupported
}
{{end}}

{{define "setter"}}
{{- if (intersection .field.type "string" "int64" "float64") }}
// Set{{.field.name}} set {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Set{{.field.name}}(v {{.type}}) {
    {{ .structname | firstLower }}.{{.field.name}} = v
}
{{ end }} 
{{ end }}

{{define "getter"}}
{{- if (intersection .field.type "string" "int64" "float64")  }}
// Get{{.field.name}} get {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Get{{.field.name}}() v {{.type}} {
    return {{ .structname | firstLower }}.{{.field.name}}
}
{{ end }} 
{{ end }}

{{define "structtype"}}
func (s *{{.struct}}) {{.name}}() string {
    return "{{.struct}}"
} 
{{end}}

{{ define "stringer" }}
{{range $index, $name := . -}}
var {{$name | toUpper}} = "{{$name | toUpper}}"
{{ end }} 
{{ end }}