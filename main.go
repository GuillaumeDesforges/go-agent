package main

import (
	"fmt"

	"github.com/openai/openai-go"
)

type App struct {
	Controllers *Controller
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
		Client: openai.NewClient(),
		Model:  model.LlmModel,
	}
	agent := &Agent{
		Llm: llm,
	}
	controller := &Controller{
		Model: model,
		Agent: agent,
	}

	tui := NewTui(controller)
	controller.WithUpdateView(func() {
		tui.Update(model)
	})
	tui.Update(model)

	app := &App{
		Controllers: controller,
		Tui:         tui,
		Model:       model,
	}
	err := app.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", model)
}
