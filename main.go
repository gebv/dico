package main

import (
    "os"
    "github.com/BurntSushi/toml"
    "github.com/codegangsta/cli"
    "bufio"
    "fmt"
    "strings"
    "bytes"
    "io/ioutil"
    "text/template"
)

var (
    Version = "0.0.1"
    
    CommentLine = "//"
    CommentBlockStart = "/*"
    CommentBlockEnd = "*/"
    
    PrefixConfig = CommentLine+"dico"
    PrefixStartGen = CommentLine+"AUTOGENERATE.DICO>>>"
    PrefixEndGen = CommentLine+"<<<AUTOGENERATE.DICO"
    
    PrefixStartConfig = CommentLine+"config."
    PrefixEndConfig = CommentLine+"config."
    
    HelpMessage = "The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it"
    
    DEBUG = true
)

var HelpfullTemplateFuncs = template.FuncMap{
    
}

func NewGenerator(args []string, config *Config) (*Generator, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("not valid command arguments");
    }
    
    // // TODO: check template errors
    var tpl *template.Template
    var err error 
    
    // TODO: check file 
    
    b, err := ioutil.ReadFile(args[1])
    
    if err != nil {
        return nil, err
    }
    
    s := string(b)
    
    // parse execute file
    tpl, err = template.New("").Funcs(HelpfullTemplateFuncs).Parse(s)
    
    if err != nil {
        return nil, err
    }
    
    // parse tempaltes
    tpl, err = tpl.ParseGlob(args[0])
    
    if err != nil {
        return nil, err
    }
    
    tpl = template.Must(tpl, err)
    
    return &Generator{config, tpl}, nil
}

type Generator struct {
    Config *Config
    Tpl *template.Template
}

func (g *Generator) Compile() (string, error) {
    var b = bytes.NewBufferString("")
    
    g.Config.Config["ENV"] = g.Config.Env
    
    if err := g.Tpl.Execute(b, g.Config.Config); err != nil {
        
        return "", err
    }
    
    return b.String(), nil
}

type ConfigParser func (raw string, v interface{}) (error)

var registredConfigTypes = map[string]ConfigParser{
    "toml": func (raw string, config interface{}) (error) {
        
        _ , err := toml.Decode(raw, config)
        
        return err
    },
}

func NewConfig(typeName string, configRaw string) *Config {
    // TODO: Check allowable types
    
    parser, supported := registredConfigTypes[typeName]
    
    return &Config{
        Type: typeName,
        ConfigRaw: configRaw,
        Config: make(map[string]interface{}),
        parser: parser,
        supportedType: supported,
        Env: make(map[interface{}]interface{}), // TODO: values for the application
    }
}

type Config struct {
    Type string
    ConfigRaw string
    
    Config map[string]interface{}
    Env map[interface{}]interface{}
    
    parser ConfigParser
    supportedType bool
}

func (g *Config) IsSupportedType() bool {
    
    return g.supportedType 
}

func (g *Config) BuildConfig() error {
    if g.IsSupportedType() {
        
        return g.parser(g.ConfigRaw, &g.Config)
    }
    
    return fmt.Errorf("not supported type")
}

func (g *Config) String() string {
    
    return "generated text";
}

func removeCommentSymbols(str string) string {
    return strings.Replace(str, CommentLine, "", -1);
}

func getTypeConfig(str string) string {
    return strings.Replace(str, PrefixStartConfig, "", 1);
}

func isStartDicoConfig(str string) bool {
    return strings.HasPrefix(str, PrefixStartConfig)
}

func isEndDicoConfig(str string) bool {
    return strings.HasPrefix(str, PrefixEndConfig)
}

func isDicoCommand(str string) bool {
    return strings.HasPrefix(str, PrefixConfig)
}

func isStartCodeGenerated(str string) bool {
    return strings.HasPrefix(str, PrefixStartGen)
}

func isEndCodeGenerated(str string) bool {
    return strings.HasPrefix(str, PrefixEndGen)
}

