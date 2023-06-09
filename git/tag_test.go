package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tags = []*Tag{{
	Name: "v0.1.0",
}, {
	Name: "v0.2.0",
}, {
	Name: "v0.3.0",
}, {
	Name: "v0.4.0",
}, {
	Name: "v0.5.0",
}, {
	Name: "v0.6.0",
}}

func TestSelect(t *testing.T) {
	args := []struct {
		query string
		want  []*Tag
	}{
		{
			query: "",
			want: []*Tag{{
				Name: "v0.1.0",
			}, {
				Name: "v0.2.0",
			}, {
				Name: "v0.3.0",
			}, {
				Name: "v0.4.0",
			}, {
				Name: "v0.5.0",
			}, {
				Name: "v0.6.0",
			}},
		},
		{
			query: "..v0.4.0",
			want: []*Tag{{
				Name: "v0.1.0",
			}, {
				Name: "v0.2.0",
			}, {
				Name: "v0.3.0",
			}, {
				Name: "v0.4.0",
			}},
		},
		{
			query: "v0.2.0..v0.4.0",
			want: []*Tag{{
				Name: "v0.2.0",
			}, {
				Name: "v0.3.0",
			}, {
				Name: "v0.4.0",
			}},
		},
		{
			query: "v0.2.0..",
			want: []*Tag{{
				Name: "v0.2.0",
			}, {
				Name: "v0.3.0",
			}, {
				Name: "v0.4.0",
			}, {
				Name: "v0.5.0",
			}, {
				Name: "v0.6.0",
			}},
		},
		{
			query: "v0.2.0",
			want: []*Tag{{
				Name: "v0.2.0",
			}},
		},
	}
	t.Parallel()
	for _, tt := range args {
		tt := tt
		t.Run("", func(t *testing.T) {
			target, _, err := Select(tags, tt.query)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, target)
		})
	}
}
