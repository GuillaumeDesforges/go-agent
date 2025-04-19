package agent

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
	"github.com/samber/lo"
)

type Tool struct {
	Name        string
	Description string
	Parameters  []ToolParameters
	Handler     func(args ...any) (any, error)
}

type ToolParameters struct {
	Name     string
	Type     string
	Required bool
}

func (t *Tool) ToParam() responses.ToolUnionParam {
	return responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Name:        t.Name,
			Description: openai.String(t.Description),
			Strict:      true,
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": lo.FromEntries(lo.Map(t.Parameters, func(p ToolParameters, _ int) lo.Entry[string, any] {
					return lo.Entry[string, any]{
						Key: p.Name,
						Value: map[string]any{
							"type": p.Type,
						},
					}
				})),
				"required": lo.FilterMap(t.Parameters, func(p ToolParameters, _ int) (string, bool) {
					if p.Required {
						return p.Name, true
					}
					return "", false
				}),
				"additionalProperties": false,
			},
		},
	}
}
