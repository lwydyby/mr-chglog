package poe

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	poe_api "github.com/lwydyby/poe-api"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	mc := &poe_api.Client{}
	ch := make(chan map[string]interface{}, 1)
	ch <- map[string]interface{}{
		"text": "hello",
	}
	close(ch)
	mockey.Mock(poe_api.NewClient).Return(mc).Build()
	mockey.Mock(mockey.GetMethod(mc, "SendMessage")).Return(ch, nil).Build()
	client := NewPOEClient("", "")
	result := client.SendMessage(context.Background(), "test")
	assert.Equal(t, "hello", result)
}
