package bot

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/lwydyby/mr-chglog/config"
)

func TestGetTenantAccessToken(t *testing.T) {
	// 构造测试用例
	conf := &config.MRChLogConfig{
		AppID:     "test_appid",
		AppSecret: "test_appsecret",
	}

	mockey.Mock(http.NewRequest).Return(nil, nil).Build()
	mockey.Mock(mockey.GetMethod(&http.Client{}, "Do")).Return(&http.Response{}, nil).Build()
	mockey.Mock(io.ReadAll).Return([]byte(`{}`), nil).Build()
	token, err := GetTenantAccessToken(context.Background(), conf)
	assert.Nil(t, err)
	assert.Equal(t, "", token)
}
