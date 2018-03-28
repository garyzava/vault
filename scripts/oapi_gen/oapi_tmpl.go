package main

const tmpl = `swagger: "2.0"
info:
  title: HashiCorp Vault API
  version: {{ .Version }}
  contact:
    name: HashiCorp Vault
    url: https://www.vaultproject.io/api/index.html
  license:
    name: Mozilla Public License 2.0
    url: https://www.mozilla.org/en-US/MPL/2.0

paths:{{ range .Paths }}
  {{ .Pattern }}:{{ range .Methods }}
    {{ .HTTPMethod }}:{{ if .Summary }}
      summary: {{ .Summary }}{{ end }}
      produces:
        - application/json
      tags:
	  {{- range .Tags }}
        - {{ . }}
      {{- end }}{{ if (or .Parameters .BodyProps) }}
      parameters:{{ range .Parameters }}
        - name: {{ .Property.Name }}
          description: {{ .Property.Description }}
          in: {{ .In }}
          type: {{ .Property.Type }}
          required: true{{ end -}}
      {{ end -}}
      {{ if .BodyProps }}
        - name: Data
          in: body
          schema:
            type: object
            properties:{{ range .BodyProps }}
              {{ .Name }}:
                description: {{ .Description }}
                type: {{ .Type }}{{ end }}{{ end }}
      responses:
        200:
          description: OK{{ end }}
{{ end }}`
