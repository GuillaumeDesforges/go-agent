package main

type ConversationController struct {
	Model      *Model
	updateFunc func()
}

func (c *ConversationController) WithUpdateFunc(updateFunc func()) *ConversationController {
	c.updateFunc = updateFunc
	return c
}

func (c *ConversationController) UserInputSent(message string) {
	c.Model.Conversation = append(c.Model.Conversation, message)
	if c.updateFunc == nil {
		panic("updateFunc is nil")
	} else {
		c.updateFunc()
	}
}
