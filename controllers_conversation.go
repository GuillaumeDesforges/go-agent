package main

import "fmt"

type ConversationController struct {
	Model *Model

	controllers *Controllers
	updateView  func()
}

func (c *ConversationController) WithUpdateView(updateFunc func()) *ConversationController {
	c.updateView = updateFunc
	return c
}

func (c *ConversationController) UserInputSent(message string) {
	c.Model.Conversation = append(c.Model.Conversation, message)
	c.updateView()
}

func (c *ConversationController) AgentTextSent(message string) {
	text := fmt.Sprintf("Agent: %s", message)
	c.Model.Conversation = append(c.Model.Conversation, text)
	c.updateView()
}

func (c *ConversationController) AgentErrorSent(err error) {
	text := fmt.Sprintf("ERROR: %s", err.Error())
	c.Model.Conversation = append(c.Model.Conversation, text)
	c.updateView()
}
