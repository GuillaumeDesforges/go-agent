package main

import (
	"encoding/json"
	"time"

	"github.com/GuillaumeDesforges/go-agent"
)

type ClockTools struct {
	tools []agent.Tool
}

func NewClockTools() *ClockTools {
	return &ClockTools{
		tools: []agent.Tool{
			{
				Name:        "clock",
				Description: "Get current date and time (RFC3339)",
				Parameters:  nil,
				Handler: func(args map[string]any) (any, error) {
					jb, err := json.Marshal(time.Now().Format(time.RFC3339))
					return string(jb), err
				},
			},
		},
	}
}

func (t *ClockTools) RegisterTools(a *agent.ReAct) {
	for _, tool := range t.tools {
		a.AddTool(tool)
	}
}
