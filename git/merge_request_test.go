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
	sql := "--查询用户" + "SELECT * FROM users;\n"
	mr := &MergeRequest{
		Description: fmt.Sprintf(md, fmt.Sprintf("```sql\n%s```\n", sql)),
	}
	assert.Equal(t, sql, mr.GetChangeSQL())
}

func TestBreakChange(t *testing.T) {
	mr := &MergeRequest{
		Description: fmt.Sprintf(md, fmt.Sprintf("```sql\n%s```\n", "sql")),
	}
	br := mr.GetHeadChange("break")
	assert.Equal(t, "\n- test1\n- test2\n", br)
}
