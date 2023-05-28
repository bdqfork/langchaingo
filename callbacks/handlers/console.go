package handlers

import (
	"context"
	"fmt"

	"github.com/fatih/color"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/schema"
)

// ConsoleCallbackHandler can be used to track the input and output of the entire call chain.
type ConsoleCallbackHandler struct {
}

var _ callbacks.CallbackHandler = (*ConsoleCallbackHandler)(nil)

func (c *ConsoleCallbackHandler) HandleLLMStart(ctx context.Context, llm string, prompts []string, extraParams map[string]any) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s Entering LLM run with input: %+v\n", color.GreenString("[llm/start]"), traceId, prompts)
}

func (c *ConsoleCallbackHandler) HandleLLMNewToken(ctx context.Context, token string) {}

func (c *ConsoleCallbackHandler) HandleLLMError(ctx context.Context, err error) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s LLM run errored with error: %+v\n", color.RedString("[llm/error]"), traceId, err)
}

func (c *ConsoleCallbackHandler) HandleLLMEnd(ctx context.Context, output schema.LLMResult) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s Exiting LLM run with output: %+v\n", color.CyanString("[llm/end]"), traceId, output)
}

func (c *ConsoleCallbackHandler) HandleChatModelStart(ctx context.Context, llm string, messages [][]schema.ChatMessage, extraParams map[string]any) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s Entering LLM run with input: %+v\n", color.GreenString("[llm/start]"), traceId, messages)
}

func (c *ConsoleCallbackHandler) HandleChainStart(ctx context.Context, chain string, inputs map[string]any) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s Entering Chain run with input: %+v\n", color.GreenString("[chain/start]"), traceId, inputs)
}

func (t *ConsoleCallbackHandler) HandleChainError(ctx context.Context, err error) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s Chain run errored with error: %+v\n", color.RedString("[chain/error]"), traceId, err)
}

func (t *ConsoleCallbackHandler) HandleChainEnd(ctx context.Context, outputs map[string]any) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s Exiting Chain run with output: %+v\n", color.CyanString("[chain/end]"), traceId, outputs)
}

func (c *ConsoleCallbackHandler) HandleToolStart(ctx context.Context, tool string, input string) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s Entering Tool run with input: %+v\n", color.GreenString("[tool/start]"), traceId, input)
}

func (c *ConsoleCallbackHandler) HandleToolError(ctx context.Context, err error) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s Tool run errored with error: %+v\n", color.RedString("[tool/error]"), traceId, err)
}

func (c *ConsoleCallbackHandler) HandleToolEnd(ctx context.Context, output string) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s Exiting Tool run with output: %+v\n", color.CyanString("[tool/end]"), traceId, output)
}

func (c *ConsoleCallbackHandler) HandleText(ctx context.Context, text string) {}

func (c *ConsoleCallbackHandler) HandleAgentAction(ctx context.Context, action schema.AgentAction) {
	traceId := ""
	if val, ok := ctx.Value(callbacks.TraceId).(string); ok {
		traceId = val
	}
	traceId = color.MagentaString(fmt.Sprintf("[%s]", traceId))
	fmt.Printf("%s %s Agent selected action: %+v\n", color.BlueString("[agent/action]"), traceId, action)
}

func (c *ConsoleCallbackHandler) HandleAgentEnd(ctx context.Context, action schema.AgentFinish) {
}
