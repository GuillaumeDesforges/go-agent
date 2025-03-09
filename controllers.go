package main

type Controllers struct {
	Model        *Model
	Conversation *ConversationController
}

func NewControllers(model *Model) *Controllers {
	return &Controllers{
		Conversation: &ConversationController{
			Model: model,
		},
	}
}

func (c *Controllers) WithUpdateFunc(updateFunc func()) *Controllers {
	c.Conversation.WithUpdateFunc(updateFunc)
	return c
}
