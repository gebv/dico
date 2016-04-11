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

{{- if and (eq (hasPrefix .field.type "map") false) (eq (hasPrefix .field.type "[]") false) }}
// Set{{.field.name}} set {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Set{{.field.name}}(v {{.field.type}}) {
    {{ .structname | firstLower }}.{{.field.name}} = v
}
{{ end }} 

{{- if (hasPrefix .field.type "[]") }}
// Add{{.field.name}} add element {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Add{{.field.name}}(v {{substring .field.type 2}}) {
    if {{ .structname | firstLower }}.Include{{.field.name}}(v) {
        return
    }
    
    {{ .structname | firstLower }}.{{.field.name}} = append({{ .structname | firstLower }}.{{.field.name}}, v)
}

// Remove{{.field.name}} remove element {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Remove{{.field.name}}(v {{substring .field.type 2}}) {
    if !{{ .structname | firstLower }}.Include{{.field.name}}(v) {
        return
    }
    
    _i := {{ .structname | firstLower }}.Index{{.field.name}}(v)
    
    {{ .structname | firstLower }}.{{.field.name}} = append({{ .structname | firstLower }}.{{.field.name}}[:_i], {{ .structname | firstLower }}.{{.field.name}}[_i+1:]...)
}
{{ end }}

{{- if (hasPrefix .field.type "map") }}
{{ $regexp := "map\\[(?P<key>[a-zA-Z0-9{}]+)\\](?P<item>[a-zA-Z0-9{}]+)" }}
{{ $keyType := (index (regexp .field.type $regexp) 1)}}
{{ $valueType := (index (regexp .field.type $regexp) 2)}}

// Set{{.field.name}} set all elements {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Set{{.field.name}}(v {{.field.type}}) {
    {{ .structname | firstLower }}.{{.field.name}} = make({{.field.type}})
    
    for key, value := range v {
        {{ .structname | firstLower }}.{{.field.name}}[key] = value
    }
}

// Add{{.field.name}} add element by key
func ({{ .structname | firstLower }} *{{.structname}}) SetOne{{.field.name}}(k {{ $keyType}}, v {{ $valueType }}) {
    {{ .structname | firstLower }}.{{.field.name}}[k] = v
}

// Remove{{.field.name}} remove element by key
func ({{ .structname | firstLower }} *{{.structname}}) Remove{{.field.name}}(k {{ $keyType}}) {
    if _, exist := {{ .structname | firstLower }}.{{.field.name}}[k]; exist {
        delete({{ .structname | firstLower }}.{{.field.name}}, k)  
    } 
}
{{ end }}

{{ end }}

{{define "getter"}}

// Get{{.field.name}} get {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Get{{.field.name}}() {{.field.type}} {
    return {{ .structname | firstLower }}.{{.field.name}}
}

{{- if (hasPrefix .field.type "[]") }}
// Index{{.field.name}} get index element {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Index{{.field.name}}(v {{substring .field.type 2}}) int {
    for _index, _v := range {{ .structname | firstLower }}.{{.field.name}} {
        if _v == v {
            return _index
        }
    }
    return -1
}

// Include{{.field.name}} has exist value {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Include{{.field.name}}(v {{substring .field.type 2}}) bool {
    return {{ .structname | firstLower }}.Index{{.field.name}}(v) > -1
}
{{ end }}

{{- if (hasPrefix .field.type "map") }}
{{ $regexp := "map\\[(?P<key>[a-zA-Z0-9{}]+)\\](?P<item>[a-zA-Z0-9{}]+)" }}
{{ $keyType := (index (regexp .field.type $regexp) 1)}}
{{ $valueType := (index (regexp .field.type $regexp) 2)}}

// Exist{{.field.name}} has exist key {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) ExistKey{{.field.name}}(k {{ $keyType}}) bool {
     _, exist := {{ .structname | firstLower }}.{{.field.name}}[k]
     
     return exist
}

func ({{ .structname | firstLower }} *{{.structname}}) GetOne{{.field.name}}(k {{ $keyType}}) {{ $valueType }} {
    return {{ .structname | firstLower }}.{{.field.name}}[k]
}
{{ end }}

{{ end }}

{{ define "stringer" }}
{{range $index, $name := . -}}
var {{$name | toUpper}} = "{{$name | toUpper}}"
{{ end }} 
{{ end }}