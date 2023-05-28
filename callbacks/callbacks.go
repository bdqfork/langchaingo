package callbacks

import (
	"context"

	"github.com/tmc/langchaingo/schema"
)

// CallbackHandler is the interface all callbacks must implement.
type CallbackHandler interface {

	// HandleLLMStart is called at the start of an LLM or Chat Model run, with the prompt(s).
	HandleLLMStart(ctx context.Context, llm string, prompts []string, extraParams map[string]any)

	// HandleLLMNewToken is called when an LLM/ChatModel in `streaming` mode produces a new token.
	HandleLLMNewToken(ctx context.Context, token string)

	// HandleLLMError is called if an LLM/ChatModel run encounters an error.
	HandleLLMError(ctx context.Context, err error)

	// HandleLLMEnd is called at the end of an LLM/ChatModel run, with the output.
	HandleLLMEnd(ctx context.Context, output schema.LLMResult)

	// HandleChatModelStart is called at the start of a Chat Model run, with the prompt(s).
	HandleChatModelStart(ctx context.Context, llm string, messages [][]schema.ChatMessage, extraParams map[string]any)

	// HandleChainStart is called at the start of a Chain run, with the chain name and inputs.
	HandleChainStart(ctx context.Context, chain string, inputs map[string]any)

	// HandleChainError is called if a Chain run encounters an error.
	HandleChainError(ctx context.Context, err error)

	// HandleChainEnd is called at the end of a Chain run, with the outputs.
	HandleChainEnd(ctx context.Context, outputs map[string]any)

	// HandleToolStart is called at the start of a Tool run, with the tool name and input.
	HandleToolStart(ctx context.Context, tool string, input string)

	// HandleToolError is called if a Tool run encounters an error.
	HandleToolError(ctx context.Context, err error)

	// HandleToolEnd is called at the end of a Tool run, with the tool output.
	HandleToolEnd(ctx context.Context, output string)

	HandleText(ctx context.Context, text string)

	// HandleAgentAction is called when an agent is about to execute an action,
	// with the action.
	HandleAgentAction(ctx context.Context, action schema.AgentAction)

	// HandleAgentEnd is called when an agent finishes execution, before it exits.
	// with the final output.
	HandleAgentEnd(ctx context.Context, action schema.AgentFinish)
}

type CallbackList []CallbackHandler

type Id string

const (
	TraceId = "traceId"
)
