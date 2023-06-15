package cmd

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/tsuyoshiwada/go-gitcmd"
)

func TestAsk(t *testing.T) {
	mockey.PatchConvey("ask", t, func() {
		c := gitcmd.New(nil)
		fs := &mockFileSystem{}
		q := NewQuestioner(c, fs)
		mockey.Mock(survey.Ask).Return(nil).Build()
		mockey.Mock(survey.AskOne).To(func(p survey.Prompt, res interface{}, opts ...survey.AskOpt) error {
			if res != nil {
				b, ok := res.(*bool)
				if ok {
					*b = true
				}
			}
			return nil
		}).Build()
		mockey.Mock(mockey.GetMethod(fs, "Exists")).Return(true).Build()
		ans, err := q.Ask()
		assert.Nil(t, err)
		assert.Equal(t, &Answer{
			ChatID:     []string{""},
			NeedRobot:  true,
			SkipTpl:    false,
			SkipConfig: false,
		}, ans)
	})
}
