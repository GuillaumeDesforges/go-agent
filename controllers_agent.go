package main

type AgentController struct {
	Model *Model
	Agent *Agent

	controllers *Controllers
	updateView  func()
}

func (c *AgentController) WithUpdateView(updateFunc func()) *AgentController {
	c.updateView = updateFunc
	return c
}

func (c *AgentController) UserQuerySent(message string) {
	response, err := c.Agent.query(message)
	if err != nil {
		c.controllers.Conversation.AgentErrorSent(err)
		return
	}
	c.controllers.Conversation.AgentTextSent(response)
}

func (c *AgentController) LlmModelChanged(model string) {
	c.Model.LlmModel = model
	c.controllers.Agent.Agent.Llm.UpdateModel(model)
	c.updateView()
}
