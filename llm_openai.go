package main

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
)

type OpenaiLlm struct {
	Client *openai.Client
	Model  string
	Tools  []Tool
}

var _ ILlm = (*OpenaiLlm)(nil)

func (a *OpenaiLlm) UpdateModel(model string) error {
	a.Model = model
	return nil
}

func (a *OpenaiLlm) WithTools(tools []Tool) ILlm {
	return &OpenaiLlm{
		Client: a.Client,
		Model:  a.Model,
		Tools:  tools,
	}
}

func (a *OpenaiLlm) Query(input string) (*LlmQueryResult, error) {
	body := openai.ChatCompletionNewParams{
		Model: openai.F(a.Model),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		}),
		Tools: openai.F(lo.Map(a.Tools, func(tool Tool, i int) openai.ChatCompletionToolParam {
			return openai.ChatCompletionToolParam{
				Type: openai.F(openai.ChatCompletionToolTypeFunction),
				Function: openai.F(openai.FunctionDefinitionParam{
					Name:        openai.F(tool.Name),
					Description: openai.F(tool.Description),
					Parameters: openai.F(shared.FunctionParameters{
						"type": "object",
						"properties": lo.FromEntries(lo.Map(tool.Parameters, func(p Parameter, i int) lo.Entry[string, any] {
							return lo.Entry[string, any]{
								Key: p.Name,
								Value: map[string]any{
									"type":        p.Type,
									"description": p.Description,
								},
							}
						})),
						"required": lo.Filter(lo.Map(tool.Parameters, func(p Parameter, i int) string {
							if p.Required {
								return p.Name
							}
							return ""
						}), func(s string, i int) bool {
							return s != ""
						}),
					}),
				}),
			}
		})),
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
