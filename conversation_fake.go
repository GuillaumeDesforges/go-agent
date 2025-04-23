package agent

import "context"

type ConversationFake struct {
	tools []Tool
}

func NewConversationFake() *ConversationFake {
	return &ConversationFake{}
}

var _ Conversation = (*ConversationFake)(nil)

func (c *ConversationFake) AddTool(tool Tool) {
	c.tools = append(c.tools, tool)
}

func (c *ConversationFake) SendMessage(
	ctx context.Context,
	message string,
	replies chan Reply,
	toolCalls chan []ToolCall,
) error {
	return nil
}

func (c *ConversationFake) SendToolResults(
	ctx context.Context,
	toolCallResults []ToolCallResult,
	replies chan Reply,
	toolCalls chan []ToolCall,
) error {
	return nil
}
