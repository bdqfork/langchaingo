package tools

import (
	"context"

	"github.com/tmc/langchaingo/callbacks"
)

// Tool is a tool for the llm agent to interact with different applications.
type Tool interface {
	Name() string
	Description() string
	Call(ctx context.Context, input string, handlers callbacks.CallbackList) (string, error)
}
