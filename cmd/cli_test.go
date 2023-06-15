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
		cli := NewCLI(&CLIContext{}, mockFS, NewConfigLoader(), NewGenerator())
		mockey.Mock(mockey.GetMethod(cli.configLoader, "Load")).Return(&config.MRChLogConfig{}, nil).Build()
		mockey.Mock(mockey.GetMethod(cli.generator, "Generate")).Return(nil).Build()
		assert.Equal(t, 1, cli.Run())
	})
}
