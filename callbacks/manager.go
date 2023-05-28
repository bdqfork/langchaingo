package callbacks

import (
	"context"

	"github.com/tmc/langchaingo/schema"
)

type callbackManager struct {
	handlers CallbackList
}

var _ CallbackHandler = (*callbackManager)(nil)

func (m *callbackManager) HandleLLMStart(ctx context.Context, llm string, prompts []string, extraParams map[string]any) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleLLMStart(ctx, llm, prompts, extraParams)
	}
}

func (m *callbackManager) HandleLLMNewToken(ctx context.Context, token string) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleLLMNewToken(ctx, token)
	}
}

func (m *callbackManager) HandleLLMError(ctx context.Context, err error) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleLLMError(ctx, err)
	}
}

func (m *callbackManager) HandleLLMEnd(ctx context.Context, output schema.LLMResult) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleLLMEnd(ctx, output)
	}
}

func (m *callbackManager) HandleChatModelStart(ctx context.Context, llm string, messages [][]schema.ChatMessage, extraParams map[string]any) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleChatModelStart(ctx, llm, messages, extraParams)
	}
}

func (m *callbackManager) HandleChainStart(ctx context.Context, chain string, inputs map[string]any) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleChainStart(ctx, chain, inputs)
	}
}

func (m *callbackManager) HandleChainError(ctx context.Context, err error) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleChainError(ctx, err)
	}
}

func (m *callbackManager) HandleChainEnd(ctx context.Context, outputs map[string]any) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleChainEnd(ctx, outputs)
	}
}

func (m *callbackManager) HandleToolStart(ctx context.Context, tool string, input string) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleToolStart(ctx, tool, input)
	}
}

func (m *callbackManager) HandleToolError(ctx context.Context, err error) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleToolError(ctx, err)
	}
}

func (m *callbackManager) HandleToolEnd(ctx context.Context, output string) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleToolEnd(ctx, output)
	}
}

func (m *callbackManager) HandleText(ctx context.Context, text string) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleText(ctx, text)
	}
}

func (m *callbackManager) HandleAgentAction(ctx context.Context, action schema.AgentAction) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleAgentAction(ctx, action)
	}
}

func (m *callbackManager) HandleAgentEnd(ctx context.Context, action schema.AgentFinish) {
	if len(m.handlers) == 0 {
		return
	}

	for _, handler := range m.handlers {
		handler.HandleAgentEnd(ctx, action)
	}
}

func NewCallbackManager(handlers CallbackList) CallbackHandler {
	manager := &callbackManager{handlers: handlers}
	return manager
}
