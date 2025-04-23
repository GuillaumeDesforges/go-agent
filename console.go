package agent

import (
	"fmt"
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

func (c *ConsoleObserver) NotifyReply(r Reply) {
	fmt.Printf("%s%s\n", c.PrefixMessage, r.Text)
}
