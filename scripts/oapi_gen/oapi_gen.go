package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
)

func funcIndent(count int, text string) string {
	var buf bytes.Buffer
	prefix := strings.Repeat(" ", count)
	scan := bufio.NewScanner(strings.NewReader(text))
	for scan.Scan() {
		buf.WriteString(prefix + scan.Text() + "\n")
	}

	return strings.TrimRight(buf.String(), "\n")
}

type Doc struct {
	Version string
	Paths   []PathDoc
}

type Parameter struct {
	Name string
	Type string
	In   string
}

type PathDoc struct {
	Pattern    string
	Method     string
	Summary    string
	Parameters []Parameter
}

func cleanPattern(pat string) string {
	pat = "/sys/" + pat
	pat = strings.TrimRight(pat, "$")
	return pat
}

func main() {
	c := vault.Core{}
	b := vault.NewSystemBackend(&c)

	doc := Doc{
		Version: version.GetVersion().Version,
		Paths:   make([]PathDoc, 0),
	}

	for _, path := range b.Backend.Paths {
		method := "get"
		if _, ok := path.Callbacks[logical.DeleteOperation]; ok {
			method = "delete"
		}

		pd := PathDoc{
			Pattern:    cleanPattern(path.Pattern),
			Method:     method,
			Summary:    path.HelpSynopsis,
			Parameters: []Parameter{{"a", "b", "c"}, {"1", "2", "3"}},
		}
		doc.Paths = append(doc.Paths, pd)
	}

	//const prefix = "/sys/"
	//sort.Strings(paths)
	//for _, path := range paths {
	//	fmt.Printf("%s%s\n", prefix, path)
	//}
	// Define the functions
	funcs := map[string]interface{}{
		"indent": funcIndent,
	}

	// Parse the help template
	tmpl, _ := template.New("root").Funcs(funcs).Parse(tmpl)
	tmpl.Execute(os.Stdout, doc)

}
