package main

import (
	"fmt"

	"github.com/openai/openai-go"
)

type App struct {
	Controllers *Controllers
	Tui         *Tui
	Model       *Model
}

func (a *App) Run() error {
	return a.Tui.Run()
}

func main() {
	model := &Model{
		LlmModel: openai.ChatModelGPT4oMini,
	}

	llm := &OpenaiLlm{
		Client:   openai.NewClient(),
		LlmModel: model.LlmModel,
	}
	agent := &Agent{
		Llm: llm,
	}
	controllers := NewControllers(ControllersParams{
		Model: model,
		Agent: agent,
	})

	tui := NewTui(*controllers)
	controllers.WithUpdateFunc(func() {
		tui.Update(model)
	})
	tui.Update(model)

	app := &App{
		Controllers: controllers,
		Tui:         tui,
		Model:       model,
	}
	err := app.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", model)
}
