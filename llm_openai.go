package main

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
)

type OpenaiLlm struct {
	Client *openai.Client
	Model  string
}

var _ ILlm = (*OpenaiLlm)(nil)

func (a *OpenaiLlm) UpdateModel(model string) error {
	a.Model = model
	return nil
}

func (a *OpenaiLlm) Query(input string) (*LlmQueryResult, error) {
	body := openai.ChatCompletionNewParams{
		Model: openai.F(a.Model),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		}),
	}
	chatCompletion, err := a.Client.Chat.Completions.New(
		context.TODO(),
		body,
	)
	if err != nil {
		return nil, eris.Wrap(err, "a.Client.Chat.Completions.New")
	}
	if len(chatCompletion.Choices) == 0 {
		return nil, eris.New("chatCompletion.Choices is empty")
	}
	message := chatCompletion.Choices[0].Message
	return &LlmQueryResult{
		Content: message.Content,
		ToolCalls: lo.Map(message.ToolCalls, func(tc openai.ChatCompletionMessageToolCall, i int) ToolCall {
			return ToolCall{
				ID:        tc.ID,
				Type:      string(tc.Type),
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			}
		}),
	}, nil
}
