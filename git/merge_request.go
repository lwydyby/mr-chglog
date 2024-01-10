package git

import (
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

func (m *MergeRequest) GetMarkdownInfo(head string) string {
	gm := goldmark.New(
		goldmark.WithExtensions(),
	)
	p := gm.Parser()
	reader := text.NewReader([]byte(m.Description))
	doc := p.Parse(reader)

	contentBuilder := &strings.Builder{}

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch node := n.(type) {
			case *ast.Heading:
				return walkHead(head, node, contentBuilder, reader), nil
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		panic(err)
	}
	return contentBuilder.String()
}

func walkHead(title string, head *ast.Heading, build *strings.Builder, reader text.Reader) ast.WalkStatus {
	headingText := string(head.Text(reader.Source()))
	if headingText != title {
		return ast.WalkContinue
	}
	for child := head.NextSibling(); child != nil && child.Kind() != ast.KindHeading; child = child.NextSibling() {
		if breakList, ok := child.(*ast.List); ok {
			for child := breakList.FirstChild(); child != nil; child = child.NextSibling() {
				if listItem, ok := child.(*ast.ListItem); ok {
					segment := listItem.Text(reader.Source())
					build.WriteByte(breakList.Marker)
					build.WriteByte(' ')
					build.Write(segment)
					build.WriteString("\n")
				}
			}
			return ast.WalkStop
		}
	}
	for child := head.NextSibling(); child != nil && child.Kind() != ast.KindHeading; child = child.NextSibling() {
		if codeBlock, ok := child.(*ast.FencedCodeBlock); ok {
			lines := codeBlock.Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				build.Write(line.Value(reader.Source()))
			}
			return ast.WalkStop
		}
	}
	return ast.WalkContinue
}

func GetHead(mrs []*MergeRequest, head string) string {
	var contentBuilder strings.Builder
	for _, mr := range mrs {
		content := mr.GetMarkdownInfo(head)
		if content != "" && !strings.HasSuffix(content, "\n") {
			content = content + "\n"
		}
		if content != "" {
			contentBuilder.WriteString(content)
		}
	}
	return contentBuilder.String()
}

func GroupByPrefix(mrs []*MergeRequest) map[string][]*MergeRequest {
	grouped := make(map[string][]*MergeRequest)

	for _, mr := range mrs {
		// 需要将skipped从title中移除
		mr.Title = strings.ReplaceAll(mr.Title, "[skipped]", "")
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
