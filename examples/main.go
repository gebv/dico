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
// name = "M"
// type = "map[string]interface{}"
// [[struct.fields]]
// name = "I"
// type = "int"
// [[struct.fields]]
// name = "S"
// type = "string"
// [[struct.fields]]
// name = "A"
// type = "[]string"
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
 



func NewFooBarStruct() *FooBarStruct {
    model := new(FooBarStruct)
    
    
    model.M = make(map[string]interface{})
      
    
      
    
      
    
      
    
    return model
}
type FooBarStruct struct {
    
    
    M map[string]interface{}   
    
    
    I int   
    
    
    S string   
    
    
    A []string   
    
}





// SetM set all elements M
func (f *FooBarStruct) SetM(v map[string]interface{}) {
    f.M = make(map[string]interface{})
    
    for key, value := range v {
        f.M[key] = value
    }
}

// AddM add element by key
func (f *FooBarStruct) SetOneM(k string, v interface{}) {
    f.M[k] = v
}

// RemoveM remove element by key
func (f *FooBarStruct) RemoveM(k string) {
    if _, exist := f.M[k]; exist {
        delete(f.M, k)  
    } 
}




// GetM get M
func (f *FooBarStruct) GetM() map[string]interface{} {
    return f.M
}




// ExistM has exist key M
func (f *FooBarStruct) ExistKeyM(k string) bool {
     _, exist := f.M[k]
     
     return exist
}

func (f *FooBarStruct) GetOneM(k string) interface{} {
    return f.M[k]
}



// SetI set I
func (f *FooBarStruct) SetI(v int) {
    f.I = v
}




// GetI get I
func (f *FooBarStruct) GetI() int {
    return f.I
}


// SetS set S
func (f *FooBarStruct) SetS(v string) {
    f.S = v
}




// GetS get S
func (f *FooBarStruct) GetS() string {
    return f.S
}


// AddA add element A
func (f *FooBarStruct) AddA(v string) {
    if f.IncludeA() {
        return
    }
    
    f.A = append(f.A, v)
}

// RemoveA remove element A
func (f *FooBarStruct) RemoveA(v string) {
    if !f.IncludeA(v) {
        return
    }
    
    _i := f.IndexA(v)
    
    f.A = append(f.A[:_i], f.A[_i+1:]...)
}




// GetA get A
func (f *FooBarStruct) GetA() []string {
    return f.A
}
// IndexA get index element A
func (f *FooBarStruct) IndexA(v string) int {
    for _index, _v := range f.A {
        if _v == v {
            return _index
        }
    }
    return -1
}

// IncludeA has exist value A
func (f *FooBarStruct) IncludeA(v string) bool {
    return f.IndexA(v) > -1
}





func main() {
    s := NewFooBarStruct()
    fmt.Printf("%v", s)
}

//<<<AUTOGENERATE.DICO

