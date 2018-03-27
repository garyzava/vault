package main

const tmpl = `openapi: "3.0.0"
info:
  version: {{ .Version }}
  title: HashiCorp
  license:
    name: Mozilla Public License 2.0

paths:{{ range .Paths }}
  {{ .Pattern }}:{{ range .Methods }}
    {{ .HTTPMethod }}:{{ if .Summary }}
      summary: {{ .Summary }}{{ end }}{{ if .Parameters }}
      parameters: {{ range .Parameters }}
        - name: {{ .Property.Name }}
          description: {{ .Property.Description }}
          in: {{ .In }}
          type: {{ .Property.Type }}{{ end }}
      {{ end }}
      responses:
        '200':
          description: Yay!{{ end }}
{{ end }}`
