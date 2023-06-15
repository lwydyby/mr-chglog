package cmd

import (
	"bytes"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestInitializer(t *testing.T) {
	mockey.PatchConvey("initializer", t, func() {
		assert := assert.New(t)

		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		mockFs := &mockFileSystem{
			ReturnMkdirP: func(path string) error {
				return nil
			},
			ReturnWriteFile: func(path string, content []byte) error {
				return nil
			},
		}

		questioner := NewQuestioner(nil, nil)
		configBuilder := NewConfigBuilder()
		init := NewInitializer(
			&InitContext{
				WorkingDir: "/test",
				Stdout:     stdout,
				Stderr:     stderr,
			},
			mockFs,
			questioner,
			configBuilder,
			templateBuilderFactory,
		)
		mockey.Mock(mockey.GetMethod(questioner, "Ask")).Return(&Answer{}, nil).Build()
		mockey.Mock(mockey.GetMethod(configBuilder, "Build")).Return("", nil).Build()
		assert.Equal(0, init.Run())
		assert.Equal("", stderr.String())
	})
}
