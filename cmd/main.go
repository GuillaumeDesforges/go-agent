package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/GuillaumeDesforges/go-agent"
	slogjson "github.com/veqryn/slog-json"
)

func main() {
	level := slog.LevelDebug
	h := slogjson.NewHandler(os.Stderr, &slogjson.HandlerOptions{AddSource: false, Level: level})
	slog.SetDefault(slog.New(h))

	conversation := agent.NewOpenaiLlmConversation()
	conversation.SetModel("o3-mini")
	a := agent.NewReAct(conversation)

	calendarTools := NewCalendarTools()
	calendarTools.RegisterTools(a)

	clockTools := NewClockTools()
	clockTools.RegisterTools(a)

	observer := agent.NewConsoleObserver()
	a.RegisterObserver(observer)

	ctx := context.Background()
	err := a.Tell(ctx, "What are the events in my calendar today?")
	if err != nil {
		panic(err)
	}
}
