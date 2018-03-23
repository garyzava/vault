package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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
	Paths   []Path
}

type Path struct {
	Pattern string
	Methods []Method
}

type Method struct {
	HTTPMethod string
	Summary    string
	Parameters []Parameter
}

type Parameter struct {
	Name string
	Type string
	In   string
}

type pathlet struct {
	pattern string
	params  map[string]bool
}

func parsePattern(root, pat string) []pathlet {
	paths := make([]pathlet, 0)
	pat = strings.TrimRight(pat, "$")

	reqd_re := regexp.MustCompile(`\(\?P<(\w+)>[^)]*\)`)

	params := make(map[string]bool)
	result := reqd_re.FindAllStringSubmatch(pat, -1)
	if result != nil {
		//println(pat)
		//println(len(result))
		for _, p := range result {
			par := p[1]
			params[par] = true
			pat = strings.Replace(pat, p[0], fmt.Sprintf("{%s}", par), 1)
		}
	}
	pat = fmt.Sprintf("/%s/%s", root, pat)
	paths = append(paths, pathlet{pat, params})
	return paths
}

func cleanSynopsis(syn string) string {
	if idx := strings.Index(syn, "\n"); idx != -1 {
		syn = syn[0:idx] + "…"
	}
	return syn
}

func procLogicalPath(p *framework.Path) []Path {
	var docPaths []Path
	methods := []Method{}

	paths := parsePattern("sys", p.Pattern)

	for opType := range p.Callbacks {
		m := Method{}
		switch opType {
		case logical.UpdateOperation:
			m.HTTPMethod = "post"
		case logical.DeleteOperation:
			m.HTTPMethod = "delete"
		default:
			m.HTTPMethod = "get"
		}

		//println(p.Pattern)

		d := make(map[string]bool)
		for _, path := range paths {
			for param := range path.params {
				d[param] = true
				m.Parameters = append(m.Parameters, Parameter{
					Name: param,
					Type: "string",
					In:   "path",
				})
			}
		}

		for name, field := range p.Fields {
			if _, ok := d[name]; !ok {
				m.Parameters = append(m.Parameters, Parameter{
					Name: name,
					Type: string(field.Type),
					In:   "body",
				})
			}
		}
		methods = append(methods, m)
	}
	pd := Path{
		Pattern: paths[0].pattern,
		//Summary: cleanSynopsis(p.HelpSynopsis),
		Methods: methods,
	}
	docPaths = append(docPaths, pd)

	return docPaths
}

func main() {
	c := vault.Core{}
	b := vault.NewSystemBackend(&c)

	doc := Doc{
		Version: version.GetVersion().Version,
		Paths:   make([]Path, 0),
	}

	for _, p := range b.Backend.Paths {
		paths := procLogicalPath(p)
		doc.Paths = append(doc.Paths, paths...)
	}

	// Define the functions
	funcs := map[string]interface{}{
		"indent": funcIndent,
	}

	// Parse the help template
	tmpl, _ := template.New("root").Funcs(funcs).Parse(tmpl)
	tmpl.Execute(os.Stdout, doc)

}
