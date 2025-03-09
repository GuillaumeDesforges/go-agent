package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Tui struct {
	app               *tview.Application
	info              *tview.TextView
	conversation      *tview.TextView
	userInputTextArea *tview.TextArea
}

func NewTui(controllers Controllers) *Tui {
	app := tview.NewApplication()

	info := tview.NewTextView()
	info.SetBorderPadding(2, 2, 2, 2)
	info.SetText(strings.Join(
		[]string{
			"go-agent",
			"Version 0.1.0",
			"",
			"Provider: OpenAI",
			"Model: mini-o1",
		}, "\n",
	))

	conversation := tview.NewTextView()
	conversation.
		SetBorder(true).
		SetTitle("Conversation").
		SetTitleAlign(tview.AlignLeft)
	conversation.SetChangedFunc(func() { app.Draw() })

	userInputTextArea := tview.NewTextArea()
	userInputTextArea.
		SetPlaceholder("Ctrl+Space to send").
		SetBorder(true).
		SetTitle("User").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1)
	userInputTextArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlSpace {
			userInput := userInputTextArea.GetText()
			controllers.Conversation.UserInputSent(userInput)
			userInputTextArea.SetText("", true)
			return nil
		}
		return event
	})

	mainView := tview.NewGrid().
		SetColumns(-1, -2).
		SetRows(-1, -1).
		AddItem(info, 0, 0, 1, 1, 0, 0, false).
		AddItem(conversation, 0, 1, 2, 1, 0, 0, false).
		AddItem(userInputTextArea, 1, 0, 1, 1, 0, 0, true)

	pages := tview.NewPages()
	pages.AddAndSwitchToPage("main", mainView, true)

	app.
		SetRoot(pages, true).
		EnableMouse(true)

	tui := &Tui{
		app:               app,
		info:              info,
		conversation:      conversation,
		userInputTextArea: userInputTextArea,
	}
	return tui
}

func (t *Tui) Run() error {
	return t.app.Run()
}

func (t *Tui) Update(model *Model) {
	t.conversation.SetText(strings.Join(model.Conversation, "\n\n"))
}
