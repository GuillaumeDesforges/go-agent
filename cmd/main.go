package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/GuillaumeDesforges/go-agent"
	slogjson "github.com/veqryn/slog-json"
)

func main() {
	level := slog.LevelError
	h := slogjson.NewHandler(os.Stderr, &slogjson.HandlerOptions{AddSource: false, Level: level})
	slog.SetDefault(slog.New(h))

	a := agent.NewReAct()

	clockTool := agent.Tool{
		Name:        "clock",
		Description: "Get current time",
		Parameters:  nil,
		Handler: func(args ...any) (any, error) {
			jb, err := json.Marshal(time.Now().Format(time.Kitchen))
			return string(jb), err
		},
	}
	a.AddTool(clockTool)

	observer := agent.NewConsoleObserver()
	a.RegisterObserver(observer)

	ctx := context.Background()
	err := a.Tell(ctx, "What time is it?")
	if err != nil {
		panic(err)
	}
}
