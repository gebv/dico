# dico
embedded code generator, text generator from templates and сonfigs

Генератор кода, конфигов, текста на основе преднастроек и шаблонов.

![demo](https://s3.amazonaws.com/idheap/ss/ezgif.com-video-to-gif_1.gif)

Для запуска генератора следует выполнить следующую команду

```
dico ./src *.go
```

> В случае языка golang рекомендуется после генерации выполнить форматирование
> `gofmt -w ./src`
> [gofmt](https://blog.golang.org/go-fmt-your-code)

первым аругментом `./src` указана дириктория с исходным кодом программы.
`*.go` указана маска по которой будут анализироваться файлы.
В нашем примере будут подвергнуты все файлы с расширением `go` в дириктории `src`.

* [Overview](https://github.com/gebv/dico#Описание)
* [Setup](https://github.com/gebv/dico#Установка)
* [Templates](https://github.com/gebv/dico#Шаблоны)
* [Examples](https://github.com/gebv/dico#examples)
 
# Описание

Положим мы пишем код программы. 
Встраиваете в ваш исходный код специальный текст начинающийся с новой строки `//dico output_name path_templates [, ...]`. Условно назовем это командой.
`output_name` определяет исполняемый шаблон,
один и более `path_templates` определяет дириктории с шаблонами.

```
//dico main templates/golang/*
```

В момент анализа вашего файла генератор анализирует комманду и генерирует код.

Имеется возможность встраивать настройки для генерируемого кода.
Настройки начинаются и заканчиваются c `//config.extension`.
`extension` определяет формат настроек. Поддерживается настройки в формате [toml](https://github.com/toml-lang/toml)).

```
//config.toml
//pkg = "main"
// errors = ["ErrNotFound", "err not allowed"]
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
// comment = "comment 2"
// name = "Bar"
// type = "string"
//config.toml
```

После генерации исходная команда и настройки остаются без изменений. Ниже добавляется секция с авто сгенерированным кодом.
Авто сгенерированный код находится в секции начинающийся с `//AUTOGENERATE.DICO>>>` и заканчивающийся `//<<<AUTOGENERATE.DICO`.

```
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
func (f *FooBarStruct) GetFoo() v string {
    return f.Foo
}
 

// SetBar set Bar
func (f *FooBarStruct) SetBar(v string) {
    f.Bar = v
}
 

// GetBar get Bar
func (f *FooBarStruct) GetBar() v string {
    return f.Bar
}
 



func main() {
    s := &FooBarStruct{"foo", "bar"}
    fmt.Printf("%v", s)
}

//<<<AUTOGENERATE.DICO

```

В секции автосгенерированного кода находится отладочная информация отражающая версию `[DICO.VERSION]` генератора и исходную команду `[DICO.COMMAND]`.
`[DICO.ERRORS.COMPLIE_TEXT]` с описанием ошоибки генерации кода и `[DICO.ERRORS.COMPLIE_CONFIG]` с описание ошибки анализа конфига.

Например

```
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  main
//[DICO.ERRORS.COMPLIE_TEXT]:	not valid command arguments [main]
//<<<AUTOGENERATE.DICO
```

Сказано что не верная команда генерации, не валидные аргументы команды.

# Установка

[Установить gоlang](https://golang.org/doc/install).

```
$ mkdir -p dico
$ cd dico
$ export GOPATH=$(pwd)
$ go get github.com/gebv/dico
$ tree -L 2
.
├── bin
│   └── dico
...

```

`bin/dico` является бинарным приложением 
Рекомендуется добавить дирикторию `bin/` в переменную окружения `PATH`

# Шаблоны

Шаблоны это [golang template](https://golang.org/pkg/text/template/).
Для удобства они находятся в различных папках и файлах разбитые по смыслу. [Например](templates).
Шаблоны необхдоимы для генерации кода. Следует разбить генерируемый код на простые конфигурируемые вставки.

Например для генерации структуры (язык golang)
```
{{define "struct" }}
{{with .comment }}// {{.name}} {{.comment}}{{end}}
type {{.name}} struct {
    {{ range $key, $field := .fields }}
    {{with $field.comment}}// {{$field.comment}}{{end}}
    {{$field.name}} {{$field.type}} {{template "structtags" $field.tag}}  
    {{ end }}
}
{{end}}

{{define "structtags" }}{{with .}}`{{.}}`{{end}}{{end}}
```

Структура ожидает на входе конфиг описывающий
* имя структуры `name`
* описание структуры `comment`
* описание полей структуры `fields`
* каждое поле отражает тип и описание поля

# Examples

## nginx настройки

Генерация настроек nginx по [шаблонам](templates)

```
//dico server templates/nginx/*
//config.toml
// listen="8080"
// root="/var/www/site"
// [[locations]]
// path="/images"
// root="/data"
// [[locations]]
// path="/"
// proxy_pass = '127.0.0.1:8081'
//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  server templates/nginx/*

server {
    listen 8080;
    
        
location /images {
    
    root /data;
}

    
        
location / {
    proxy_pass 127.0.0.1:8081;
    
}

    
}

//<<<AUTOGENERATE.DICO


```

[Больше примеров](examples)

TODO
* Автогенерация документации по шаблонам (парсить define из шаблонов и брать комментарии)