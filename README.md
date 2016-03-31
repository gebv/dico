# dico
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

# Examples

Генерация nginx настроек

Создадим файл `nginx.conf`

```
//dico ../templates/nginx/* nginx.conf.tpl
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
```

Разберем файл

`//dico ../templates/nginx/* nginx.conf` говорит парсеру что есть команда для генерации.
Первый аргумент `../templates/nginx/*` указывает путь к шаблонам (заготовкам), которыми мы оперируем в `nginx.conf`.
Вторым аргументом `nginx.conf` указывается файл который следует пропарсить и сгенерирорвать код в определенных местах (все что начинается с `//dico` парсер определяет как комманда для генерации).

`//config.toml` обозначает начало и конец настроек (все что после точки определяет формат в котором представлены настройки, в нашем примере [toml](https://github.com/toml-lang/toml))
Настроки передаются в 

После работы парсера
```
$ dico nginx.config
```

Исходный файл выглядит следующим образом
```
//dico ../templates/nginx/* nginx.conf.tpl
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
//[DICO.VERSION]:	 0.0.1
//[DICO.COMMAND]:	  ../templates/nginx/* nginx.conf.tpl
//[DICO.CONFIG]:	map[listen:8080 root:/var/www/site locations:[map[path:/images root:/data] map[proxy_pass:127.0.0.1:8081 path:/]]]

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

Секция начинающаяся с `AUTOGENERATE.DICO>>>` и заканчивающая `<<<AUTOGENERATE.DICO` определяет автоматически сгенерированный код\текст. Внутри секции находятся вспомогательные параметры (отражающие версию генератора и прочее) и сам сгенерированный код\текст.

[Больше примеров](examples)

# TODO

* вторым аргументом указывать дирикторию в которой следует анализировать файлы по шаблону (например `dico src/*.go`)
