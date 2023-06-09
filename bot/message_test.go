package bot

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestSendAlertMessage(t *testing.T) {
	mockey.PatchConvey("SendAlertMessage", t, func() {
		mockey.Mock(mockey.GetMethod(&http.Client{}, "Do")).Return(&http.Response{}, nil).Build()
		mockey.Mock(io.ReadAll).Return([]byte(`{"code":0,"data":{}}`), nil).Build()
		err := SendAlertMessage(context.Background(), "", "", "", "")
		assert.Nil(t, err)
	})
}
