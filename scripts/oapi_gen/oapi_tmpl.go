package main

const tmpl = `swagger: "2.0"
info:
  version: {{ .Version }}
  title: HashiCorp
  license:
    name: Mozilla Public License 2.0

paths:{{ range .Paths }}
  {{ .Pattern }}:{{ range .Methods }}
    {{ .HTTPMethod }}:{{ if .Summary }}
      summary: {{ .Summary }}{{ end }}{{ if (or .Parameters .BodyProps) }}
      parameters: {{ range .Parameters }}
        - name: {{ .Property.Name }}
          description: {{ .Property.Description }} BLAH
          in: {{ .In }}
          type: {{ .Property.Type }}
          required: true{{ end -}}
      {{ end -}}
      {{ if .BodyProps }}
        - name: Delete me
          in: body
          schema:
            type: object
            properties:{{ range .BodyProps }}
              {{ .Name }}:
                description: {{ .Description }}
                type: {{ .Type }}{{ end }}{{ end }}
      responses:
        '200':
          description: Yay!{{ end }}
{{ end }}`

//summary: {{ .Summary }}{{ end }}{{ if or((ge (len .Parameters) 0) (ge (len .BodyProperties) 0)) }}
