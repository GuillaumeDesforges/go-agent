package agent

import (
	"fmt"

	"github.com/openai/openai-go/responses"
)

type ConsoleObserver struct {
	PrefixMessage string
}

func NewConsoleObserver() *ConsoleObserver {
	return &ConsoleObserver{
		PrefixMessage: "> ",
	}
}

var _ Observer = (*ConsoleObserver)(nil)

func (c *ConsoleObserver) NotifyResponse(r *responses.Response) {
	for _, output := range r.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				fmt.Printf("%s%s\n", c.PrefixMessage, content.Text)
			}
		}
	}
}
