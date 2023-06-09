package poe

import (
	"context"
	"time"

	"github.com/lwydyby/poe-api"
)

type POEClient struct {
	client *poe_api.Client
	bot    string
}

func NewPOEClient(token string, bot string) *POEClient {
	return &POEClient{
		client: poe_api.NewClient(token, nil),
		bot:    bot,
	}
}

func (p *POEClient) SendMessage(ctx context.Context, message string) string {
	ch, err := p.client.SendMessage(p.bot, message, true, 1*time.Minute)
	if err != nil {
		panic(err)
	}
	var lastMessage string
	for c := range ch {
		lastMessage = c["text"].(string)
	}
	return lastMessage
}