func analyzeAndGenerate(path string) ([]byte, error) {
    buff := bytes.NewBuffer([]byte{});
    
    inFile, err := os.Open(path)
    
    if err != nil {
        return buff.Bytes(), err
    }
    
    defer inFile.Close()
    
    scanner := bufio.NewScanner(inFile)
    scanner.Split(bufio.ScanLines)
    
    lineCounter := 0
    
    for scanner.Scan() {
        lineCounter++;
        
NEXT:       
        if isDicoCommand(scanner.Text()) {
            var command = strings.Replace(scanner.Text(), PrefixConfig, "", 1)
            var args = strings.Fields(command)
            var config = bytes.NewBuffer([]byte{});
            
            fmt.Printf("%s:%d\t%s\n", path, lineCounter, command)
            
            fmt.Fprintln(buff, scanner.Text()) // for reuse
            
            var EmbeddedConfig *Config
            
            // next row
            if scanner.Scan() {
                
                // detected embedded config
                if isStartDicoConfig(scanner.Text()) {
                    // first row config
                    fmt.Fprintln(buff, scanner.Text())
                    
                    // type config
                    var typeConfig = getTypeConfig(scanner.Text())
                    
                    for scanner.Scan() {
                        fmt.Fprintln(buff, scanner.Text())
                        
                        if isEndDicoConfig(scanner.Text()) {
                            break
                        }
                        
                        fmt.Fprintln(config, removeCommentSymbols(scanner.Text()))
                    }
                    
                    EmbeddedConfig = NewConfig(typeConfig, config.String());
                    
                    if !scanner.Scan() {
                        return buff.Bytes(), fmt.Errorf("unexpected ending")
                    }
                }
                
                if isStartCodeGenerated(scanner.Text()) {
                    for scanner.Scan() {
                        if isEndCodeGenerated(scanner.Text()) {
                            // clear the previous generated code
                            break
                        }
                    }       
                }
            }
            
            fmt.Fprintln(buff, PrefixStartGen)
            fmt.Fprintln(buff, CommentLine + "\t" + HelpMessage)
            fmt.Fprintln(buff, CommentLine + "[DICO.VERSION]:\t", Version)
            fmt.Fprintln(buff, CommentLine + "[DICO.COMMAND]:\t", command)
            
            if EmbeddedConfig != nil {
                
                err := EmbeddedConfig.BuildConfig()
                
                fmt.Fprintln(buff, CommentLine + "[DICO.CONFIG]:\t" + fmt.Sprintf("%+v", EmbeddedConfig.Config))
                
                if err != nil && DEBUG {
                    fmt.Fprintln(buff, CommentLine + "[DICO.ERRORS.COMPLIE_CONFIG]:\t" + err.Error())    
                }
                
                var g *Generator
                g, err = NewGenerator(args, EmbeddedConfig)
                
                if g != nil && err == nil{
                    gentext, err := g.Compile()
                    
                    if err == nil {
                        fmt.Fprintln(buff, gentext)
                            
                    } else {
                        fmt.Fprintln(buff, CommentLine + "[DICO.ERRORS.COMPLIE_TEXT]:\t" + err.Error())    
                    }
                    
                        
                } else if (err != nil) {
                    fmt.Fprintln(buff, CommentLine + "[DICO.ERRORS.COMPLIE_TEXT]:\t" + err.Error())
                }
            }
            
            
            fmt.Fprintln(buff, PrefixEndGen)
            
            if !isEndCodeGenerated(scanner.Text()) {
                
                // if the beginning of the next line Config
                if isDicoCommand(scanner.Text()) {
                    goto NEXT
                }
                
                fmt.Fprintln(buff, scanner.Text())
            }
            
            continue;
        }
        
        fmt.Fprintln(buff, scanner.Text())    
    }
    
    return buff.Bytes(), nil
}

func main() {    
    app := cli.NewApp()
    app.Action = func(c *cli.Context) {
        file := c.Args().Get(0)
        
        if out, err := analyzeAndGenerate(file); err == nil {
            err := ioutil.WriteFile(file, out, 0644)
            if err != nil {
                fmt.Println(err)
            }
        } else {
            fmt.Println(err)
        }
        
    }

    app.Run(os.Args)
}


/*
//dico --t=example

console.log("start");
//dico --t=example2

console.log("end");
*/