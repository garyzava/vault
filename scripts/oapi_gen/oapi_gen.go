package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
)

// regexen
var optRe = regexp.MustCompile(`(?U)\(.*\)\?`)
var cleanRe = regexp.MustCompile("[()$?]")
var reqdRe = regexp.MustCompile(`\(\?P<(\w+)>[^)]*\)`)

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
	BodyProps  []Property
}

type Property struct {
	Name        string
	Type        string
	Description string
}

type Parameter struct {
	Property Property
	In       string
}

type pathlet struct {
	pattern string
	params  map[string]bool
}

func deRegex(s string) string {
	return cleanRe.ReplaceAllString(s, "")
}

// expandPattern expands a regex pattern by generating permutations of any optional parameters
// and changing named parameters into their {open_api} style equivalents.
func expandPattern(root, pat string) []pathlet {
	paths := make([]pathlet, 0)
	toppaths := []string{pat}

	// expand all optional elements into two paths. This apporach really only useful up to 2 optional
	// groups, but we probably don't want to deal with the exponential increase for the general case.
	for i := 0; i < len(toppaths); i++ {
		p := toppaths[i]
		match := optRe.FindStringIndex(p)
		if match != nil {
			toppaths[i] = p[0:match[0]] + p[match[0]+1:match[1]-2] + p[match[1]:]
			toppaths = append(toppaths, p[0:match[0]]+p[match[1]:])
			i--
		}

	}

	sort.Strings(toppaths)

	for _, pat := range toppaths {
		params := make(map[string]bool)
		result := reqdRe.FindAllStringSubmatch(pat, -1)
		if result != nil {
			for _, p := range result {
				par := p[1]
				params[par] = true
				pat = strings.Replace(pat, p[0], fmt.Sprintf("{%s}", par), 1)
			}
		}
		pat = fmt.Sprintf("/%s/%s", root, pat)
		pat = deRegex(pat)
		paths = append(paths, pathlet{pat, params})
	}
	return paths
}

// TODO: this is conservative. Should omit surrounding quotes if not needed.
func escapeYAML(syn string) string {
	if idx := strings.Index(syn, "\n"); idx != -1 {
		syn = syn[0:idx] + "…"
	}
	syn = strings.Replace(syn, `"`, `\"`, -1)
	return fmt.Sprintf(`"%s"`, syn)
}

func procLogicalPath(p *framework.Path) []Path {
	var docPaths []Path
	var verbose bool

	if strings.Contains(p.Pattern, "revoke-prefix") {
		//verbose = true
	}
	paths := expandPattern("sys", p.Pattern)
	if verbose {
		fmt.Println(paths)
	}

	for _, path := range paths {
		methods := []Method{}
		for opType := range p.Callbacks {
			m := Method{
				Summary: "Yay, a summary!", // TODO escapify
			}
			switch opType {
			case logical.UpdateOperation:
				m.HTTPMethod = "post"
			case logical.DeleteOperation:
				m.HTTPMethod = "delete"
			case logical.ReadOperation:
				m.HTTPMethod = "get"
			case logical.ListOperation:
				continue
				//m.HTTPMethod = "get"
			default:
				panic(fmt.Sprintf("unknown operation type %v", opType))
			}

			d := make(map[string]bool)
			for param := range path.params {
				d[param] = true
				m.Parameters = append(m.Parameters, Parameter{
					In: "path",
					Property: Property{
						Name: param,
						Type: "string",
					},
				})
			}

			//m.BodyProps = make([]Property, 0)
			for name, field := range p.Fields {
				//fmt.Printf("Processing field %s\n", name)
				// TODO don't need ", ok"
				if _, ok := d[name]; !ok {
					m.BodyProps = append(m.BodyProps, Property{
						Name:        name,
						Description: escapeYAML(field.Description),
						Type:        "string", //field.Type.String(),
					})
				}
			}
			methods = append(methods, m)
		}
		if len(methods) > 0 {
			pd := Path{
				Pattern: path.pattern,
				Methods: methods,
			}
			docPaths = append(docPaths, pd)
		}
	}

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
		if !strings.Contains(p.Pattern, "revoke-prefix") {
			//continue
		}
		paths := procLogicalPath(p)
		doc.Paths = append(doc.Paths, paths...)
	}

	r := OAPIRenderer{
		output:   os.Stdout,
		template: tmpl,
	}
	_ = r
	r.render(doc)
}
