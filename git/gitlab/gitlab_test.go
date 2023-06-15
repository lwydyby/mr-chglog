package gitlab

import (
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"

	"github.com/lwydyby/mr-chglog/git"
)

func TestNewGit(t *testing.T) {
	mockey.PatchConvey("new_git", t, func() {
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
		client := NewGit("123", "http://github.com/xxxx/sxsxs")
		assert.Equal(t, 123, client.(*Gitlab).projectID)
	})
}

func TestGetTag(t *testing.T) {
	mockey.PatchConvey("get_tag", t, func() {
		c := &gitlab.Client{
			Search: &gitlab.SearchService{},
			Tags:   &gitlab.TagsService{},
		}
		mockey.Mock(gitlab.NewClient).Return(c, nil).Build()
		mockey.Mock(mockey.GetMethod(c.Search, "Projects")).Return([]*gitlab.Project{
			{
				PathWithNamespace: "xxxx/sxsxs",
				ID:                123,
			},
		}, nil, nil).Build()
		now := time.Now()
		mockey.Mock(mockey.GetMethod(c.Tags, "ListTags")).Return([]*gitlab.Tag{
			{
				Name: "123",
				Commit: &gitlab.Commit{
					Title:         "123",
					CommittedDate: &now,
				},
			},
		}, nil, nil).Build()
		client := NewGit("123", "http://github.com/xxxx/sxsxs")
		assert.Equal(t, &git.Tag{
			Name:    "123",
			Subject: "123",
			Date:    now,
		}, client.GetTags()[0])
	})
}

func TestGetMergeRequests(t *testing.T) {
	mockey.PatchConvey("get_mq", t, func() {
		c := &gitlab.Client{
			Search:        &gitlab.SearchService{},
			Tags:          &gitlab.TagsService{},
			Commits:       &gitlab.CommitsService{},
			MergeRequests: &gitlab.MergeRequestsService{},
		}
		mockey.Mock(gitlab.NewClient).Return(c, nil).Build()
		mockey.Mock(mockey.GetMethod(c.Search, "Projects")).Return([]*gitlab.Project{
			{
				PathWithNamespace: "xxxx/sxsxs",
				ID:                123,
			},
		}, nil, nil).Build()
		mockey.Mock(mockey.GetMethod(c.MergeRequests, "ListProjectMergeRequests")).Return([]*gitlab.MergeRequest{}, &gitlab.Response{
			CurrentPage: 1,
			TotalPages:  1,
		}, nil).Build()
		client := NewGit("123", "http://github.com/xxxx/sxsxs")
		assert.Equal(t, 0, len(client.GetMergeRequests(nil, nil)))
	})
}

func TestGetCreateTag(t *testing.T) {
	mockey.PatchConvey("create_tag", t, func() {
		c := &gitlab.Client{
			Search: &gitlab.SearchService{},
			Tags:   &gitlab.TagsService{},
		}
		mockey.Mock(gitlab.NewClient).Return(c, nil).Build()
		mockey.Mock(mockey.GetMethod(c.Search, "Projects")).Return([]*gitlab.Project{
			{
				PathWithNamespace: "xxxx/sxsxs",
				ID:                123,
			},
		}, nil, nil).Build()
		mockey.Mock(mockey.GetMethod(c.Tags, "CreateTag")).Return(nil, &gitlab.Response{}, nil).Build()
		client := NewGit("123", "http://github.com/xxxx/sxsxs")
		client.CreateTag("", "")
	})
}

func TestGetMRChanges(t *testing.T) {
	mockey.PatchConvey("get_mr_change", t, func() {
		c := &gitlab.Client{
			Search:        &gitlab.SearchService{},
			Tags:          &gitlab.TagsService{},
			Commits:       &gitlab.CommitsService{},
			MergeRequests: &gitlab.MergeRequestsService{},
		}
		mockey.Mock(gitlab.NewClient).Return(c, nil).Build()
		mockey.Mock(mockey.GetMethod(c.Search, "Projects")).Return([]*gitlab.Project{
			{
				PathWithNamespace: "xxxx/sxsxs",
				ID:                123,
			},
		}, nil, nil).Build()
		mockey.Mock(mockey.GetMethod(c.MergeRequests, "GetMergeRequestChanges")).Return(&gitlab.MergeRequest{}, &gitlab.Response{
			CurrentPage: 1,
			TotalPages:  1,
		}, nil).Build()
		client := NewGit("123", "http://github.com/xxxx/sxsxs")
		client.GetMRChanges(&git.MergeRequest{})
	})
}
