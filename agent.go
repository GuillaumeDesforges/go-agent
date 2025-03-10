package main

import "github.com/rotisserie/eris"

type IAgent interface {
	query(input string) (string, error)
}

type Agent struct {
	Llm ILlm
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

func (a *Agent) query(input string) (string, error) {
	output, err := a.Llm.Query(input)
	if err != nil {
		return "", eris.Wrap(err, "a.Llm.Query")
	}
	return output, nil
}
