package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"

	"github.com/lwydyby/mr-chglog/bot"
	"github.com/lwydyby/mr-chglog/config"
	"github.com/lwydyby/mr-chglog/git"
	m_gitlab "github.com/lwydyby/mr-chglog/git/gitlab"
)

func TestGenerate(t *testing.T) {
	mockey.PatchConvey("generate", t, func() {
		c := &gitlab.Client{
			Search: &gitlab.SearchService{},
		}
		mockey.Mock(gitlab.NewClient).Return(c, nil).Build()
		mockey.Mock(mockey.GetMethod(c.Search, "Projects")).Return([]*gitlab.Project{
			{
				PathWithNamespace: "xxxx/sxsxs",
				ID:                123,
			},
		}, nil, nil).Build()
		client := m_gitlab.NewGit("123", "http://github.com/xxxx/sxsxs")
		g := NewGenerator()
		temp := &template.Template{}
		mockey.Mock(m_gitlab.NewGit).Return(client).Build()
		mockey.Mock(mockey.GetMethod(client, "GetTags")).Return([]*git.Tag{
			{
				Name: "v0.0.1",
			},
		}).Build()
		mockey.Mock(mockey.GetMethod(client, "GetMergeRequests")).Return([]*git.MergeRequest{
			{
				Title: "feat: test",
			},
		}).Build()
		mockey.Mock(mockey.GetMethod(client, "CreateTag")).Return().Build()
		mockey.Mock(mockey.GetMethod(client, "GetMRChanges")).Return().Build()
		aiBuilder := &AICommitBuilder{}
		mockey.Mock(commitBuilderFactory).Return(aiBuilder).Build()
		mockey.Mock(mockey.GetMethod(aiBuilder, "BuildCommit")).Return("test").Build()
		mockey.Mock(os.Stat).Return(nil, nil).Build()
		mockey.Mock(filepath.Base).Return("").Build()
		mockey.Mock(template.Must).Return(temp).Build()
		mockey.Mock(template.New).Return(temp).Build()
		mockey.Mock(bot.GetTenantAccessToken).Return("123", nil).Build()
		mockey.Mock(bot.SendAlertMessage).Return(nil).Build()
		mockey.Mock(mockey.GetMethod(temp, "Funcs")).Return(temp).Build()
		mockey.Mock(mockey.GetMethod(temp, "ParseFiles")).Return(temp, nil).Build()
		mockey.Mock(mockey.GetMethod(&template.Template{}, "Execute")).Return(nil).Build()
		assert.Nil(t, g.Generate(&bytes.Buffer{}, &CLIContext{
			Token:         "123",
			RepositoryURL: "http://github.com/xxxx/sxsxs",
			NextTag:       "v0.0.2",
			PushBot:       true,
		}, &config.MRChLogConfig{
			ChatID: []string{"123"},
		}))
	})
}
