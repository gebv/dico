{{define "main"}}
package {{.pkg}}

import (
    "fmt"
)

{{template "stringer" .errors}}

{{template "struct" .struct}}

func main() {
    s := New{{.struct.name}}()
    fmt.Printf("%v", s)
}
{{end}}