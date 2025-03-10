package main

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/rotisserie/eris"
)

type OpenaiLlm struct {
	Client   *openai.Client
	LlmModel string
}

var _ ILlm = (*OpenaiLlm)(nil)

type OpenaiCompletionRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

func (a *OpenaiLlm) UpdateModel(model string) error {
	a.LlmModel = model
	return nil
}

func (a *OpenaiLlm) Query(input string) (string, error) {
	body := openai.ChatCompletionNewParams{
		Model: openai.F(a.LlmModel),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		}),
	}
	chatCompletion, err := a.Client.Chat.Completions.New(
		context.TODO(),
		body,
	)
	if err != nil {
		return "", eris.Wrap(err, "a.Client.Chat.Completions.New")
	}
	if len(chatCompletion.Choices) == 0 {
		return "", eris.New("chatCompletion.Choices is empty")
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
