{{ if .Unreleased -}}
## [Unreleased]

{{- range $prefix, $mergeRequests := .Unreleased }}

### {{ $prefix }}
{{- range $mergeRequests }}
  - {{ .Title }}([{{ .SHA }}]({{ .WebURL }}))@{{ .Author }}
{{- end}}

{{ end -}}



{{ if .SQL }}
### SQL变更
```sql
{{ .SQL }}
```
{{ end }}

{{ end -}}

{{ if .Version }}

{{ range .Version }}
## {{ if .Tag.Previous }}[{{ .Tag.Name }} release note ]({{ $.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}){{ else }}{{ .Tag.Name }}{{ end }}  ({{ datetime "2006-01-02" .Tag.Date }})
{{- range $prefix, $mergeRequests := .MRs }}

### {{ $prefix }}
{{- range $mergeRequests }}
  - {{ .Title }}([{{ .SHA }}]({{ .WebURL }}))@{{ .Author }}
{{- end}}

{{ end -}}

{{ if .SQL }}
### SQL变更
```sql
{{ .SQL }}
```
{{ end }}

{{ end -}}
{{ end -}}