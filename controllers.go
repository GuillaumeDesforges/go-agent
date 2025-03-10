package main

type Controllers struct {
	Model        *Model
	Conversation *ConversationController
	Agent        *AgentController
}

type ControllersParams struct {
	Model *Model
	Agent *Agent
}

func NewControllers(params ControllersParams) *Controllers {
	controllers := &Controllers{
		Conversation: &ConversationController{
			Model: params.Model,
		},
		Agent: &AgentController{
			Model: params.Model,
			Agent: params.Agent,
		},
	}
	controllers.Conversation.controllers = controllers
	controllers.Agent.controllers = controllers
	return controllers
}

func (c *Controllers) WithUpdateFunc(updateFunc func()) *Controllers {
	c.Conversation.WithUpdateView(updateFunc)
	c.Agent.WithUpdateView(updateFunc)
	return c
}
