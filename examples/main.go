//dico main templates/golang/*
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
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  main templates/golang/*

package main

import (
    "fmt"
)



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

