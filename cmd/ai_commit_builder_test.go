package cmd

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/lwydyby/mr-chglog/ai/poe"
	"github.com/lwydyby/mr-chglog/git"
)

func TestAICommitBuilder_BuildCommit(t *testing.T) {
	type args struct {
		ctx context.Context
		mr  *git.MergeRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple",
			args: args{
				ctx: context.Background(),
				mr:  &git.MergeRequest{},
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &poe.POEClient{}
			mockey.Mock(mockey.GetMethod(client, "SendMessage")).Return(tt.want).Build()
			a := &AICommitBuilder{
				client: client,
			}
			if got := a.BuildCommit(tt.args.ctx, tt.args.mr); got != tt.want {
				t.Errorf("BuildCommit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAICommitBuilder(t *testing.T) {
	type args struct {
		tp    string
		token string
	}
	tests := []struct {
		name string
		args args
		want CommitBuilder
	}{
		{
			name: "simple",
			args: args{
				tp:    "poe",
				token: "",
			},
			want: &AICommitBuilder{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockey.Mock(poe.NewPOEClient).Return(nil).Build()
			assert.True(t, NewAICommitBuilder(tt.args.tp, tt.args.token) != nil)
		})
	}
}
