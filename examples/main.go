//dico ../templates/golang/* main.go.tpl
//config.toml
//pkg = "main"
//[values]
// foo="Foo"
// bar="Bar"
//[struct]
//name="FooBarStruct"
// [[struct.fields]]
// comment = "comment"
// name = "Foo"
// type = "string"
// tag = '''json:"Foo"'''
// [[struct.fields]]
// name = "Bar"
// type = "string"
//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.1
//[DICO.COMMAND]:	  ../templates/golang/* main.go.tpl
//[DICO.CONFIG]:	map[pkg:main values:map[foo:Foo bar:Bar] struct:map[name:FooBarStruct fields:[map[comment:comment name:Foo type:string tag:json:"Foo"] map[name:Bar type:string]]]]
package main

import (
    "fmt"
)


// FooBarStruct <no value>
type FooBarStruct struct {
    
    // comment
    Foo string `json:"Foo"`  
    
    
    Bar string   
    
}


func main() {
    s := &FooBarStruct{"Foo", "Bar"}
    fmt.Printf("%v", s)
}
//<<<AUTOGENERATE.DICO

