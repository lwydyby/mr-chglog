package cmd

import (
	"os"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/lwydyby/mr-chglog/config"
)

func TestLoad(t *testing.T) {
	mockey.PatchConvey("load", t, func() {
		loader := NewConfigLoader()
		mockey.Mock(os.ReadFile).Return([]byte(`
style: gitlab
template: MRCHANGELOG.tpl.md
title: CHANGELOG
repository_url:
token:
poe_token:
app_id:
app_secret:
bot_title:
# 获取路径: https://open.feishu.cn/api-explorer/cli_xxxxxxxx?apiName=list&from=op_doc&project=im&resource=chat&version=v1
chat_id:
`), nil).Build()
		conf, err := loader.Load("/etc/config/mr.yaml")
		assert.Nil(t, err)
		assert.Equal(t, &config.MRChLogConfig{
			Style:    "gitlab",
			Template: "MRCHANGELOG.tpl.md",
			Title:    "CHANGELOG",
		}, conf)
	})
}
