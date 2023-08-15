package cmd

import (
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/lwydyby/mr-chglog/config"
)

func TestCLIRun(t *testing.T) {
	mockey.PatchConvey("cli_run", t, func() {
		mockFS := &mockFileSystem{}
		cli := NewCLI(&CLIContext{
			AppID:         "123",
			AppSecret:     "12213",
			ChatID:        []string{"12312"},
			BotTitle:      "asdqw",
			Template:      "/etc",
			RepositoryURL: "http://github.com/123",
			Token:         "12213123",
		}, mockFS, NewConfigLoader(), NewGenerator())
		mockey.Mock(mockey.GetMethod(cli.configLoader, "Load")).Return(&config.MRChLogConfig{}, nil).Build()
		mockey.Mock(mockey.GetMethod(cli.generator, "Generate")).Return(nil).Build()
		assert.Equal(t, 1, cli.Run())
	})
}
