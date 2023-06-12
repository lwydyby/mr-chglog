package cmd

import (
	"fmt"
	"strings"
)

// ConfigBuilder ...
type ConfigBuilder interface {
	Builder
}

type configBuilderImpl struct{}

// NewConfigBuilder ...
func NewConfigBuilder() ConfigBuilder {
	return &configBuilderImpl{}
}

// Build ...
func (*configBuilderImpl) Build(ans *Answer) (string, error) {
	repoURL := strings.TrimRight(ans.RepositoryURL, "/")
	if repoURL == "" {
		repoURL = "\"\""
	}

	config := fmt.Sprintf(`style: %s
template: %s
title: CHANGELOG
repository_url: %s
token: %s
poe_token: %s
need_robot: %v
app_id: %s
app_secret: %s
bot_title: %s
# 获取路径: https://open.feishu.cn/api-explorer/cli_a2e211279cbb900c?apiName=list&from=op_doc&project=im&resource=chat&version=v1
chat_id:
`,
		ans.Style,
		defaultTemplateFilename,
		repoURL,
		ans.Token,
		ans.POEToken,
		ans.NeedRobot,
		ans.AppID,
		ans.AppSecret,
		ans.BotTitle,
	)
	for i := range ans.ChatID {
		config += fmt.Sprintf("  - %s\n", ans.ChatID[i])
	}
	return config, nil
}
