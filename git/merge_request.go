package git

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type MergeRequest struct {
	ID          int
	IID         int
	Title       string     `json:"title"`
	MergedAt    *time.Time `json:"merged_at"`
	Description string     `json:"description"`
	Changes     []Diff
	Author      string `json:"author"`
	SHA         string `json:"sha"`
	WebURL      string `json:"web_url"`
}

type Diff struct {
	OldPath     string `json:"old_path"`
	NewPath     string `json:"new_path"`
	AMode       string `json:"a_mode"`
	BMode       string `json:"b_mode"`
	Diff        string `json:"diff"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
}

func (d Diff) String() string {
	return d.Diff
}

func (m *MergeRequest) GetChangeSQL() string {
	gm := goldmark.New(
		goldmark.WithExtensions(),
	)
	p := gm.Parser()
	reader := text.NewReader([]byte(m.Description))
	doc := p.Parse(reader)

	var sqlContentBuilder strings.Builder

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch node := n.(type) {
			case *ast.Heading:
				headingText := string(node.Text(reader.Source()))

				if headingText == "sql" {
					for child := node.NextSibling(); child != nil && child.Kind() != ast.KindHeading; child = child.NextSibling() {
						if codeBlock, ok := child.(*ast.FencedCodeBlock); ok {
							lines := codeBlock.Lines()
							for i := 0; i < lines.Len(); i++ {
								line := lines.At(i)
								sqlContentBuilder.Write(line.Value(reader.Source()))
							}
							return ast.WalkStop, nil
						}
					}
				}
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		panic(err)
	}
	return sqlContentBuilder.String()
}

func (m *MergeRequest) GetHeadChange(title string) string {
	return getTextUnderHeading(m.Description, title)
}

func getTextUnderHeading(markdown, heading string) string {
	lines := strings.Split(markdown, "\n")
	start := -1
	end := -1

	for i, line := range lines {
		if strings.TrimSpace(line) == fmt.Sprintf("# %s", heading) {
			start = i
			continue
		}
		if start != -1 && end == -1 && strings.HasPrefix(strings.TrimSpace(line), "#") {
			end = i
			break
		}
	}

	if start == -1 {
		return ""
	}

	if end == -1 {
		return strings.Join(lines[start+1:], "\n")
	}

	return strings.Join(lines[start+1:end], "\n")
}

func GetSQL(mrs []*MergeRequest) string {
	var sqlContentBuilder strings.Builder
	for _, mr := range mrs {
		sql := mr.GetChangeSQL()
		if sql != "" {
			sqlContentBuilder.Write([]byte(sql))
		}
	}
	return sqlContentBuilder.String()
}

func GetHead(mrs []*MergeRequest, head string) string {
	var sqlContentBuilder strings.Builder
	for _, mr := range mrs {
		sql := mr.GetHeadChange(head)
		if !strings.HasSuffix(sql, "\n") {
			sql = sql + "\n"
		}
		if sql != "" {
			sqlContentBuilder.Write([]byte(sql))
		}
	}
	return sqlContentBuilder.String()
}

func GroupByPrefix(mrs []*MergeRequest) map[string][]*MergeRequest {
	grouped := make(map[string][]*MergeRequest)

	for _, mr := range mrs {
		parts := strings.SplitN(mr.Title, ":", 2)
		prefix := "Other"
		if len(parts) > 1 {
			prefix = strings.TrimSpace(parts[0])
			mr.Title = strings.TrimSpace(parts[1])

			// 将前缀的首字母大写
			for i, r := range prefix {
				prefix = string(unicode.ToTitle(r)) + prefix[i+1:]
				break
			}
		}

		grouped[prefix] = append(grouped[prefix], mr)
	}

	return grouped
}
