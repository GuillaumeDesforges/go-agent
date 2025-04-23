package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
	"github.com/samber/lo"
)

type Reply struct {
	Text string
}

type Conversation interface {
	AddTool(t Tool)
	SendMessage(ctx context.Context, message string, replies chan Reply, toolCalls chan []ToolCall) error
	SendToolResults(ctx context.Context, results []ToolCallResult, replies chan Reply, toolCalls chan []ToolCall) error
}

type OpenaiLlmConversation struct {
	client openai.Client
	model  string
	items  []ConversationItem
	tools  []responses.ToolUnionParam
}

func (c *OpenaiLlmConversation) AddTool(t Tool) {
	tool := responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Name:        t.Name,
			Description: openai.String(t.Description),
			Strict:      true,
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": lo.FromEntries(lo.Map(t.Parameters, func(p ToolParameters, _ int) lo.Entry[string, any] {
					return lo.Entry[string, any]{
						Key: p.Name,
						Value: map[string]any{
							"type":        p.Type,
							"description": p.Description,
						},
					}
				})),
				"required": lo.FilterMap(t.Parameters, func(p ToolParameters, _ int) (string, bool) {
					if p.Required {
						return p.Name, true
					}
					return "", false
				}),
				"additionalProperties": false,
			},
		},
	}

	c.tools = append(c.tools, tool)
}

func NewOpenaiLlmConversation() *OpenaiLlmConversation {
	client := openai.NewClient()
	return &OpenaiLlmConversation{
		client: client,
		items:  []ConversationItem{},
	}
}

func (c *OpenaiLlmConversation) SetModel(model string) {
	c.model = model
}

type ConversationItem struct {
	params   *responses.ResponseNewParams
	response *responses.Response
}

var _ Conversation = (*OpenaiLlmConversation)(nil)

func (c *OpenaiLlmConversation) SendMessage(
	ctx context.Context,
	message string,
	replies chan Reply,
	toolCalls chan []ToolCall,
) error {
	var previousConversationId *string
	if len(c.items) > 0 {
		lastItem := c.items[len(c.items)-1]
		previousConversationId = &lastItem.response.ID
	}

	params := responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				{
					OfMessage: &responses.EasyInputMessageParam{
						Role: "user",
						Content: responses.EasyInputMessageContentUnionParam{
							OfString: openai.String(message),
						},
					},
				},
			},
		},
		Tools: c.tools,
	}
	if previousConversationId != nil {
		params.PreviousResponseID = openai.String(*previousConversationId)
	}

	response, err := c.client.Responses.New(ctx, params)
	if err != nil {
		return fmt.Errorf("a.client.Responses.New: %w", err)
	}

	err = c.sendResponse(response, replies, toolCalls)
	if err != nil {
		return fmt.Errorf("c.sendResponse: %w", err)
	}
	c.items = append(c.items, ConversationItem{
		params:   &params,
		response: response,
	})

	return nil
}

func (c *OpenaiLlmConversation) SendToolResults(
	ctx context.Context,
	results []ToolCallResult,
	replies chan Reply,
	toolCalls chan []ToolCall,
) error {
	var inputItemList responses.ResponseInputParam
	for _, result := range results {
		var output string
		if result.Error != nil {
			output = fmt.Sprintf("%v", result.Error)
		}
		_, err := json.Marshal(result.Result)
		if err != nil {
			return fmt.Errorf("json.Marshall: %w", err)
		}
		item := responses.ResponseInputItemUnionParam{
			OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
				CallID: result.CallId,
				Output: output,
			},
		}
		inputItemList = append(inputItemList, item)
	}

	params := responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: inputItemList,
		},
	}

	response, err := c.client.Responses.New(ctx, params)
	if err != nil {
		return fmt.Errorf("a.client.Responses.New: %w", err)
	}

	err = c.sendResponse(response, replies, toolCalls)
	if err != nil {
		return fmt.Errorf("c.sendResponse: %w", err)
	}
	c.items = append(c.items, ConversationItem{
		params:   &params,
		response: response,
	})

	return nil
}

func (c *OpenaiLlmConversation) sendResponse(
	response *responses.Response,
	outReplies chan Reply,
	outToolCalls chan []ToolCall,
) error {
	var replies []Reply
	var toolCalls []ToolCall
	for _, output := range response.Output {
		if output.Type == "message" {
			message := output.AsMessage()
			for _, content := range message.Content {
				reply := Reply{Text: content.Text}
				replies = append(replies, reply)
			}
		}
		if output.Type == "function_call" {
			call := output.AsFunctionCall()
			var arguments map[string]any
			err := json.Unmarshal([]byte(call.Arguments), &arguments)
			if err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}
			toolCall := ToolCall{
				CallId:    call.CallID,
				ToolName:  call.Name,
				Arguments: arguments,
			}
			toolCalls = append(toolCalls, toolCall)
		}
	}
	for _, reply := range replies {
		outReplies <- reply
	}
	outToolCalls <- toolCalls
	return nil
}
