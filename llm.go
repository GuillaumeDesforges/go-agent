package main

type LlmQueryResult struct {
	Content   string
	ToolCalls []ToolCall
}

type ToolCall struct {
	ID        string
	Type      string
	Name      string
	Arguments string
}

type ILlm interface {
	UpdateModel(model string) error
	Query(input string) (*LlmQueryResult, error)
}
