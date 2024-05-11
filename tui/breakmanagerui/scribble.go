package breakmanagerui

import (
	"github.com/SamD2021/boba-break/internal/breaklog"
	"github.com/charmbracelet/huh"
)

type scribble struct {
	text   string
	form   huh.Form
	logger *breaklog.FileBreakLogger
}

func New() *scribble {
	l, err := breaklog.NewFileBreakLogger("data/entry.json")
	if err != nil {
		panic(err)
	}
	s := scribble{
		text:   "",
		logger: l,
	}
	s.form = *huh.NewForm(huh.NewGroup(huh.NewText().Title("Current Thoughts").CharLimit(1000).Value(&s.text)))
	return &s
}

func (s scribble) log() {
	s.logger.AddLogEntry(breaklog.NewBreakLogEntry(s.text, ""))
}
