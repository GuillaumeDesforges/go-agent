package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
	"github.com/samber/lo"
)

type ReAct struct {
	client openai.Client
	tools  []Tool
}

func NewReAct() *ReAct {
	client := openai.NewClient()
	return &ReAct{
		client: client,
	}
}

func (a *ReAct) AddTool(tool Tool) {
	a.tools = append(a.tools, tool)
}

func (a *ReAct) Tell(ctx context.Context, message string) error {
	model := "o3-mini"
	tools := lo.Map(a.tools, func(tool Tool, _ int) responses.ToolUnionParam {
		return tool.ToParam()
	})

	var params responses.ResponseNewParams
	params = responses.ResponseNewParams{
		Model: model,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(message),
		},
		Tools: tools,
	}
	var response *responses.Response
	var err error
	response, err = a.client.Responses.New(ctx, params)
	if err != nil {
		return fmt.Errorf("a.client.Responses.New: %w", err)
	}
	a.onResponse(response)

	// iterate on tools
	params = responses.ResponseNewParams{
		PreviousResponseID: openai.String(response.ID),
		Model:              model,
		Tools:              tools,
	}
	for _, output := range response.Output {
		if output.Type != "function_call" {
			continue
		}
		call := output.AsFunctionCall()
		slog.Debug("Tell", "call", call)
		for _, tool := range a.tools {
			if tool.Name == call.Name {
				observation, err := tool.Handler()
				if err != nil {
					return fmt.Errorf("tool.Handler: %w", err)
				}
				observationJson, err := json.Marshal(observation)
				if err != nil {
					return fmt.Errorf("json.Marshal: %w", err)
				}
				params.Input = responses.ResponseNewParamsInputUnion{
					OfInputItemList: []responses.ResponseInputItemUnionParam{
						{
							OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
								CallID: call.CallID,
								Output: string(observationJson),
							},
						},
					},
				}
			}
		}
	}

	response, err = a.client.Responses.New(ctx, params)
	if err != nil {
		return fmt.Errorf("a.client.Responses.New: %w", err)
	}
	a.onResponse(response)

	return nil
}

func (a *ReAct) onResponse(response *responses.Response) {
	slog.Debug("onResponse", "response.Reasoning.GenerateSummary", response.Reasoning.GenerateSummary)
	for _, output := range response.Output {
		for _, content := range output.Content {
			if content.Type == "output_text" {
				fmt.Printf("%s\n", content.Text)
			}
		}
	}
}
