package main

const tmpl = `
openapi: "3.0.0"
info:
  version: {{ .Version }}
  title: Hashicorp
  license:
    name: Mozilla Public License 2.0

paths:{{ range .Paths }}
  {{ .Pattern }}:{{ range .Methods }}
    {{ .HTTPMethod }}:
      summary: {{ .Summary }}
      parameters: {{ range .Parameters }}
        - name: {{ .Name }}
          in: {{ .In }}
          type: {{ .Type }}{{ end }}
      responses:
        '200':
          description: Yay!{{ end }}
{{ end }}`
