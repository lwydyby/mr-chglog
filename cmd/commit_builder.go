package cmd

import (
	"context"

	"github.com/lwydyby/mr-chglog/git"
)

type CommitBuilder interface {
	BuildCommit(ctx context.Context, mr *git.MergeRequest) string
}

func commitBuilderFactory(aiType string, token string) CommitBuilder {
	return NewAICommitBuilder(aiType, token)
}
