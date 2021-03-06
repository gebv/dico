package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
)

var (
	Version = "0.0.3"

	CommentLine       = "//"
	CommentBlockStart = "/*"
	CommentBlockEnd   = "*/"

	PrefixConfig   = CommentLine + "dico"
	PrefixStartGen = CommentLine + "AUTOGENERATE.DICO>>>"
	PrefixEndGen   = CommentLine + "<<<AUTOGENERATE.DICO"

	PrefixStartConfig = CommentLine + "config."
	PrefixEndConfig   = CommentLine + "config."

	HelpMessage = "The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it"

	DEBUG = false
)

var HelpfullTemplateFuncs = template.FuncMap{
	"regexp": func(s, r string) []string {
		// map\\[(?P<key>[a-zA-Z0-9{}]+)\\](?P<item>[a-zA-Z0-9{}]+)
		re := regexp.MustCompile(r)
		// re.SubexpNames()
		return re.FindStringSubmatch(s)
	},
	"substring": func(s string, ii ...int) string {
		if len(ii) == 0 {
			return s
		}

		var start = 0
		var end = len(s)

		if len(ii) >= 1 {
			start = ii[0]
		}

		if start < 0 || start > len(s) {
			start = 0
		}

		if len(ii) == 2 {
			end = ii[1]
		}

		if end < 0 || end > len(s) {
			end = len(s)
		}

		if start > end && len(ii) == 2 {
			start = 0
			end = len(s)
		}

		return s[start:end]
	},
	"intersection": func(v interface{}, vv ...interface{}) bool {
		switch vv[0].(type) {
		case func(string, string) bool:
			fn := vv[0].(func(string) bool)

			// TODO: check type string

			for _, _v := range vv {
				if fn(_v.(string)) {
					return true
				}
			}
		default:
		}

		for _, _v := range vv {
			if _v == v {
				return true
			}
		}

		return false
	},
	"hasPrefix": func(s1, s2 string) bool {

		return strings.HasPrefix(s1, s2)
	},
	"fnHasPrefix": func(prefix string) func(string) bool {
		return func(s1 string) bool {
			return strings.HasPrefix(s1, prefix)
		}
	},
	"hasSuffix": func(s1, s2 string) bool {

		return strings.HasSuffix(s1, s2)
	},
	"fnHasSuffix": func(suffix string) func(string) bool {
		return func(s1 string) bool {
			return strings.HasSuffix(s1, suffix)
		}
	},
	"map": func(vv ...interface{}) (res map[interface{}]interface{}) {
		res = make(map[interface{}]interface{})

		if len(vv) == 0 || len(vv)%2 != 0 {

			return
		}

		for i := 0; i < len(vv); i += 2 {
			res[vv[i]] = vv[i+1]
		}

		return
	},
	"array": func(vv ...interface{}) []interface{} {
		return vv
	},
	"setter": func(m map[interface{}]interface{}, key, value interface{}) interface{} {
		m[key] = value
		return value
	},
	"getter": func(m map[interface{}]interface{}, key interface{}) interface{} {
		return m[key]
	},

	// toUpper transform "str str" to "Str<separator>Str"
	"toUpper": func(ss ...string) string {
		if len(ss) == 0 {
			return ""
		}

		res := bytes.NewBufferString("")
		hasSep := true // first symbol up
		replace := ""

		if len(ss) == 2 {
			replace = ss[1]
		}

		for index, runeValue := range ss[0] {
			if !unicode.IsLetter(runeValue) {
				hasSep = true
				continue
			}

			if hasSep {
				if index > 0 {
					res.WriteString(replace)
				}

				runeValue = unicode.ToUpper(runeValue)
			}

			res.WriteRune(runeValue)

			hasSep = false
		}

		return res.String()
	},
	// toLower transform "str str" to "str<separator>str"
	"toLower": func(ss ...string) string {
		if len(ss) == 0 {
			return ""
		}

		res := bytes.NewBufferString("")
		hasSep := false
		replace := ""

		if len(ss) == 2 {
			replace = ss[1]
		}

		for index, runeValue := range ss[0] {
			if !unicode.IsLetter(runeValue) {
				hasSep = true
				continue
			}

			if unicode.IsUpper(runeValue) {

				runeValue = unicode.ToLower(runeValue)
				hasSep = true
			}

			if index > 0 && hasSep {
				res.WriteString(replace)
			}

			// IsLower
			res.WriteRune(runeValue)

			hasSep = false
		}

		return res.String()
	},
	"firstLower": func(s string) string {
		if len(s) == 0 {
			return ""
		}

		return string(unicode.ToLower([]rune(s)[0]))
	},
}

