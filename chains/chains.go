package chains

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/schema"
)

// Chain is the interface all chains must implement.
type Chain interface {
	// Call runs the logic of the chain and returns the output. This method should
	// not be called directly. Use rather the Call function that handles the memory
	// of the chain.
	Call(ctx context.Context, inputs map[string]any, options ...ChainCallOption) (map[string]any, error)
	// GetMemory gets the memory of the chain.
	GetMemory() schema.Memory
	// InputKeys returns the input keys the chain expects.
	GetInputKeys() []string
	// OutputKeys returns the output keys the chain expects.
	GetOutputKeys() []string
}

type ChainCallOptions struct {
	StopWords []string
	Callbacks callbacks.CallbackList
}

// WithStopWords is a ChainCallOption that can be used to set the stop words of the chain.
func WithStopWords(stopWords []string) ChainCallOption {
	return func(options *ChainCallOptions) {
		options.StopWords = stopWords
	}
}

func WithCallbacks(handlers callbacks.CallbackList) ChainCallOption {
	return func(options *ChainCallOptions) {
		options.Callbacks = handlers
	}
}

// ChainCallOption is a function that can be used to modify the behavior of the Call function.
type ChainCallOption func(*ChainCallOptions)

// Call is the function used for calling chains.
func Call(ctx context.Context, c Chain, inputValues map[string]any, options ...ChainCallOption) (map[string]any, error) { //nolint: lll
	if err := validateInputs(c, inputValues); err != nil {
		return nil, err
	}

	fullValues := make(map[string]any, 0)
	for key, value := range inputValues {
		fullValues[key] = value
	}

	newValues, err := c.GetMemory().LoadMemoryVariables(inputValues)
	if err != nil {
		return nil, err
	}

	for key, value := range newValues {
		fullValues[key] = value
	}

	outputValues, err := c.Call(ctx, fullValues, options...)
	if err != nil {
		return nil, err
	}
	if err := validateOutputs(c, outputValues); err != nil {
		return nil, err
	}

	err = c.GetMemory().SaveContext(inputValues, outputValues)
	if err != nil {
		return nil, err
	}

	return outputValues, nil
}

// Run can be used to call a chain if the chain only expects one string input
// and one string output.
func Run(ctx context.Context, c Chain, input string, options ...ChainCallOption) (string, error) {
	inputKeys := c.GetInputKeys()
	if len(inputKeys) != 1 {
		return "", ErrMultipleInputsInRun
	}

	outputKeys := c.GetOutputKeys()
	if len(outputKeys) != 1 {
		return "", ErrMultipleOutputsInRun
	}

	inputValues := map[string]any{inputKeys[0]: input}
	outputValues, err := Call(ctx, c, inputValues, options...)
	if err != nil {
		return "", err
	}

	outputValue, ok := outputValues[outputKeys[0]].(string)
	if !ok {
		return "", ErrWrongOutputTypeInRun
	}

	return outputValue, nil
}

func validateInputs(c Chain, inputValues map[string]any) error {
	for _, k := range c.GetInputKeys() {
		if _, ok := inputValues[k]; !ok {
			return fmt.Errorf("%w: %w: %v", ErrInvalidInputValues, ErrMissingInputValues, k)
		}
	}
	return nil
}

func validateOutputs(c Chain, outputValues map[string]any) error {
	for _, k := range c.GetOutputKeys() {
		if _, ok := outputValues[k]; !ok {
			return fmt.Errorf("%w: %v", ErrInvalidOutputValues, k)
		}
	}
	return nil
}
