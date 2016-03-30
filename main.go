package main

import (
    "os"
    // "github.com/BurntSushi/toml"
    "github.com/codegangsta/cli"
    "bufio"
    "fmt"
    "strings"
    "bytes"
    "io/ioutil"
)

var (
    CommentLine = "//"
    CommentBlockStart = "/*"
    CommentBlockEnd = "*/"
    
    PrefixConfig = CommentLine+"dico"
    PrefixStartGen = CommentLine+">>>dico:autogenerate>>>"
    PrefixEndGen = CommentLine+"<<<dico:autogenerate<<<"
    
)

func isConfigLine(str string) bool {
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
        
        var text = scanner.Text();
NEXTCONFIG:       
        if isConfigLine(text) {
            var config = strings.Replace(text, PrefixConfig, "", 2)
            
            fmt.Printf("%s:%d\t%s\n", path, lineCounter, text)
            
            fmt.Fprintln(buff, text) // for reuse
            
            if scanner.Scan() {
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
            fmt.Fprintln(buff, CommentLine + "\tfrom:\t", config)
            fmt.Fprintln(buff, "hello")
            fmt.Fprintln(buff, PrefixEndGen)
            
            if !isEndCodeGenerated(scanner.Text()) {
                
                // if the beginning of the next line generator
                if isConfigLine(scanner.Text()) {
                    goto NEXTCONFIG
                }
                
                fmt.Fprintln(buff, scanner.Text())
            }
            
            continue;
        }
        
        fmt.Fprintln(buff, text)    
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