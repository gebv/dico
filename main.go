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
    "path/filepath"
)

var (
    Version = "0.0.2"
    
    CommentLine = "//"
    CommentBlockStart = "/*"
    CommentBlockEnd = "*/"
    
    PrefixConfig = CommentLine+"dico"
    PrefixStartGen = CommentLine+"AUTOGENERATE.DICO>>>"
    PrefixEndGen = CommentLine+"<<<AUTOGENERATE.DICO"
    
    PrefixStartConfig = CommentLine+"config."
    PrefixEndConfig = CommentLine+"config."
    
    HelpMessage = "The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it"
    
    DEBUG = false
)

var HelpfullTemplateFuncs = template.FuncMap{
    
}

func NewGenerator(args []string, config *Config) (*Generator, error) {
    if len(args) < 2 {
        return nil, fmt.Errorf("not valid command arguments %v", args);
    }
    
    // // TODO: check template errors
    var tpl *template.Template
    var err error 
    
    // TODO: check file 
    
    // b, err := ioutil.ReadFile(args[0])
    
    // if err != nil {
    //     return nil, err
    // }
    
    // s := string(b)
    
    // // parse execute file
    tpl = template.New(args[0]).Funcs(HelpfullTemplateFuncs)
    // tpl, err = template.New(args[0]).Funcs(HelpfullTemplateFuncs).Parse(s)
    
    // if err != nil {
    //     return nil, err
    // }
    
    // parse tempaltes
    pwd, _ := os.Getwd()
    for _, tplpath := range args[1:]{
        tpl, err = tpl.ParseGlob(pwd + string(os.PathSeparator) + tplpath)
    }
    
    if err != nil {
        return nil, err
    }
    
    tpl = template.Must(tpl, err)
    
    config.TplName = args[0]
    
    return &Generator{config, tpl}, nil
}

type Generator struct {
    Config *Config
    Tpl *template.Template
}

func (g *Generator) Compile() (string, error) {
    var b = bytes.NewBufferString("")
    
    g.Config.Config["ENV"] = g.Config.Env
    
    // if err := g.Tpl.Execute(b, g.Config.Config); err != nil {
        
    //     return "", err
    // }
    
    if err := g.Tpl.ExecuteTemplate(b, g.Config.TplName, g.Config.Config); err != nil {
        
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
    TplName string // output template name
     
    Type string // type config toml or ...
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
                    
                    scanner.Scan()
                    // if !scanner.Scan() {
                        
                    //     return buff.Bytes(), fmt.Errorf("unexpected ending")
                    // }
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
                
                // for debug
                if DEBUG {
                   fmt.Fprintln(buff, CommentLine + "[DICO.CONFIG]:\t" + fmt.Sprintf("%+v", EmbeddedConfig.Config)) 
                } 
                
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

func analyzeFile(pattern string) func (fp string, fi os.FileInfo, err error) error {
    return func (fp string, fi os.FileInfo, err error) error {
        // https://rosettacode.org/wiki/Walk_a_directory/Recursively#Go
        
        if err != nil {
            fmt.Println(err) // can't walk here,
            return nil       // but continue walking elsewhere
        }
        if fi.IsDir() {
            return nil // not a file.  ignore.
        }
        
        matched, err := filepath.Match(pattern, fi.Name())
        
        if err != nil {
            fmt.Println(err) // malformed pattern
            return err       // this is fatal.
        }
        if !matched {
            return nil
        }
        
        if out, err := analyzeAndGenerate(fp); err == nil {
            err := ioutil.WriteFile(fp, out, 0644)
            if err != nil {
                fmt.Println(err)
            }
        } else {
            fmt.Println(err)
        }
        
        return nil
    }
} 

func main() {    
    app := cli.NewApp()
    app.Action = func(c *cli.Context) {
        path := c.Args().Get(0)
        
        filepath.Walk(path, analyzeFile(c.Args().Get(1)))
    }
    
    fmt.Printf("%v\n", os.Args)
    
    app.Run(os.Args)
}


/*
//dico main templates/golang/*
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.1
//[DICO.COMMAND]:	  main templates/golang/*
//<<<AUTOGENERATE.DICO
*/
