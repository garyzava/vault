package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"text/template"
)

type OAPIRenderer struct {
	output   io.Writer
	template string
}

func (r *OAPIRenderer) render(doc Doc) {
	funcs := map[string]interface{}{
		"indent": funcIndent,
	}

	tmpl, _ := template.New("root").Funcs(funcs).Parse(r.template)
	tmpl.Execute(r.output, doc)
}

func funcIndent(count int, text string) string {
	var buf bytes.Buffer
	prefix := strings.Repeat(" ", count)
	scan := bufio.NewScanner(strings.NewReader(text))
	for scan.Scan() {
		buf.WriteString(prefix + scan.Text() + "\n")
	}

	return strings.TrimRight(buf.String(), "\n")
}
