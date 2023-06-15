package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_configBuilderImpl_Build(t *testing.T) {
	type args struct {
		ans *Answer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple",
			args: args{
				ans: &Answer{},
			},
			want: `style: 
template: MRCHANGELOG.tpl.md
title: CHANGELOG
repository_url: ""
token: 
poe_token: 
need_robot: false
app_id: 
app_secret: 
bot_title: 
# 获取路径: https://open.feishu.cn/api-explorer/cli_a2e211279cbb900c?apiName=list&from=op_doc&project=im&resource=chat&version=v1
chat_id:
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			co := NewConfigBuilder()
			got, _ := co.Build(tt.args.ans)
			assert.Equalf(t, tt.want, got, "Build(%v)", tt.args.ans)
		})
	}
}
