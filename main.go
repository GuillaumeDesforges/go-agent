package main

import "fmt"

type App struct {
	Controllers *Controllers
	Tui         *Tui
	Model       *Model
}

func (a *App) Run() error {
	return a.Tui.Run()
}

func main() {
	model := &Model{}

	controllers := NewControllers(model)

	tui := NewTui(*controllers)
	controllers.WithUpdateFunc(func() {
		fmt.Printf("UPDATE")
		tui.Update(model)
	})

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
