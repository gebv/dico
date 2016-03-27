# dynamo
text generator from template and сonfig

Any text generation. For example code of programm
Give a description of the app module (with the specific dynamics) and generate code component

Example for golang

Config:
``` toml
pkg_name = 'entity'

[[fields]]
name = 'Foo'
type = 'string'
```

Template:
``` tpl
package {{ .PkgName }}

{{ $field := ragne .Fields }}
var {{$field.Name }} {{$field.Type}}
{{ end }}
```

Output
``` go
package entity

var Foo string
```
