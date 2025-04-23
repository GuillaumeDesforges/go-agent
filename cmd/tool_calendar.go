package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/GuillaumeDesforges/go-agent"
)

type CalendarTools struct {
	tools []agent.Tool
}

func NewCalendarTools() *CalendarTools {
	return &CalendarTools{
		tools: []agent.Tool{
			{
				Name:        "calendar-search",
				Description: "Search events in the user's calendar",
				Parameters: []agent.ToolParameters{
					{
						Name:        "from",
						Type:        "string",
						Description: "Date and time from which to search events from (ISO string)",
						Required:    true,
					},
					{
						Name:        "until",
						Type:        "string",
						Description: "Date and time until which to search events from (ISO string)",
						Required:    true,
					},
				},
				Handler: func(args map[string]any) (any, error) {
					fromAny, ok := args["from"]
					if !ok {
						return nil, errors.New("calendar-search: missing argument 'from'")
					}
					fromStr, ok := fromAny.(string)
					if !ok {
						return nil, errors.New("calendar-search: arg 1 must be a string")
					}
					from, err := time.Parse(time.RFC3339, fromStr)
					if err != nil {
						return nil, fmt.Errorf("calendar-search: time.Parse(from): %w", err)
					}

					toAny, ok := args["until"]
					if !ok {
						return nil, errors.New("calendar-search: missing argument 'until'")
					}
					toStr, ok := toAny.(string)
					if !ok {
						return nil, errors.New("calendar-search: arg 2 must be a string")
					}
					to, err := time.Parse(time.RFC3339, toStr)
					if err != nil {
						return nil, fmt.Errorf("calendar-search: time.Parse(to): %w", err)
					}

					var events []map[string]any
					someEventTime := time.Now().Truncate(24 * time.Hour).Add(12*time.Hour + 30*time.Minute)
					if someEventTime.After(from) && someEventTime.Before(to) {
						events = append(events, map[string]any{
							"name": "Lunch with Theo",
						})
					}

					return events, nil
				},
			},
		},
	}
}

func (t *CalendarTools) RegisterTools(a *agent.ReAct) {
	for _, tool := range t.tools {
		a.AddTool(tool)
	}
}
