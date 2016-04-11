//dico main templates/golang/*
//config.toml
//pkg = "main"
// errors = ["ErrNotFound", "err not allowed"]
// [values]
// foo = "foo"
// bar = "bar"
//[struct]
// name ="FooBarStruct"
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


var ErrNotFound = "ErrNotFound"
var ErrNotAllowed = "ErrNotAllowed"
 



type FooBarStruct struct {
    
    // comment
    Foo string `json:"Foo"`  
    
    
    Bar string   
    
}

// SetFoo set Foo
func (f *FooBarStruct) SetFoo(v string) {
    f.Foo = v
}
 

// GetFoo get Foo
func (f *FooBarStruct) GetFoo() string {
    return f.Foo
}
 

// SetBar set Bar
func (f *FooBarStruct) SetBar(v string) {
    f.Bar = v
}
 

// GetBar get Bar
func (f *FooBarStruct) GetBar() string {
    return f.Bar
}
 



func main() {
    s := &FooBarStruct{"foo", "bar"}
    fmt.Printf("%v", s)
}

//<<<AUTOGENERATE.DICO

