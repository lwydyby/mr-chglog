package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/tsuyoshiwada/go-gitcmd"
)

// Answer ...
type Answer struct {
	RepositoryURL string `survey:"repository_url"`
	Style         string `survey:"style"`
	Template      string `survey:"template"`
	ConfigDir     string `survey:"config_dir"`
	Token         string `survey:"token"`
	POEToken      string `survey:"poe_token"`
}

// Questioner ...
type Questioner interface {
	Ask() (*Answer, error)
}

type questionerImpl struct {
	client gitcmd.Client
	fs     FileSystem
}

// NewQuestioner ...
func NewQuestioner(client gitcmd.Client, fs FileSystem) Questioner {
	return &questionerImpl{
		client: client,
		fs:     fs,
	}
}

// Ask ...
func (q *questionerImpl) Ask() (*Answer, error) {
	ans, err := q.ask()
	if err != nil {
		return nil, err
	}

	config := filepath.Join(ans.ConfigDir, defaultConfigFilename)
	tpl := filepath.Join(ans.ConfigDir, defaultTemplateFilename)
	c := q.fs.Exists(config)
	t := q.fs.Exists(tpl)
	msg := ""

	switch {
	case c && t:
		msg = fmt.Sprintf("\"%s\" and \"%s\" already exists. Do you want to overwrite?", config, tpl)
	case c:
		msg = fmt.Sprintf("\"%s\" already exists. Do you want to overwrite?", config)
	case t:
		msg = fmt.Sprintf("\"%s\" already exists. Do you want to overwrite?", tpl)
	}

	if msg != "" {
		overwrite := false
		err = survey.AskOne(&survey.Confirm{
			Message: msg,
			Default: true,
		}, &overwrite, nil)

		if err != nil || !overwrite {
			return nil, errors.New("creation of the file was interrupted")
		}
	}

	return ans, nil
}

func (q *questionerImpl) ask() (*Answer, error) {
	ans := &Answer{}
	// todo 自定义模板
	tpls := q.getPreviewableList(templates)

	var previewableTransform = func(ans interface{}) (newAns interface{}) {
		if s, ok := ans.(survey.OptionAnswer); ok {
			newAns = survey.OptionAnswer{
				Value: q.parsePreviewableList(s.Value),
				Index: s.Index,
			}
		}
		return
	}

	questions := []*survey.Question{
		{
			Name: "repository_url",
			Prompt: &survey.Input{
				Message: "What is the URL of your repository?",
				Default: q.getRepositoryURL(),
			},
		},
		{
			Name: "style",
			Prompt: &survey.Select{
				Message: "What is your favorite style?",
				Options: styles,
				Default: styles[0],
			},
		},
		{
			Name: "template",
			Prompt: &survey.Select{
				Message: "What is your favorite template style?",
				Options: tpls,
				Default: tpls[0],
			},
			Transform: previewableTransform,
		},
		{
			Name: "config_dir",
			Prompt: &survey.Input{
				Message: "In which directory do you output configuration files and templates?",
				Default: defaultConfigDir,
			},
		},
		{
			Name: "token",
			Prompt: &survey.Input{
				Message: "what is your gitlab token?",
			},
		},
		{
			Name: "poe_token",
			Prompt: &survey.Input{
				Message: "what is your poe token?",
			},
		},
	}

	err := survey.Ask(questions, ans)
	if err != nil {
		return nil, err
	}

	return ans, nil
}

func (*questionerImpl) getPreviewableList(list []Previewable) []string {
	arr := make([]string, len(list))
	max := 0

	for _, p := range list {
		l := len(p.Display())
		if max < l {
			max = l
		}
	}

	for i, p := range list {
		arr[i] = fmt.Sprintf(
			"%s -- %s",
			p.Display()+strings.Repeat(" ", max-len(p.Display())),
			p.Preview(),
		)
	}

	return arr
}

func (*questionerImpl) parsePreviewableList(input string) string {
	return strings.TrimSpace(strings.Split(input, "--")[0])
}

func (q *questionerImpl) getRepositoryURL() string {
	if q.client.CanExec() != nil || q.client.InsideWorkTree() != nil {
		return ""
	}

	rawurl, err := q.client.Exec("config", "--get", "remote.origin.url")
	if err != nil {
		return ""
	}

	return remoteOriginURLToHTTP(rawurl)
}
