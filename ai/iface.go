package ai

import (
	"context"
)

type AIClient interface {
	SendMessage(ctx context.Context, message string) string
}
