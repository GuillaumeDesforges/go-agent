package agent

import (
	"context"
	"fmt"
	"log/slog"
)

const (
	CHAN_REPLIES_BUFF_SIZE   = 32
	CHAN_TOOLCALLS_BUFF_SIZE = 32
)

type Observer interface {
	NotifyReply(Reply)
}

type ReAct struct {
	tools        []Tool
	observers    []Observer
	conversation Conversation
}

func NewReAct(conversation Conversation) *ReAct {
	return &ReAct{
		conversation: conversation,
	}
}

func (a *ReAct) AddTool(t Tool) {
	a.tools = append(a.tools, t)
	a.conversation.AddTool(t)
}

func (a *ReAct) RegisterObserver(o Observer) {
	a.observers = append(a.observers, o)
}

func (a *ReAct) Tell(ctx context.Context, message string) error {
	// WARNING: may get stuck if the conversation fills up the buffers and we don't consume in a goroutine
	replies := make(chan Reply, CHAN_REPLIES_BUFF_SIZE)
	toolCalls := make(chan []ToolCall, CHAN_TOOLCALLS_BUFF_SIZE)

	err := a.conversation.SendMessage(
		ctx,
		message,
		replies,
		toolCalls,
	)
	close(replies)
	close(toolCalls)
	if err != nil {
		return fmt.Errorf("a.conversation.SendMessage: %w", err)
	}

	for {
		for reply := range replies {
			slog.Debug("a.Tell", "reply", reply)
			for _, o := range a.observers {
				o.NotifyReply(reply)
			}
		}
		var toolCallResults []ToolCallResult
		for toolCallBatch := range toolCalls {
			for _, toolCall := range toolCallBatch {
				slog.Debug("a.Tell", "toolCall", toolCall)
				for _, t := range a.tools {
					toolName := t.Name
					if toolName == toolCall.ToolName {
						result, err := t.Handler(toolCall.Arguments)
						toolCallResults = append(toolCallResults, ToolCallResult{
							ToolName: toolName,
							Result:   result,
							Error:    err,
						})
					}
				}
			}
		}

		if len(toolCallResults) == 0 {
			// no more tool call results to react on
			// NOTE: probably a good time to ask what was the objective and if it fullfilled it?
			break
		}

		replies := make(chan Reply, CHAN_REPLIES_BUFF_SIZE)
		toolCalls := make(chan []ToolCall, CHAN_TOOLCALLS_BUFF_SIZE)
		err = a.conversation.SendToolResults(ctx, toolCallResults, replies, toolCalls)
		close(replies)
		close(toolCalls)
		if err != nil {
			return fmt.Errorf("a.conversation.SendToolResults: %w", err)
		}
	}

	return nil
}
