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
`,
		ans.Style,
		defaultTemplateFilename,
		repoURL,
		ans.Token,
		ans.POEToken,
	)

	return config, nil
}
