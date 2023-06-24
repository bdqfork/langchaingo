package openai

import (
	"context"
	"errors"
	"io"
	"os"

	openai "github.com/sashabaranov/go-openai"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

const (
	defaultChatModel       = "gpt-3.5-turbo"
	defaultCompletionModel = "text-davinci-003"
	defaultEmbeddingModel  = "text-embedding-ada-002"

	defaultMaxTokens = 1024
)

var (
	ErrEmptyResponse = errors.New("no response")
	ErrMissingToken  = errors.New("missing the OpenAI API key, set it in the OPENAI_API_KEY environment variable")

	ErrUnexpectedResponseLength = errors.New("unexpected length of response")

	ErrUnexpectedEmbeddingModel = errors.New("unexpected embedding model")
)

type LLM struct {
	client *openai.Client
}

var _ llms.LLM = (*LLM)(nil)

var _ llms.ChatLLM = (*LLM)(nil)

// Call requests a completion for the given prompt.
func (o *LLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	r, err := o.Generate(ctx, []string{prompt}, options...)
	if err != nil {
		return "", err
	}
	if len(r) == 0 {
		return "", ErrEmptyResponse
	}
	return r[0].Text, nil
}

func (o *LLM) Generate(ctx context.Context, prompts []string, options ...llms.CallOption) ([]*llms.Generation, error) {
	opts := llms.CallOptions{MaxTokens: defaultMaxTokens}
	for _, opt := range options {
		opt(&opts)
	}

	model := opts.Model
	if len(model) == 0 {
		model = defaultCompletionModel
	}

	request := openai.CompletionRequest{
		Model:            model,
		Prompt:           prompts[0],
		MaxTokens:        opts.MaxTokens,
		Temperature:      float32(opts.Temperature),
		TopP:             float32(opts.TopP),
		Stream:           opts.StreamingFunc != nil,
		Stop:             opts.StopWords,
		FrequencyPenalty: float32(opts.RepetitionPenalty),
	}

	if request.Stream {
		stream, err := o.client.CreateCompletionStream(ctx, request)
		if err != nil {
			return nil, err
		}
		defer stream.Close()

		output := ""

		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				return nil, err
			}

			if len(resp.Choices) == 0 {
				return nil, ErrEmptyResponse
			}

			text := resp.Choices[0].Text

			err = opts.StreamingFunc(ctx, []byte(text))
			if err != nil {
				return nil, err
			}
			output += text
		}

		return []*llms.Generation{
			{Text: output},
		}, nil
	}

	resp, err := o.client.CreateCompletion(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, ErrEmptyResponse
	}

	return []*llms.Generation{
		{Text: resp.Choices[0].Text},
	}, nil
}

// Chat requests a chat response for the given prompt.
func (o *LLM) Chat(ctx context.Context, messages []schema.ChatMessage, options ...llms.CallOption) (*llms.ChatGeneration, error) { // nolint: lll
	opts := llms.CallOptions{MaxTokens: defaultMaxTokens}
	for _, opt := range options {
		opt(&opts)
	}

	msgs := make([]openai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		msg := openai.ChatCompletionMessage{
			Content: m.GetText(),
		}
		typ := m.GetType()
		switch typ {
		case schema.ChatMessageTypeSystem:
			msg.Role = openai.ChatMessageRoleSystem
		case schema.ChatMessageTypeAI:
			msg.Role = openai.ChatMessageRoleAssistant
		case schema.ChatMessageTypeHuman:
			msg.Role = openai.ChatMessageRoleUser
		case schema.ChatMessageTypeGeneric:
			msg.Role = openai.ChatMessageRoleUser
			// TODO: support name
		}
		msgs[i] = msg
	}

	model := opts.Model
	if len(model) == 0 {
		model = defaultChatModel
	}

	request := openai.ChatCompletionRequest{
		Model:            model,
		Messages:         msgs,
		MaxTokens:        opts.MaxTokens,
		Temperature:      float32(opts.Temperature),
		TopP:             float32(opts.TopP),
		Stream:           opts.StreamingFunc != nil,
		Stop:             opts.StopWords,
		FrequencyPenalty: float32(opts.RepetitionPenalty),
	}

	if request.Stream {
		stream, err := o.client.CreateChatCompletionStream(ctx, request)
		if err != nil {
			return nil, err
		}
		defer stream.Close()

		output := ""

		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				return nil, err
			}

			if len(resp.Choices) == 0 {
				return nil, ErrEmptyResponse
			}

			content := resp.Choices[0].Delta.Content

			err = opts.StreamingFunc(ctx, []byte(content))
			if err != nil {
				return nil, err
			}
			output += content
		}

		return &llms.ChatGeneration{
			Message: &schema.AIChatMessage{
				Text: output,
			},
			// TODO: fill in generation info
		}, nil
	}

	resp, err := o.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, ErrEmptyResponse
	}

	return &llms.ChatGeneration{
		Message: &schema.AIChatMessage{
			Text: resp.Choices[0].Message.Content,
		},
		// TODO: fill in generation info
	}, nil
}

// CreateEmbedding creates embeddings for the given input texts.
func (o *LLM) CreateEmbedding(ctx context.Context, model string, inputTexts []string) ([][]float64, error) {
	if len(model) == 0 {
		model = defaultEmbeddingModel
	}
	embeddingModel, ok := stringToEmbeddingModel[model]
	if !ok {
		return [][]float64{}, ErrUnexpectedEmbeddingModel
	}

	resp, err := o.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: embeddingModel,
		Input: inputTexts,
	})

	if err != nil {
		return [][]float64{}, err
	}

	data := resp.Data

	if len(data) == 0 {
		return [][]float64{}, ErrEmptyResponse
	}

	if len(inputTexts) != len(data) {
		return [][]float64{}, ErrUnexpectedResponseLength
	}

	embeddings := make([][]float64, len(data))
	for i := range data {
		embedding := make([]float64, len(data[i].Embedding))
		for j := range data[i].Embedding {
			embedding[j] = float64(data[i].Embedding[j])
		}
		embeddings[i] = embedding
	}

	return embeddings, nil
}

// New returns a new OpenAI LLM.
func New(opts ...Option) (*LLM, error) {
	options := &options{
		token:   os.Getenv(tokenEnvVarName),
		model:   os.Getenv(modelEnvVarName),
		baseURL: os.Getenv(baseURLEnvVarName),
	}

	for _, opt := range opts {
		opt(options)
	}

	if len(options.token) == 0 {
		return nil, ErrMissingToken
	}

	config := openai.DefaultConfig(options.token)
	config.BaseURL = options.baseURL
	client := openai.NewClientWithConfig(config)

	return &LLM{
		client: client,
	}, nil
}
