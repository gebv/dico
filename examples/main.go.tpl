package {{.pkg}}

import (
    "fmt"
)

{{ template "struct" .struct }}

func main() {
    s := &{{.struct.name}}{"{{.values.foo}}", "{{.values.bar}}"}
    fmt.Printf("%v", s)
}