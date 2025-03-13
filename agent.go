package main

import "github.com/rotisserie/eris"

type IAgent interface {
	UpdateModel(model string) error
	Query(input string, responses chan string) error
}

type Agent struct {
	ILlm
}

var _ IAgent = (*Agent)(nil)

type Tool struct {
	Name        string
	Description string
	Parameters  []Parameter
}

type Parameter struct {
	Properties map[string]interface{}
	Required   []string
}

func (a *Agent) Query(input string, responses chan string) error {
	output, err := a.ILlm.Query(input)
	if err != nil {
		return eris.Wrap(err, "a.Llm.Query")
	}
	responses <- output.Content
	return nil
}
