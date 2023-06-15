package bot

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestSendAlertMessage(t *testing.T) {
	mockey.PatchConvey("SendAlertMessage", t, func() {
		mockey.Mock(sendMessage).Return(&MessageItem{}, nil).Build()
		mockey.Mock(pinMessage).Return().Build()
		err := SendAlertMessage(context.Background(), "", "", "", "")
		assert.Nil(t, err)
	})
}