func NewGenerator(args []string, config *Config) (*Generator, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("not valid command arguments %v", args)
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

	path := os.Getenv("DICO_TEMPLATES")

	if len(path) > 0 {
		tpl, err = tpl.ParseGlob(path)

		if err != nil {
			return nil, err
		}
	}

	pwd, _ := os.Getwd()
	for _, tplpath := range args[1:] {
		tpl, err = tpl.ParseGlob(pwd + string(os.PathSeparator) + tplpath)

		if err != nil {
			return nil, err
		}
	}

	tpl = template.Must(tpl, err)

	config.TplName = args[0]

	if DEBUG {
		fmt.Printf("defined tempaltes: %v\n", tpl.DefinedTemplates())
	}

	return &Generator{config, tpl}, nil
}

type Generator struct {
	Config *Config
	Tpl    *template.Template
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

type ConfigParser func(raw string, v interface{}) error

var registredConfigTypes = map[string]ConfigParser{
	"toml": func(raw string, config interface{}) error {

		_, err := toml.Decode(raw, config)

		return err
	},
}

func NewConfig(typeName string, configRaw string) *Config {
	// TODO: Check allowable types

	parser, supported := registredConfigTypes[typeName]

	return &Config{
		Type:          typeName,
		ConfigRaw:     configRaw,
		Config:        make(map[string]interface{}),
		parser:        parser,
		supportedType: supported,
		Env:           make(map[interface{}]interface{}), // TODO: values for the application
	}
}

type Config struct {
	TplName string // output template name

	Type      string // type config toml or ...
	ConfigRaw string

	Config map[string]interface{}
	Env    map[interface{}]interface{}

	parser        ConfigParser
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

	return "generated text"
}

func removeCommentSymbols(str string) string {
	return strings.Replace(str, CommentLine, "", -1)
}

func getTypeConfig(str string) string {
	return strings.Replace(str, PrefixStartConfig, "", 1)
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
	buff := bytes.NewBuffer([]byte{})

	inFile, err := os.Open(path)

	if err != nil {
		return buff.Bytes(), err
	}

	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	lineCounter := 0

	for scanner.Scan() {
		lineCounter++

	NEXT:
		if isDicoCommand(scanner.Text()) {
			var command = strings.Replace(scanner.Text(), PrefixConfig, "", 1)
			var args = strings.Fields(command)
			var config = bytes.NewBuffer([]byte{})

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

					EmbeddedConfig = NewConfig(typeConfig, config.String())

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
			fmt.Fprintln(buff, CommentLine+"\t"+HelpMessage)
			fmt.Fprintln(buff, CommentLine+"[DICO.VERSION]:\t", Version)
			fmt.Fprintln(buff, CommentLine+"[DICO.COMMAND]:\t", command)

			if EmbeddedConfig != nil {

				err := EmbeddedConfig.BuildConfig()

				// for debug
				if DEBUG {
					fmt.Fprintln(buff, CommentLine+"[DICO.CONFIG]:\t"+fmt.Sprintf("%+v", EmbeddedConfig.Config))
				}

				if err != nil && DEBUG {
					fmt.Fprintln(buff, CommentLine+"[DICO.ERRORS.COMPLIE_CONFIG]:\t"+err.Error())
				}

				var g *Generator
				g, err = NewGenerator(args, EmbeddedConfig)

				if g != nil && err == nil {
					gentext, err := g.Compile()

					if err == nil {
						fmt.Fprintln(buff, gentext)

					} else {
						fmt.Fprintln(buff, CommentLine+"[DICO.ERRORS.COMPLIE_TEXT]:\t"+err.Error())
					}

				} else if err != nil {
					fmt.Fprintln(buff, CommentLine+"[DICO.ERRORS.COMPLIE_TEXT]:\t"+err.Error())
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

			continue
		}

		fmt.Fprintln(buff, scanner.Text())
	}

	return buff.Bytes(), nil
}

func analyzeFile(pattern string) func(fp string, fi os.FileInfo, err error) error {
	return func(fp string, fi os.FileInfo, err error) error {
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
	app.Action = func(c *cli.Context) error {
		path := c.Args().Get(0)

		return filepath.Walk(path, analyzeFile(c.Args().Get(1)))
	}

	if DEBUG {
		fmt.Printf("%v\n", os.Args)
	}

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
