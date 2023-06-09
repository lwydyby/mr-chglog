package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const md = `
# Title

Some content

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
