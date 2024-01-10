package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const md = `
# Title

Some content

# break

- test1
- test2

# sql

%s

# Another Title

More content under SQL heading.
Other content.


`

func TestGetChangeSQL(t *testing.T) {
	sql := "--查询用户\n" + "SELECT * FROM users;\n"
	mr := &MergeRequest{
		Description: fmt.Sprintf(md, fmt.Sprintf("```sql\n%s```\n", sql)),
	}
	assert.Equal(t, sql, GetHead([]*MergeRequest{mr}, "sql"))
}

func TestBreakChange(t *testing.T) {
	mr := &MergeRequest{
		Description: fmt.Sprintf(md, fmt.Sprintf("```sql\n%s```\n", "sql")),
	}
	br := mr.GetMarkdownInfo("break")
	assert.Equal(t, "- test1\n- test2\n", br)
}

func TestGroupByPrefix(t *testing.T) {
	mrs := []*MergeRequest{
		{
			Title: "feat: test",
		},
		{
			Title: "fix: bug",
		},
	}
	want := map[string][]*MergeRequest{
		"Feat": {
			{
				Title: "test",
			},
		},
		"Fix": {
			{
				Title: "bug",
			},
		},
	}
	assert.Equal(t, want, GroupByPrefix(mrs))
}
