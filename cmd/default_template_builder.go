package cmd

import "fmt"

type defaultTemplateBuilderImpl struct{}

func NewDefaultTemplateBuilder() TemplateBuilder {
	return &defaultTemplateBuilderImpl{}
}

// Build ...
func (t *defaultTemplateBuilderImpl) Build(ans *Answer) (string, error) {
	// unreleased
	tpl := t.unreleased()

	// version start
	tpl += "\n{{ if .Version }}\n"
	tpl += "\n{{ range .Version }}\n"

	tpl += t.versionHeader()

	// commits
	tpl += t.commits(".MRs")
	tpl += t.sql()
	tpl += t.setBreak()
	// versions end
	tpl += "\n{{ end -}}"
	tpl += "\n{{ end -}}"
	return tpl, nil
}

func (t *defaultTemplateBuilderImpl) unreleased() string {
	var (
		id      = ""
		title   = "Unreleased"
		commits = t.commits(".Unreleased")
	)

	title = fmt.Sprintf("[%s]", title)

	return fmt.Sprintf(`{{ if .Unreleased -}}
%s## %s

%s

%s
{{ end -}}
`, id, title, commits, t.sql())
}

func (t *defaultTemplateBuilderImpl) versionHeader() string {
	var (
		id      = ""
		tagName = "{{ .Tag.Name }}"
		date    = "({{ datetime \"2006-01-02\" .Tag.Date }})"
	)
	id = templateTagNameAnchor
	tagName = "{{ if .Tag.Previous }}[{{ .Tag.Name }} release note ]({{ $.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}){{ else }}{{ .Tag.Name }}{{ end }}"

	return fmt.Sprintf("%s## %s  %s\n", id, tagName, date)
}

func (t *defaultTemplateBuilderImpl) commits(commitGroups string) string {
	body := `
### {{ $prefix }}
{{- range $mergeRequests }}
  - {{ .Title }}([{{ .SHA }}]({{ .WebURL }}))@{{ .Author }}
{{- end}}
`
	return fmt.Sprintf(`{{- range $prefix, $mergeRequests := %s }}
%s
{{ end -}}
`, commitGroups, body)
}

func (t *defaultTemplateBuilderImpl) sql() string {
	body := `
{{ if .SQL }}
### SQL变更
%s
{{ .SQL }}
%s
{{ end }}
`
	return fmt.Sprintf(body, "```sql", "```")
}

func (t *defaultTemplateBuilderImpl) setBreak() string {
	body := `
{{ if .Break }}
### Break变更
{{ .Break }}

{{ end }}
`
	return body
}
