package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultTemplateBuilderImpl_Build(t1 *testing.T) {
	type args struct {
		ans *Answer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple",
			args: args{
				ans: &Answer{},
			},
			want: `{{ if .Unreleased -}}
## [Unreleased]

{{- range $prefix, $mergeRequests := .Unreleased }}

### {{ $prefix }}
{{- range $mergeRequests }}
  - {{ .Title }}([{{ .SHA }}]({{ .WebURL }}))@{{ .Author }}
{{- end}}

{{ end -}}



{{ if .SQL }}
### SQL变更
` + "```sql\n" +
				"{{ .SQL }}\n" +
				"```" + `
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
` + "```sql\n" +
				"{{ .SQL }}\n" +
				"```" + `
{{ end }}

{{ if .Break }}
### Break变更
{{ .Break }}

{{ end }}

{{ end -}}
{{ end -}}`,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &defaultTemplateBuilderImpl{}
			got, _ := t.Build(tt.args.ans)
			assert.Equalf(t1, tt.want, got, "Build(%v)", tt.args.ans)
		})
	}
}
