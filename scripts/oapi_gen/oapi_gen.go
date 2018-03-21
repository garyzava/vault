package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

type PathDoc struct {
	Pattern string
	Method  string
}

func main() {
	c := vault.Core{}
	b := vault.NewSystemBackend(&c)

	paths := []PathDoc{}

	for _, path := range b.Backend.Paths {
		method := "GET"
		if f, ok := path.Callbacks[logical.DeleteOperation]; ok {
			fmt.Println(f)
			method = "DELETE"
		}

		pd := PathDoc{
			Pattern: path.Pattern,
			Method:  method,
		}
		paths = append(paths, pd)
	}

	//const prefix = "/sys/"
	//sort.Strings(paths)
	//for _, path := range paths {
	//	fmt.Printf("%s%s\n", prefix, path)
	//}
	tmpl, _ := template.New("test").Parse(tmpl)
	tmpl.Execute(os.Stdout, paths)

}
