package cmd

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"

	"github.com/lwydyby/mr-chglog/bot"
	"github.com/lwydyby/mr-chglog/config"
	"github.com/lwydyby/mr-chglog/git"
	"github.com/lwydyby/mr-chglog/git/gitlab"
)

// Generator ...
type Generator interface {
	Generate(io.Writer, *CLIContext, *config.MRChLogConfig) error
}

type generatorImpl struct {
	client git.GitClient
	config *config.MRChLogConfig
}

// NewGenerator ...
func NewGenerator() Generator {
	return &generatorImpl{}
}

// Generate ...
func (g *generatorImpl) Generate(w io.Writer, ctx *CLIContext, c *config.MRChLogConfig) error {
	g.client = gitlab.NewGit(c.Token, c.RepositoryURL)
	g.config = c
	allTags := g.client.GetTags()
	allTags, t, err := git.Select(allTags, ctx.Query)
	if err != nil {
		return err
	}
	versions := g.getMRGroup(allTags, t, ctx, c)
	b := &bytes.Buffer{}

	if len(ctx.NextTag) != 0 || ctx.PushBot || ctx.Update {
		w = b
	}
	if len(ctx.NextTag) != 0 {
		versions[len(versions)-1].Tag = &git.Tag{
			Name: ctx.NextTag,
			Date: time.Now(),
		}
		if len(versions) > 2 {
			versions[len(versions)-1].Tag.Previous = versions[len(versions)-2].Tag
		}
		versions = versions[len(versions)-1:]
		defer func() {
			g.client.CreateTag(ctx.NextTag, b.String())
		}()
	}
	if ctx.Update && ctx.IsSingleTag() {
		defer func() {
			g.client.UpdateTagRelease(versions[0].Tag.Name, b.String())
		}()
	}
	if ctx.PushBot {
		defer func() {
			if ctx.PushBot && len(c.ChatID) != 0 {
				token, err := bot.GetTenantAccessToken(context.Background(), c)
				if err != nil {
					panic(err)
				}
				for i := range c.ChatID {
					err = bot.SendAlertMessage(context.Background(), token, c.ChatID[i], c.BotTitle, b.String())
					if err != nil {
						panic(err)
					}
				}
			}
		}()
	}
	return g.render(w, versions)
}

func (g *generatorImpl) getMRGroup(tags []*git.Tag, from *git.Tag, ctx *CLIContext, c *config.MRChLogConfig) []*Version {
	if len(tags) == 0 {
		mr := g.client.GetMergeRequests(nil, nil)
		if ctx.AI {
			ac := commitBuilderFactory(ctx.AIType, c.POEToken)
			for j := range mr {
				g.client.GetMRChanges(mr[j])
				resp := ac.BuildCommit(context.Background(), mr[j])
				mr[j].Title = mr[j].Title[:strings.Index(mr[j].Title, ":")+1] + resp
			}
		}
		return []*Version{{Tag: nil, MRs: git.GroupByPrefix(mr), SQL: git.GetSQL(mr), Break: git.GetHead(mr, "break")}}
	}

	results := make([]*Version, 0, len(tags))

	var prevTag = from
	for i := len(tags) - 1; i >= 0; i-- {
		tag := tags[i]

		mr := g.client.GetMergeRequests(prevTag, tag)
		if ctx.AI {
			ac := commitBuilderFactory(ctx.AIType, c.POEToken)
			for j := range mr {
				g.client.GetMRChanges(mr[j])
				resp := ac.BuildCommit(context.Background(), mr[j])
				mr[j].Title = mr[j].Title[:strings.Index(mr[j].Title, ":")+1] + resp
			}
		}
		results = append(results, &Version{Tag: tag, MRs: git.GroupByPrefix(mr), SQL: git.GetSQL(mr), Break: git.GetHead(mr, "break")})

		prevTag = tag
	}

	if from == nil && !ctx.IsSingleTag() {
		mr := g.client.GetMergeRequests(prevTag, nil)
		results = append(results, &Version{MRs: git.GroupByPrefix(mr), SQL: git.GetSQL(mr), Break: git.GetHead(mr, "break")})
	}

	return results
}

type Version struct {
	Tag   *git.Tag
	MRs   map[string][]*git.MergeRequest
	SQL   string
	Break string
}

func (g *generatorImpl) render(w io.Writer, versions []*Version) error {
	if _, err := os.Stat(g.config.Template); err != nil {
		return err
	}

	fmap := template.FuncMap{
		// format the input time according to layout
		"datetime": func(layout string, input time.Time) string {
			return input.Format(layout)
		},
		// upper case the first character of a string
		"upperFirst": func(s string) string {
			if len(s) > 0 {
				return strings.ToUpper(string(s[0])) + s[1:]
			}
			return ""
		},
		// indent all lines of s n spaces
		"indent": func(s string, n int) string {
			if len(s) == 0 {
				return ""
			}
			pad := strings.Repeat(" ", n)
			return pad + strings.ReplaceAll(s, "\n", "\n"+pad)
		},
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"replace":   strings.Replace,
	}

	fname := filepath.Base(g.config.Template)

	t := template.Must(template.New(fname).Funcs(sprig.TxtFuncMap()).Funcs(fmap).ParseFiles(g.config.Template))
	var unreleased map[string][]*git.MergeRequest
	if versions[len(versions)-1].Tag == nil {
		unreleased = versions[len(versions)-1].MRs
		versions = versions[:len(versions)-1]
	}
	return t.Execute(w, &RenderData{
		Title:         g.config.Title,
		RepositoryURL: g.config.RepositoryURL,
		Unreleased:    unreleased,
		Version:       versions,
	})
}

type RenderData struct {
	Title         string
	RepositoryURL string
	Unreleased    map[string][]*git.MergeRequest
	Version       []*Version
}
