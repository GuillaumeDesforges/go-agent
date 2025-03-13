package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Controller struct {
	Model *Model
	Agent *Agent

	updateView func()
}

func (c *Controller) WithUpdateView(updateView func()) *Controller {
	c.updateView = updateView
	return c
}

func (c *Controller) LlmModelChanged(model string) {
	c.Model.LlmModel = model
	c.Agent.Llm.UpdateModel(model)
	c.updateView()
}

func (c *Controller) UserQuerySent(message string) {
	c.Model.Conversation = append(c.Model.Conversation, message)
	c.Model.UserInputDisabled = true
	c.updateView()
	defer func() {
		c.Model.UserInputDisabled = false
		c.updateView()
	}()

	responses := make(chan string)

	errGroup, _ := errgroup.WithContext(context.TODO())
	errGroup.Go(func() error {
		err := c.Agent.Query(message, responses)
		close(responses)
		return err
	})
	errGroup.Go(func() error {
		for response := range responses {
			c.AgentTextSent(response)
		}
		return nil
	})
	err := errGroup.Wait()
	if err != nil {
		c.AgentErrorSent(err)
		return
	}
}

func (c *Controller) AgentTextSent(message string) {
	text := fmt.Sprintf("Agent: %s", message)
	c.Model.Conversation = append(c.Model.Conversation, text)
	c.updateView()
}

func (c *Controller) AgentErrorSent(err error) {
	text := fmt.Sprintf("ERROR: %s", err.Error())
	c.Model.Conversation = append(c.Model.Conversation, text)
	c.updateView()
}
