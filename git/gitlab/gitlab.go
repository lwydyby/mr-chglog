package gitlab

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"

	"github.com/lwydyby/mr-chglog/git"
)

type Gitlab struct {
	client      *gitlab.Client
	projectName string
	projectID   int
}

func NewGit(token string, repositoryURL string, opts ...gitlab.ClientOptionFunc) git.GitClient {
	baseURL, nameWithNameSpace, projectName, err := getBaseUrlAndProjectName(repositoryURL)
	if err != nil {
		panic(err)
	}
	opts = append(opts, gitlab.WithBaseURL(fmt.Sprintf("https://%s/api/v4", baseURL)))

	client, err := gitlab.NewClient(token, opts...)
	if err != nil {
		panic(err)
	}
	projects, _, err := client.Search.Projects(projectName, &gitlab.SearchOptions{})
	if err != nil {
		panic(err)
	}
	if len(projects) == 0 {
		panic(fmt.Errorf("project %s not found", projectName))
	}
	var project *gitlab.Project
	for i := range projects {
		if projects[i].PathWithNamespace == nameWithNameSpace {
			project = projects[i]
			break
		}
	}
	return &Gitlab{
		client:      client,
		projectName: projectName,
		projectID:   project.ID,
	}
}

func (g *Gitlab) GetTags() []*git.Tag {
	tags, _, err := g.client.Tags.ListTags(g.projectID, &gitlab.ListTagsOptions{})
	if err != nil {
		log.Fatalf("Error fetching tags: %s", err)
	}

	var tagStructs []*git.Tag

	for _, tag := range tags {
		tagStruct := &git.Tag{
			Name:    tag.Name,
			Subject: tag.Commit.Title,
			Date:    *tag.Commit.CommittedDate,
		}
		tagStructs = append(tagStructs, tagStruct)
	}

	for i, tagStruct := range tagStructs {
		if i > 0 {
			tagStruct.Previous = tagStructs[i-1]
		}
		if i < len(tagStructs)-1 {
			tagStruct.Next = tagStructs[i+1]
		}
	}
	return tagStructs
}

func (g *Gitlab) GetMergeRequests(from, end *git.Tag) []*git.MergeRequest {
	var startTime, endTime *time.Time
	if from != nil {
		startCommit, _, err := g.client.Commits.GetCommit(g.projectID, from.Name)
		if err != nil {
			log.Fatalf("Error fetching start tag commit: %s", err)
		}
		startTime = g.getCommitDate(startCommit)
	}

	if end != nil {
		endCommit, _, err := g.client.Commits.GetCommit(g.projectID, end.Name)
		if err != nil {
			log.Fatalf("Error fetching end tag commit: %s", err)
		}
		endTime = g.getCommitDate(endCommit)
	}

	options := &gitlab.ListProjectMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		State: gitlab.String("merged"),
	}
	if startTime != nil {
		options.UpdatedAfter = startTime
	}
	if endTime != nil {
		before := endTime.Add(1 * time.Second)
		options.UpdatedBefore = &before
	}

	var allMRs []*gitlab.MergeRequest
	for {
		mrs, resp, err := g.client.MergeRequests.ListProjectMergeRequests(g.projectID, options)
		if err != nil {
			log.Fatalf("Error fetching merge requests: %s", err)
		}

		allMRs = append(allMRs, mrs...)
		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		options.Page = resp.NextPage
	}
	results := make([]*git.MergeRequest, len(allMRs), len(allMRs))
	for i := range allMRs {
		results[i] = &git.MergeRequest{
			ID:          allMRs[i].ID,
			IID:         allMRs[i].IID,
			Title:       allMRs[i].Title,
			MergedAt:    allMRs[i].MergedAt,
			Description: allMRs[i].Description,
			SHA:         allMRs[i].SHA[:8],
			WebURL:      allMRs[i].WebURL,
		}
		if allMRs[i].Author != nil {
			results[i].Author = allMRs[i].Author.Username
		}
	}
	return results
}

func (g *Gitlab) CreateTag(tag string, desc string) {
	// todo 支持指定ref
	ref := "master"
	_, _, err := g.client.Tags.CreateTag(g.projectID, &gitlab.CreateTagOptions{
		TagName:            &tag,
		Ref:                &ref,
		ReleaseDescription: &desc,
	})
	if err != nil {
		panic(err)
	}
}

func (g *Gitlab) GetMRChanges(mr *git.MergeRequest) {
	m, _, err := g.client.MergeRequests.GetMergeRequestChanges(g.projectID, mr.IID, &gitlab.GetMergeRequestChangesOptions{})
	if err != nil {
		panic(err)
	}
	diff := make([]git.Diff, 0)
	for j := range m.Changes {
		d := m.Changes[j]
		diff = append(diff, git.Diff{
			OldPath:     d.OldPath,
			NewPath:     d.NewPath,
			AMode:       d.AMode,
			BMode:       d.BMode,
			Diff:        d.Diff,
			NewFile:     d.NewFile,
			RenamedFile: d.RenamedFile,
			DeletedFile: d.DeletedFile,
		})
	}
	mr.Changes = diff
}

func (g *Gitlab) UpdateTagRelease(tagName string, desc string) {
	opt := &gitlab.CreateReleaseNoteOptions{
		Description: gitlab.String(desc),
	}
	_, _, err := g.client.Tags.CreateReleaseNote(g.projectID, tagName, opt)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Gitlab) getCommitDate(commit *gitlab.Commit) *time.Time {
	options := &gitlab.ListProjectMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
		State:          gitlab.String("merged"),
		AuthorUsername: &commit.AuthorName,
		Sort:           gitlab.String("desc"),
		OrderBy:        gitlab.String("created_at"),
	}
	for {
		mrs, resp, err := g.client.MergeRequests.ListProjectMergeRequests(g.projectID, options)
		if err != nil {
			log.Fatalf("Error fetching merge requests: %s", err)
		}
		for i := range mrs {
			commits, _, err := g.client.MergeRequests.GetMergeRequestCommits(g.projectID, mrs[i].IID, &gitlab.GetMergeRequestCommitsOptions{})
			if err != nil {
				log.Fatalf("Error fetching merge requests: %s", err)
			}
			if hasCommit(commits, commit) {
				return mrs[i].UpdatedAt
			}
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		options.Page = resp.NextPage
	}
	return commit.CommittedDate
}

func hasCommit(in []*gitlab.Commit, commit *gitlab.Commit) bool {
	for i := range in {
		if in[i].ID == commit.ID {
			return true
		}
	}
	return false
}

func getBaseUrlAndProjectName(urlStr string) (host string, nameWithNameSpace string, name string, err error) {
	if urlStr == "" {
		return "", "", "", errors.New("url is empty")
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", "", "", err
	}
	if u.Scheme == "" {
		return "", "", "", errors.New("url scheme is empty")
	}
	if u.Host == "" {
		return "", "", "", errors.New("url host is empty")
	}
	if u.Path == "" || u.Path == "/" {
		return "", "", "", errors.New("project name is empty")
	}
	projectName := path.Base(strings.TrimSuffix(u.Path, "/"))
	if projectName == "" {
		return "", "", "", errors.New("project name is empty")
	}
	return u.Host, u.Path[1:], projectName, nil
}
