package main

import (
	"encoding/json"

	"github.com/rotisserie/eris"
	"github.com/samber/lo"
)

type IAgent interface {
	UpdateModel(model string) error
	Query(input string, responses chan string) error
}

type Agent struct {
	ILlm
	Tools []Tool
}

type Tool struct {
	Name        string
	Description string
	Parameters  []Parameter
	f           func(arguments map[string]any) (string, error)
}

type Parameter struct {
	Name        string
	Type        string
	Description string
	Required    bool
}

var _ IAgent = (*Agent)(nil)

func (a *Agent) Query(input string, responses chan string) error {
	output, err := a.ILlm.WithTools(a.Tools).Query(input)
	if err != nil {
		return eris.Wrap(err, "a.Llm.Query")
	}
	responses <- output.Content
	for _, toolCall := range output.ToolCalls {
		tool, found := lo.Find(a.Tools, func(item Tool) bool {
			return item.Name == toolCall.Name
		})
		if !found {
			return eris.Errorf("tool %s not found", toolCall.Name)
		}
		arguments := make(map[string]any)
		err := json.Unmarshal([]byte(toolCall.Arguments), &arguments)
		if err != nil {
			return eris.Wrap(err, "json.Unmarshal")
		}
		_, err = tool.f(arguments)
		if err != nil {
			return eris.Wrap(err, "tool.f")
		}
	}
	return nil
}
