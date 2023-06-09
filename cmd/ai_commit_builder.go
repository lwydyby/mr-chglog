package cmd

import (
	"context"
	"fmt"

	"github.com/lwydyby/mr-chglog/ai"
	"github.com/lwydyby/mr-chglog/ai/poe"
	"github.com/lwydyby/mr-chglog/git"
)

const (
	question = `下面我会传给你一个MR的commit diff信息,请帮我归纳总结为一个不超过50字的中文CHANGELOG.回答时请直接返回总结结果即可且不要产生多行数据以及格式字符。diff信息如下: `
)

type AICommitBuilder struct {
	client ai.AIClient
}

func NewAICommitBuilder(tp string, token string) CommitBuilder {
	switch tp {
	case "poe":
		return &AICommitBuilder{
			client: poe.NewPOEClient(token, "ChatGPT"),
		}
	default:
		return nil
	}
}

func (a *AICommitBuilder) BuildCommit(ctx context.Context, mr *git.MergeRequest) string {
	return a.client.SendMessage(ctx, fmt.Sprintf("%s %s", question, mr.Changes))
}
