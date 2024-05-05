package breakmanagerui

import (
	"github.com/charmbracelet/huh"
)

type scribble struct {
	text string
	form huh.Form
}

func New() scribble {
	var text string
	form := *huh.NewForm(huh.NewGroup(huh.NewText().Title("Current Thoughts").CharLimit(1000).Value(&text)))
	return scribble{
		form: form,
		text: "",
	}
}
