package main

type IAgent interface {
	query(input string) error
}

type Agent struct {
	tools []Tool
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

func (a *Agent) query(input string) error {
	for {
		// TODO: read input, send request, process them, repeat until done
		break
	}
	return nil
}
